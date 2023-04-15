// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// chatbotHandler is an autogenerated mock type for the chatbotHandler type
type chatbotHandler struct {
	mock.Mock
}

// HandleBotMessage provides a mock function with given fields: text, room
func (_m *chatbotHandler) HandleBotMessage(text string, room string) {
	_m.Called(text, room)
}

// StockHandler provides a mock function with given fields: stock, room
func (_m *chatbotHandler) StockHandler(stock string, room string) error {
	ret := _m.Called(stock, room)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(stock, room)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewChatbotHandler interface {
	mock.TestingT
	Cleanup(func())
}

// newChatbotHandler creates a new instance of chatbotHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newChatbotHandler(t mockConstructorTestingTnewChatbotHandler) *chatbotHandler {
	mock := &chatbotHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}