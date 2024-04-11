// Code generated by mockery v2.42.2. DO NOT EDIT.

package portmock

import (
	context "context"
	domain "reprocess-gui/internal/apps/api/core/domain"

	mock "github.com/stretchr/testify/mock"

	utils "reprocess-gui/internal/utils"
)

// ConsumerRepository is an autogenerated mock type for the ConsumerRepository type
type ConsumerRepository struct {
	mock.Mock
}

// GetAllConsumers provides a mock function with given fields: ctx, pageToken
func (_m *ConsumerRepository) GetAllConsumers(ctx context.Context, pageToken *utils.PaginationToken) ([]*domain.Consumer, error) {
	ret := _m.Called(ctx, pageToken)

	if len(ret) == 0 {
		panic("no return value specified for GetAllConsumers")
	}

	var r0 []*domain.Consumer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *utils.PaginationToken) ([]*domain.Consumer, error)); ok {
		return rf(ctx, pageToken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *utils.PaginationToken) []*domain.Consumer); ok {
		r0 = rf(ctx, pageToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Consumer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *utils.PaginationToken) error); ok {
		r1 = rf(ctx, pageToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertNewConsumer provides a mock function with given fields: ctx, consumer
func (_m *ConsumerRepository) InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error) {
	ret := _m.Called(ctx, consumer)

	if len(ret) == 0 {
		panic("no return value specified for InsertNewConsumer")
	}

	var r0 *domain.Consumer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Consumer) (*domain.Consumer, error)); ok {
		return rf(ctx, consumer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Consumer) *domain.Consumer); ok {
		r0 = rf(ctx, consumer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Consumer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Consumer) error); ok {
		r1 = rf(ctx, consumer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewConsumerRepository creates a new instance of ConsumerRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConsumerRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConsumerRepository {
	mock := &ConsumerRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
