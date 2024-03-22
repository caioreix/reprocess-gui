// Code generated by mockery v2.42.1. DO NOT EDIT.

package portmock

import (
	context "context"
	domain "reprocess-gui/internal/apps/api/core/domain"

	mock "github.com/stretchr/testify/mock"
)

// TableService is an autogenerated mock type for the TableService type
type TableService struct {
	mock.Mock
}

// GetAllTables provides a mock function with given fields: ctx
func (_m *TableService) GetAllTables(ctx context.Context) ([]*domain.Table, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllTables")
	}

	var r0 []*domain.Table
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*domain.Table, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*domain.Table); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Table)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTableService creates a new instance of TableService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTableService(t interface {
	mock.TestingT
	Cleanup(func())
}) *TableService {
	mock := &TableService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
