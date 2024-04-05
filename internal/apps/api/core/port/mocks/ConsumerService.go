// Code generated by mockery v2.42.1. DO NOT EDIT.

package portmock

import (
	context "context"
	domain "reprocess-gui/internal/apps/api/core/domain"

	mock "github.com/stretchr/testify/mock"
)

// ConsumerService is an autogenerated mock type for the ConsumerService type
type ConsumerService struct {
	mock.Mock
}

// InsertNewConsumer provides a mock function with given fields: ctx, consumer
func (_m *ConsumerService) InsertNewConsumer(ctx context.Context, consumer *domain.Consumer) (*domain.Consumer, error) {
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

// NewConsumerService creates a new instance of ConsumerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConsumerService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConsumerService {
	mock := &ConsumerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}