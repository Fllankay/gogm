package main

import (
	"context"
	"fmt"
	"github.com/mindstand/gogm/v2"
	"github.com/mindstand/gogm/v2/examples/movies/domain"
	"reflect"
	"time"
)

type tdString string
type tdInt int

// structs for the example (can also be found in decoder_test.go)
type VertexA struct {
	// provides required node fields
	gogm.BaseUUIDNode

	TestField         string            `gogm:"name=test_field"`
	TestTypeDefString tdString          `gogm:"name=test_type_def_string"`
	TestTypeDefInt    tdInt             `gogm:"name=test_type_def_int"`
	MapProperty       map[string]string `gogm:"name=map_property;properties"`
	SliceProperty     []string          `gogm:"name=slice_property;properties"`
	SingleA           *VertexB          `gogm:"direction=incoming;relationship=test_rel"`
	ManyA             []*VertexB        `gogm:"direction=incoming;relationship=testm2o"`
	MultiA            []*VertexB        `gogm:"direction=incoming;relationship=multib"`
	SingleSpecA       *EdgeC            `gogm:"direction=outgoing;relationship=special_single"`
	MultiSpecA        []*EdgeC          `gogm:"direction=outgoing;relationship=special_multi"`
}

type VertexB struct {
	// provides required node fields
	gogm.BaseUUIDNode

	TestField  string     `gogm:"name=test_field"`
	TestTime   time.Time  `gogm:"name=test_time"`
	Single     *VertexA   `gogm:"direction=outgoing;relationship=test_rel"`
	ManyB      *VertexA   `gogm:"direction=outgoing;relationship=testm2o"`
	Multi      []*VertexA `gogm:"direction=outgoing;relationship=multib"`
	SingleSpec *EdgeC     `gogm:"direction=incoming;relationship=special_single"`
	MultiSpec  []*EdgeC   `gogm:"direction=incoming;relationship=special_multi"`
}

// EdgeC implements Edge
type EdgeC struct {
	// provides required node fields
	gogm.BaseUUIDNode

	Start *VertexA
	End   *VertexB
	Test  string `gogm:"name=test"`
}

func (e *EdgeC) GetStartNode() interface{} {
	return e.Start
}

func (e *EdgeC) GetStartNodeType() reflect.Type {
	return reflect.TypeOf(&VertexA{})
}

func (e *EdgeC) SetStartNode(v interface{}) error {
	val, ok := v.(*VertexA)
	if !ok {
		return fmt.Errorf("unable to cast [%T] to *VertexA", v)
	}

	e.Start = val
	return nil
}

func (e *EdgeC) GetEndNode() interface{} {
	return e.End
}

func (e *EdgeC) GetEndNodeType() reflect.Type {
	return reflect.TypeOf(&VertexB{})
}

func (e *EdgeC) SetEndNode(v interface{}) error {
	val, ok := v.(*VertexB)
	if !ok {
		return fmt.Errorf("unable to cast [%T] to *VertexB", v)
	}

	e.End = val
	return nil
}

func main() {
	config := gogm.Config{
		IndexStrategy: gogm.VALIDATE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		PoolSize:      50,
		Port:          7687,
		IsCluster:     false, //tells it whether or not to use `bolt+routing`
		Host:          "10.35.19.9",
		Username:      "neo4j",
		Password:      "Un3jS85EjaARKMvu",
		//Host:     "120.27.23.246",
		//Username: "neo4j",
		//Password: "abcd1234",
	}

	_gogm, err := gogm.New(&config, gogm.UUIDPrimaryKeyStrategy,
		&domain.Device{},
		&domain.DeviceModel{},
		&domain.ServerCabinet{},
		&domain.ServerRoom{},
		&domain.CloudProvider{},
		&domain.DEVICE_SERVERCABINET{},
		&domain.DEVICE_DEVICEMODEL{},
		&domain.SERVERCABINET_SERVERROOM{},
		&domain.SERVERROOM_CLOUDPROVIDER{})
	if err != nil {
		panic(err)
	}

	//param is readonly, we're going to make stuff so we're going to do read write
	sess, err := _gogm.NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer sess.Close()

	cloudProvider := &domain.CloudProvider{
		Name: "gogm云提供商",
	}

	var serverRoomList []*domain.ServerRoom

	for i := 0; i < 100; i++ {
		serverRoom := &domain.ServerRoom{
			Address:           "机房地址",
			Area:              "区域",
			CommissioningDate: time.Now().UTC(),
			Maintainer:        "维护人员",
			Name:              "机房名称",
			Number:            "机房编号",
			Remark:            "备注",
			Telephone:         "联系电话",
		}
		serverRoom.CloudProvider = cloudProvider
		serverRoomList = append(serverRoomList, serverRoom)

		var serverCabinetList []*domain.ServerCabinet
		for j := 0; j < 50; j++ {
			serverCabinet := &domain.ServerCabinet{
				Capacity:      "机柜容量",
				ElectricPower: "电力",
				Maintainer:    "维护人员",
				Name:          "机柜名称",
				Number:        "机柜编号",
				StartDate:     time.Now().UTC(),
				EndDate:       time.Now().UTC(),
				Remark:        "备注",
			}
			serverCabinet.ServerRoom = serverRoom
			serverCabinetList = append(serverCabinetList, serverCabinet)

			var deviceList []*domain.Device
			for k := 0; k < 20; k++ {
				device := &domain.Device{
					AssetBelong:  "资产归属",
					AssetNumber:  "资产编号",
					RegisterTime: time.Now().UTC(),
					Sn:           "sn",
					UseType:      "服务器",
					Remark:       "备注",
				}
				device.ServerCabinet = serverCabinet
				deviceList = append(deviceList, device)
			}
			serverCabinet.Device = deviceList
		}
		serverRoom.ServerCabinet = serverCabinetList
	}
	cloudProvider.ServerRoom = serverRoomList

	err = sess.SaveDepth(context.Background(), cloudProvider, 4)
	if err != nil {
		panic(err)
	}

	//var cloudProviderQuery []*domain.CloudProvider
	//err = sess.LoadAll(context.Background(), &cloudProviderQuery)
	//if err != nil {
	//	panic(err)
	//}

	//for index, value := range cloudProviderQuery {
	//	list := *value
	//	for i, v := range list.ServerRoom {
	//		serverRoom := *v
	//		for j := 0; j < 50; j++ {
	//			serverCabinet := &domain.ServerCabinet{
	//				Maintainer: "维护人员",
	//				Name:       "机房名称",
	//				Number:     "机房编号",
	//				Remark:     "备注",
	//			}
	//			serverCabinet.ServerRoom = &serverRoom
	//			serverRoomList = append(serverRoomList, serverRoom)
	//		}
	//		cloudProvider.ServerRoom = serverRoomList
	//	}
	//}
}

func exp(config gogm.Config) {
	// register all vertices and edges
	// this is so that GoGM doesn't have to do reflect processing of each edge in real time
	// use nil or gogm.DefaultPrimaryKeyStrategy if you only want graph ids
	_gogm, err := gogm.New(&config, gogm.UUIDPrimaryKeyStrategy, &VertexA{}, &VertexB{}, &EdgeC{})
	if err != nil {
		panic(err)
	}

	//param is readonly, we're going to make stuff so we're going to do read write
	sess, err := _gogm.NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer sess.Close()

	aVal := &VertexA{
		TestField: "woo neo4j",
	}

	bVal := &VertexB{
		TestTime: time.Now().UTC(),
	}

	//set bi directional pointer
	bVal.Single = aVal
	aVal.SingleA = bVal

	err = sess.SaveDepth(context.Background(), aVal, 2)
	if err != nil {
		panic(err)
	}

	//load the object we just made (save will set the uuid)
	var readin VertexA
	err = sess.Load(context.Background(), &readin, aVal.UUID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", readin)
}
