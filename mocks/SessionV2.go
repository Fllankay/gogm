// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	go_cypherdsl "github.com/mindstand/go-cypherdsl"
	gogm "github.com/mindstand/gogm/v2"

	mock "github.com/stretchr/testify/mock"

	neo4j "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// SessionV2 is an autogenerated mock type for the SessionV2 type
type SessionV2 struct {
	mock.Mock
}

// Begin provides a mock function with given fields: ctx
func (_m *SessionV2) Begin(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *SessionV2) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Commit provides a mock function with given fields: ctx
func (_m *SessionV2) Commit(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, deleteObj
func (_m *SessionV2) Delete(ctx context.Context, deleteObj interface{}) error {
	ret := _m.Called(ctx, deleteObj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, deleteObj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUUID provides a mock function with given fields: ctx, uuid
func (_m *SessionV2) DeleteUUID(ctx context.Context, uuid string) error {
	ret := _m.Called(ctx, uuid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Load provides a mock function with given fields: ctx, respObj, id
func (_m *SessionV2) Load(ctx context.Context, respObj interface{}, id interface{}) error {
	ret := _m.Called(ctx, respObj, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}) error); ok {
		r0 = rf(ctx, respObj, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadAll provides a mock function with given fields: ctx, respObj
func (_m *SessionV2) LoadAll(ctx context.Context, respObj interface{}) error {
	ret := _m.Called(ctx, respObj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, respObj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadAllDepth provides a mock function with given fields: ctx, respObj, depth
func (_m *SessionV2) LoadAllDepth(ctx context.Context, respObj interface{}, depth int) error {
	ret := _m.Called(ctx, respObj, depth)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, int) error); ok {
		r0 = rf(ctx, respObj, depth)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadAllDepthFilter provides a mock function with given fields: ctx, respObj, depth, filter, params
func (_m *SessionV2) LoadAllDepthFilter(ctx context.Context, respObj interface{}, depth int, filter go_cypherdsl.ConditionOperator, params map[string]interface{}) error {
	ret := _m.Called(ctx, respObj, depth, filter, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, int, go_cypherdsl.ConditionOperator, map[string]interface{}) error); ok {
		r0 = rf(ctx, respObj, depth, filter, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadAllDepthFilterPagination provides a mock function with given fields: ctx, respObj, depth, filter, params, pagination
func (_m *SessionV2) LoadAllDepthFilterPagination(ctx context.Context, respObj interface{}, depth int, filter go_cypherdsl.ConditionOperator, params map[string]interface{}, pagination *gogm.Pagination) error {
	ret := _m.Called(ctx, respObj, depth, filter, params, pagination)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, int, go_cypherdsl.ConditionOperator, map[string]interface{}, *gogm.Pagination) error); ok {
		r0 = rf(ctx, respObj, depth, filter, params, pagination)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadDepth provides a mock function with given fields: ctx, respObj, id, depth
func (_m *SessionV2) LoadDepth(ctx context.Context, respObj interface{}, id interface{}, depth int) error {
	ret := _m.Called(ctx, respObj, id, depth)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, int) error); ok {
		r0 = rf(ctx, respObj, id, depth)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadDepthFilter provides a mock function with given fields: ctx, respObj, id, depth, filter, params
func (_m *SessionV2) LoadDepthFilter(ctx context.Context, respObj interface{}, id interface{}, depth int, filter go_cypherdsl.ConditionOperator, params map[string]interface{}) error {
	ret := _m.Called(ctx, respObj, id, depth, filter, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, int, go_cypherdsl.ConditionOperator, map[string]interface{}) error); ok {
		r0 = rf(ctx, respObj, id, depth, filter, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadDepthFilterPagination provides a mock function with given fields: ctx, respObj, id, depth, filter, params, pagination
func (_m *SessionV2) LoadDepthFilterPagination(ctx context.Context, respObj interface{}, id interface{}, depth int, filter go_cypherdsl.ConditionOperator, params map[string]interface{}, pagination *gogm.Pagination) error {
	ret := _m.Called(ctx, respObj, id, depth, filter, params, pagination)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, int, go_cypherdsl.ConditionOperator, map[string]interface{}, *gogm.Pagination) error); ok {
		r0 = rf(ctx, respObj, id, depth, filter, params, pagination)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManagedTransaction provides a mock function with given fields: ctx, work
func (_m *SessionV2) ManagedTransaction(ctx context.Context, work gogm.TransactionWork) error {
	ret := _m.Called(ctx, work)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, gogm.TransactionWork) error); ok {
		r0 = rf(ctx, work)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Query provides a mock function with given fields: ctx, query, properties, respObj
func (_m *SessionV2) Query(ctx context.Context, query string, properties map[string]interface{}, respObj interface{}) error {
	ret := _m.Called(ctx, query, properties, respObj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}, interface{}) error); ok {
		r0 = rf(ctx, query, properties, respObj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryRaw provides a mock function with given fields: ctx, query, properties
func (_m *SessionV2) QueryRaw(ctx context.Context, query string, properties map[string]interface{}) ([][]interface{}, neo4j.ResultSummary, error) {
	ret := _m.Called(ctx, query, properties)

	var r0 [][]interface{}
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) [][]interface{}); ok {
		r0 = rf(ctx, query, properties)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]interface{})
		}
	}

	var r1 neo4j.ResultSummary
	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]interface{}) neo4j.ResultSummary); ok {
		r1 = rf(ctx, query, properties)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(neo4j.ResultSummary)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, map[string]interface{}) error); ok {
		r2 = rf(ctx, query, properties)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Rollback provides a mock function with given fields: ctx
func (_m *SessionV2) Rollback(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RollbackWithError provides a mock function with given fields: ctx, err
func (_m *SessionV2) RollbackWithError(ctx context.Context, err error) error {
	ret := _m.Called(ctx, err)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, error) error); ok {
		r0 = rf(ctx, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: ctx, saveObj
func (_m *SessionV2) Save(ctx context.Context, saveObj interface{}) error {
	ret := _m.Called(ctx, saveObj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, saveObj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveDepth provides a mock function with given fields: ctx, saveObj, depth
func (_m *SessionV2) SaveDepth(ctx context.Context, saveObj interface{}, depth int) error {
	ret := _m.Called(ctx, saveObj, depth)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, int) error); ok {
		r0 = rf(ctx, saveObj, depth)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
