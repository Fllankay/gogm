package gogm

import (
	dsl "github.com/mindstand/go-cypherdsl"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

//session version 2 is experimental to start trying breaking changes
type SessionV2 interface {
	//transaction functions
	TransactionV2

	// Begin begins transaction
	Begin() error

	// ManagedTransaction runs tx work managed for retry
	ManagedTransaction(work TransactionWork) error

	// closes session
	Close() error
}

// TransactionV2 specifies functions for Neo4j ACID transactions
type TransactionV2 interface {
	// Rollback rolls back transaction
	Rollback() error
	// RollbackWithError wraps original error into rollback error if there is one
	RollbackWithError(err error) error
	// Commit commits transaction
	Commit() error

	// functions the tx can do
	ogmFunctions
}

type ogmFunctions interface {
	//load single object
	Load(respObj interface{}, id string) error

	//load object with depth
	LoadDepth(respObj interface{}, id string, depth int) error

	//load with depth and filter
	LoadDepthFilter(respObj interface{}, id string, depth int, filter *dsl.ConditionBuilder, params map[string]interface{}) error

	//load with depth, filter and pagination
	LoadDepthFilterPagination(respObj interface{}, id string, depth int, filter dsl.ConditionOperator, params map[string]interface{}, pagination *Pagination) error

	//load slice of something
	LoadAll(respObj interface{}) error

	//load all of depth
	LoadAllDepth(respObj interface{}, depth int) error

	//load all of type with depth and filter
	LoadAllDepthFilter(respObj interface{}, depth int, filter dsl.ConditionOperator, params map[string]interface{}) error

	//load all with depth, filter and pagination
	LoadAllDepthFilterPagination(respObj interface{}, depth int, filter dsl.ConditionOperator, params map[string]interface{}, pagination *Pagination) error

	//load all edge query
	LoadAllEdgeConstraint(respObj interface{}, endNodeType, endNodeField string, edgeConstraint interface{}, minJumps, maxJumps, depth int, filter dsl.ConditionOperator) error

	//save object
	Save(saveObj interface{}) error

	//save object with depth
	SaveDepth(saveObj interface{}, depth int) error

	//delete
	Delete(deleteObj interface{}) error

	//delete uuid
	DeleteUUID(uuid string) error

	//specific query, responds to slice and single objects
	Query(query string, properties map[string]interface{}, respObj interface{}) error

	//similar to query, but returns raw rows/cols
	QueryRaw(query string, properties map[string]interface{}) ([][]interface{}, neo4j.ResultSummary, error)

	//delete everything, this will literally delete everything
	PurgeDatabase() error
}

type TransactionWork func(tx TransactionV2) error
