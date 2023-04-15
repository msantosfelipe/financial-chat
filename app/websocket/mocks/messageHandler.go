// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// messageHandler is an autogenerated mock type for the messageHandler type
type messageHandler struct {
	mock.Mock
}

// HandleBotMessage provides a mock function with given fields: text, room
func (_m *messageHandler) HandleBotMessage(text string, room string) {
	_m.Called(text, room)
}

// StockHandler provides a mock function with given fields: stock, room
func (_m *messageHandler) StockHandler(stock string, room string) error {
	ret := _m.Called(stock, room)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(stock, room)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMessageHandler interface {
	mock.TestingT
	Cleanup(func())
}

// newMessageHandler creates a new instance of messageHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMessageHandler(t mockConstructorTestingTnewMessageHandler) *messageHandler {
	mock := &messageHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
