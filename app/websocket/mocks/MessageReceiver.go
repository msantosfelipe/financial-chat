// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MessageReceiver is an autogenerated mock type for the MessageReceiver type
type MessageReceiver struct {
	mock.Mock
}

// HandleReceivedMessages provides a mock function with given fields:
func (_m *MessageReceiver) HandleReceivedMessages() {
	_m.Called()
}

type mockConstructorTestingTNewMessageReceiver interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageReceiver creates a new instance of MessageReceiver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageReceiver(t mockConstructorTestingTNewMessageReceiver) *MessageReceiver {
	mock := &MessageReceiver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
