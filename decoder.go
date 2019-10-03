package gogm

import (
	"errors"
	"fmt"
	dsl "github.com/mindstand/go-cypherdsl"
	neo "github.com/mindstand/golang-neo4j-bolt-driver"
	"github.com/mindstand/golang-neo4j-bolt-driver/structures/graph"
	"reflect"
	"strings"
	"time"
)

func decodeNeoRows(rows neo.Rows, respObj interface{}) error {
	defer rows.Close()

	arr, err := dsl.RowsTo2DInterfaceArray(rows)
	if err != nil {
		return err
	}

	return decode(arr, respObj)
}

//example query `match p=(n)-[*0..5]-() return p`
//decodes raw path response from driver
func decode(rawArr [][]interface{}, respObj interface{}) (err error) {
	//check nil params
	if rawArr == nil {
		return fmt.Errorf("rawArr can not be nil, %w", ErrInvalidParams)
	}

	//check empty
	if len(rawArr) == 0 {
		return fmt.Errorf("nothing returned from driver, %w", ErrNotFound)
	}

	//we're doing reflection now, lets set up a panic recovery
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v - PANIC RECOVERY - %w", r, ErrInternal)
		}
	}()

	if respObj == nil {
		return fmt.Errorf("response object can not be nil - %w", ErrInvalidParams)
	}

	rv := reflect.ValueOf(respObj)
	rt := reflect.TypeOf(respObj)

	primaryLabel := getPrimaryLabel(rt)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("invalid resp type %T - %w", respObj, ErrInvalidParams)
	}

	//todo optimize with set array size
	var paths []*graph.Path
	var strictRels []*graph.Relationship
	var isolatedNodes []*graph.Node

	for _, arr := range rawArr {
		for _, graphType := range arr {
			switch graphType.(type) {
			case graph.Path:
				convP := graphType.(graph.Path)
				paths = append(paths, &convP)
				break
			case graph.Relationship:
				convR := graphType.(graph.Relationship)
				strictRels = append(strictRels, &convR)
				break
			case graph.Node:
				convN := graphType.(graph.Node)
				isolatedNodes = append(isolatedNodes, &convN)
				break
			default:
				return fmt.Errorf("%T unsupported type, %w", graphType, ErrInternal)
			}
		}
	}
	nodeLookup := make(map[int64]*reflect.Value)
	var pks []int64
	rels := make(map[int64]*neoEdgeConfig)
	labelLookup := map[int64]string{}

	if paths != nil && len(paths) != 0 {
		err = sortPaths(paths, &nodeLookup, &rels, &pks, primaryLabel)
		if err != nil {
			return err
		}
	}

	if isolatedNodes != nil && len(isolatedNodes) != 0 {
		err = sortIsolatedNodes(isolatedNodes, &labelLookup, &nodeLookup, &pks, primaryLabel)
		if err != nil {
			return err
		}
	}

	if strictRels != nil && len(strictRels) != 0 {
		err = sortStrictRels(strictRels, &labelLookup, &rels)
		if err != nil {
			return err
		}
	}

	//build relationships
	for _, relationConfig := range rels {
		//todo figure out why this is broken
		if relationConfig.StartNodeType == "" || relationConfig.EndNodeType == "" {
			continue
		}

		start, _, err := getValueAndConfig(relationConfig.StartNodeId, relationConfig.StartNodeType, nodeLookup)
		if err != nil {
			return err
		}

		end, _, err := getValueAndConfig(relationConfig.EndNodeId, relationConfig.EndNodeType, nodeLookup)
		if err != nil {
			return err
		}

		startConfig, endConfig, err := mappedRelations.GetConfigs(relationConfig.StartNodeType, relationConfig.EndNodeType,
			relationConfig.EndNodeType, relationConfig.StartNodeType, relationConfig.Type)
		if err != nil {
			return err
		}

		if startConfig.UsesEdgeNode {
			var typeConfig structDecoratorConfig

			it := startConfig.Type

			//get the actual type if its a slice
			if it.Kind() == reflect.Slice {
				it = it.Elem()
			}

			label := ""

			if it.Kind() == reflect.Ptr {
				label = it.Elem().Name()
			} else {
				label = it.Name()
				it = reflect.PtrTo(it)
			}

			temp, ok := mappedTypes.Get(label) // mappedTypes[boltNode.Labels[0]]
			if !ok {
				return fmt.Errorf("can not find mapping for node with label %s - %w", label, ErrInternal)
			}

			typeConfig = temp.(structDecoratorConfig)
			if !ok {
				return fmt.Errorf("unable to cast [%T] to structDecoratorConfig - %w", temp, ErrInternal)
			}

			//create value
			val, err := convertToValue(relationConfig.Id, typeConfig, relationConfig.Obj, it)
			if err != nil {
				return err
			}

			var startCall reflect.Value
			var endCall reflect.Value

			if start.Kind() != reflect.Ptr {
				startCall = start.Addr()
			} else {
				startCall = *start
			}

			if end.Kind() != reflect.Ptr {
				endCall = end.Addr()
			} else {
				endCall = *end
			}

			//can ensure that it implements proper interface if it made it this far
			res := val.MethodByName("SetStartNode").Call([]reflect.Value{startCall})
			if res == nil || len(res) == 0 {
				return fmt.Errorf("invalid response from edge callback - %w", err)
			} else if !res[0].IsNil() {
				return fmt.Errorf("failed call to SetStartNode - %w", res[0].Interface().(error))
			}

			res = val.MethodByName("SetEndNode").Call([]reflect.Value{endCall})
			if res == nil || len(res) == 0 {
				return fmt.Errorf("invalid response from edge callback - %w", err)
			} else if !res[0].IsNil() {
				return fmt.Errorf("failed call to SetEndNode - %w", res[0].Interface().(error))
			}

			//relate end-start
			if reflect.Indirect(*end).FieldByName(endConfig.FieldName).Kind() == reflect.Slice {
				reflect.Indirect(*end).FieldByName(endConfig.FieldName).Set(reflect.Append(reflect.Indirect(*end).FieldByName(endConfig.FieldName), *val))
			} else {
				//non slice relationships are already asserted to be pointers
				end.FieldByName(endConfig.FieldName).Set(*val)
			}

			//relate start-start
			if reflect.Indirect(*start).FieldByName(startConfig.FieldName).Kind() == reflect.Slice {
				reflect.Indirect(*start).FieldByName(startConfig.FieldName).Set(reflect.Append(reflect.Indirect(*start).FieldByName(startConfig.FieldName), *val))
			} else {
				start.FieldByName(startConfig.FieldName).Set(*val)
			}
		} else {
			if end.FieldByName(endConfig.FieldName).Kind() == reflect.Slice {
				end.FieldByName(endConfig.FieldName).Set(reflect.Append(end.FieldByName(endConfig.FieldName), start.Addr()))
			} else {
				end.FieldByName(endConfig.FieldName).Set(start.Addr())
			}

			//relate end-start
			if start.FieldByName(startConfig.FieldName).Kind() == reflect.Slice {
				start.FieldByName(startConfig.FieldName).Set(reflect.Append(start.FieldByName(startConfig.FieldName), end.Addr()))
			} else {
				start.FieldByName(startConfig.FieldName).Set(end.Addr())
			}
		}
	}

	//handle if its returning a slice -- validation has been done at an earlier step
	if rt.Elem().Kind() == reflect.Slice {

		reflection := reflect.MakeSlice(rt.Elem(), 0, cap(pks))

		reflectionValue := reflect.New(reflection.Type())
		reflectionValue.Elem().Set(reflection)

		slicePtr := reflect.ValueOf(reflectionValue.Interface())

		sliceValuePtr := slicePtr.Elem()

		sliceType := rt.Elem().Elem()

		for _, id := range pks {
			val, ok := nodeLookup[id]
			if !ok {
				return fmt.Errorf("cannot find value with id (%v)", id)
			}

			//handle slice of pointers
			if sliceType.Kind() == reflect.Ptr {
				sliceValuePtr.Set(reflect.Append(sliceValuePtr, val.Addr()))
			} else {
				sliceValuePtr.Set(reflect.Append(sliceValuePtr, *val))
			}
		}

		reflect.Indirect(rv).Set(sliceValuePtr)

		return err
	} else {
		//handles single -- already checked to make sure p2 is at least 1
		reflect.Indirect(rv).Set(*nodeLookup[pks[0]])

		return err
	}
}

func getPrimaryLabel(rt reflect.Type) string {
	//assume its already a pointer
	rt = rt.Elem()

	if rt.Kind() == reflect.Slice {
		rt = rt.Elem()
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
	}

	return rt.Name()
}

func sortIsolatedNodes(isolatedNodes []*graph.Node, labelLookup *map[int64]string, nodeLookup *map[int64]*reflect.Value, pks *[]int64, pkLabel string) error {
	if isolatedNodes == nil {
		return fmt.Errorf("isolatedNodes can not be nil, %w", ErrInternal)
	}

	for _, node := range isolatedNodes {
		if node == nil {
			return fmt.Errorf("node should not be nil, %w", ErrInternal)
		}

		//check if node has already been found by another process
		if _, ok := (*nodeLookup)[node.NodeIdentity]; !ok {
			//if it hasn't, map it
			val, err := convertNodeToValue(*node)
			if err != nil {
				return err
			}

			(*nodeLookup)[node.NodeIdentity] = val

			//primary to return
			if node.Labels != nil && len(node.Labels) != 0 && node.Labels[0] == pkLabel {
				*pks = append(*pks, node.NodeIdentity)
			}

			//set label map
			if _, ok := (*labelLookup)[node.NodeIdentity]; !ok && len(node.Labels) != 0 && node.Labels[0] == pkLabel {
				(*labelLookup)[node.NodeIdentity] = node.Labels[0]
			}
		}
	}

	return nil
}

func sortStrictRels(strictRels []*graph.Relationship, labelLookup *map[int64]string, rels *map[int64]*neoEdgeConfig) error {
	if strictRels == nil {
		return fmt.Errorf("paths is empty, that shouldn't have happened, %w", ErrInternal)
	}

	for _, rel := range strictRels {
		if rel == nil {
			return errors.New("path can not be nil")
		}

		if _, ok := (*rels)[rel.RelIdentity]; !ok {
			startLabel, ok := (*labelLookup)[rel.StartNodeIdentity]
			if !ok {
				return fmt.Errorf("label not found for node [%v], %w", rel.StartNodeIdentity, ErrInternal)
			}

			endLabel, ok := (*labelLookup)[rel.EndNodeIdentity]
			if !ok {
				return fmt.Errorf("label not found for node [%v], %w", rel.EndNodeIdentity, ErrInternal)
			}

			(*rels)[rel.RelIdentity] = &neoEdgeConfig{
				Id:            rel.RelIdentity,
				StartNodeId:   rel.StartNodeIdentity,
				StartNodeType: startLabel,
				EndNodeId:     rel.StartNodeIdentity,
				EndNodeType:   endLabel,
				Obj:           rel.Properties,
				Type:          rel.Type,
			}
		}
	}

	return nil
}

func sortPaths(paths []*graph.Path, nodeLookup *map[int64]*reflect.Value, rels *map[int64]*neoEdgeConfig, pks *[]int64, pkLabel string) error {
	if paths == nil {
		return fmt.Errorf("paths is empty, that shouldn't have happened, %w", ErrInternal)
	}

	for _, path := range paths {
		if path == nil {
			return errors.New("path can not be nil")
		}

		if path.Nodes == nil || len(path.Nodes) == 0 {
			return fmt.Errorf("no nodes found, %w", ErrNotFound)
		}

		for _, node := range path.Nodes {

			if _, ok := (*nodeLookup)[node.NodeIdentity]; !ok {
				//we haven't parsed this one yet, lets do that now
				val, err := convertNodeToValue(node)
				if err != nil {
					return err
				}

				(*nodeLookup)[node.NodeIdentity] = val

				//primary to return
				if node.Labels != nil && len(node.Labels) != 0 && node.Labels[0] == pkLabel {
					*pks = append(*pks, node.NodeIdentity)
				}
			}
		}

		//handle the relationships
		//sequence is [node_id, edge_id (negative if its in the wrong direction), repeat....]
		if path.Sequence != nil && len(path.Sequence) != 0 && path.Relationships != nil && len(path.Relationships) == (len(path.Sequence)/2) {
			seqPrime := append([]int{0}, path.Sequence...)
			seqPrimeLen := len(seqPrime)

			for i := 0; i+2 < seqPrimeLen; i += 2 {
				startPosIndex := seqPrime[i]
				edgeIndex := seqPrime[i+1]
				endPosIndex := seqPrime[i+2]

				var startId int
				var endId int
				var edgeId int

				//keep order
				if edgeIndex >= 0 {
					startId = startPosIndex
					endId = endPosIndex
					edgeId = edgeIndex
				} else {
					//reverse
					startId = endPosIndex
					endId = startPosIndex
					edgeId = -edgeIndex
				}

				startNode := path.Nodes[startId]
				endNode := path.Nodes[endId]
				rel := path.Relationships[edgeId-1] //offset for the array

				if _, ok := (*rels)[rel.RelIdentity]; !ok {
					(*rels)[rel.RelIdentity] = &neoEdgeConfig{
						Id:            rel.RelIdentity,
						StartNodeId:   startNode.NodeIdentity,
						StartNodeType: startNode.Labels[0],
						EndNodeId:     endNode.NodeIdentity,
						EndNodeType:   endNode.Labels[0],
						Obj:           rel.Properties,
						Type:          rel.Type,
					}
				}
			}
		}
	}

	return nil
}

func getValueAndConfig(id int64, t string, nodeLookup map[int64]*reflect.Value) (val *reflect.Value, conf structDecoratorConfig, err error) {
	var ok bool

	val, ok = nodeLookup[id]
	if !ok {
		return nil, structDecoratorConfig{}, fmt.Errorf("value for id (%v) not found", id)
	}

	temp, ok := mappedTypes.Get(t)
	if !ok {
		return nil, structDecoratorConfig{}, fmt.Errorf("no config found for type (%s)", t)
	}

	conf, ok = temp.(structDecoratorConfig)
	if !ok {
		return nil, structDecoratorConfig{}, errors.New("unable to cast to structDecoratorConfig")
	}

	return
}

var sliceOfEmptyInterface []interface{}
var emptyInterfaceType = reflect.TypeOf(sliceOfEmptyInterface).Elem()

func convertToValue(graphId int64, conf structDecoratorConfig, props map[string]interface{}, rtype reflect.Type) (valss *reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if rtype == nil {
		return nil, errors.New("rtype can not be nil")
	}

	isPtr := false
	if rtype.Kind() == reflect.Ptr {
		isPtr = true
		rtype = rtype.Elem()
	}

	val := reflect.New(rtype)

	if graphId >= 0 {
		reflect.Indirect(val).FieldByName("Id").Set(reflect.ValueOf(graphId))
	}

	for field, fieldConfig := range conf.Fields {
		if fieldConfig.Name == "id" {
			continue //id is handled above
		}

		//skip if its a relation field
		if fieldConfig.Relationship != "" {
			continue
		}

		if fieldConfig.Properties {
			mapType := reflect.MapOf(reflect.TypeOf(""), emptyInterfaceType)
			mapVal := reflect.MakeMap(mapType)

			for k, v := range props {
				if !strings.Contains(k, fieldConfig.Name) {
					//not one of our map fields
					continue
				}

				parts := strings.Split(k, ".")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid key [%s]", k)
				}

				mapKey := parts[1]

				mapVal.SetMapIndex(reflect.ValueOf(mapKey), reflect.ValueOf(v))
			}

			reflect.Indirect(val).FieldByName(field).Set(mapVal)
			continue
		}

		raw, ok := props[fieldConfig.Name]
		if !ok {
			if fieldConfig.IsTypeDef {
				log.Debugf("skipping field %s since it is typedeffed and not defined", fieldConfig.Name)
				continue
			}
			return nil, fmt.Errorf("unrecognized field [%s]", fieldConfig.Name)
		}

		if raw == nil {
			continue //its already initialized to 0 value, no need to do anything
		} else {
			if fieldConfig.IsTime {
				timeStr, ok := raw.(string)
				if !ok {
					return nil, errors.New("can not convert interface{} time to string")
				}

				convTime, err := time.Parse(time.RFC3339, timeStr)
				if err != nil {
					return nil, err
				}

				var writeVal reflect.Value

				if fieldConfig.Type.Kind() == reflect.Ptr {
					writeVal = reflect.ValueOf(convTime).Addr()
				} else {
					writeVal = reflect.ValueOf(convTime)
				}

				reflect.Indirect(val).FieldByName(field).Set(writeVal)
			} else {
				rawVal := reflect.ValueOf(raw)
				indirect := reflect.Indirect(val)
				if indirect.FieldByName(field).Type() == rawVal.Type() {
					indirect.FieldByName(field).Set(rawVal)
				} else {
					indirect.FieldByName(field).Set(rawVal.Convert(indirect.FieldByName(field).Type()))
				}
			}
		}
	}

	//if its not a pointer, dereference it
	if !isPtr {
		retV := reflect.Indirect(val)
		return &retV, nil
	}

	return &val, err
}

func convertNodeToValue(boltNode graph.Node) (*reflect.Value, error) {

	if boltNode.Labels == nil || len(boltNode.Labels) == 0 {
		return nil, errors.New("boltNode has no labels")
	}

	var typeConfig structDecoratorConfig

	temp, ok := mappedTypes.Get(boltNode.Labels[0]) // mappedTypes[boltNode.Labels[0]]
	if !ok {
		return nil, fmt.Errorf("can not find mapping for node with label %s", boltNode.Labels[0])
	}

	typeConfig, ok = temp.(structDecoratorConfig)
	if !ok {
		return nil, errors.New("unable to cast to struct decorator config")
	}

	return convertToValue(boltNode.NodeIdentity, typeConfig, boltNode.Properties, typeConfig.Type)
}
