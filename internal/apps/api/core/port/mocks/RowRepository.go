// Code generated by mockery v2.42.1. DO NOT EDIT.

package portmock

import (
	context "context"
	domain "reprocess-gui/internal/apps/api/core/domain"

	mock "github.com/stretchr/testify/mock"
)

// RowRepository is an autogenerated mock type for the RowRepository type
type RowRepository struct {
	mock.Mock
}

// InsertNewError provides a mock function with given fields: ctx, row
func (_m *RowRepository) InsertNewError(ctx context.Context, row *domain.Row) (*domain.Row, error) {
	ret := _m.Called(ctx, row)

	if len(ret) == 0 {
		panic("no return value specified for InsertNewError")
	}

	var r0 *domain.Row
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Row) (*domain.Row, error)); ok {
		return rf(ctx, row)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Row) *domain.Row); ok {
		r0 = rf(ctx, row)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Row)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Row) error); ok {
		r1 = rf(ctx, row)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRowRepository creates a new instance of RowRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRowRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *RowRepository {
	mock := &RowRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}