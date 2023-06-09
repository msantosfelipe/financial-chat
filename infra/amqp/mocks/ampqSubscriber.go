// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	amqp091 "github.com/rabbitmq/amqp091-go"
	mock "github.com/stretchr/testify/mock"
)

// ampqSubscriber is an autogenerated mock type for the ampqSubscriber type
type ampqSubscriber struct {
	mock.Mock
}

// SubscribeToQueue provides a mock function with given fields: queue
func (_m *ampqSubscriber) SubscribeToQueue(queue string) <-chan amqp091.Delivery {
	ret := _m.Called(queue)

	var r0 <-chan amqp091.Delivery
	if rf, ok := ret.Get(0).(func(string) <-chan amqp091.Delivery); ok {
		r0 = rf(queue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan amqp091.Delivery)
		}
	}

	return r0
}

type mockConstructorTestingTnewAmpqSubscriber interface {
	mock.TestingT
	Cleanup(func())
}

// newAmpqSubscriber creates a new instance of ampqSubscriber. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newAmpqSubscriber(t mockConstructorTestingTnewAmpqSubscriber) *ampqSubscriber {
	mock := &ampqSubscriber{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
