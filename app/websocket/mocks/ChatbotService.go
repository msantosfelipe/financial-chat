// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ChatbotService is an autogenerated mock type for the ChatbotService type
type ChatbotService struct {
	mock.Mock
}

// HandleBotMessage provides a mock function with given fields: text, room
func (_m *ChatbotService) HandleBotMessage(text string, room string) {
	_m.Called(text, room)
}

// StockHandler provides a mock function with given fields: stock, room
func (_m *ChatbotService) StockHandler(stock string, room string) error {
	ret := _m.Called(stock, room)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(stock, room)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewChatbotService interface {
	mock.TestingT
	Cleanup(func())
}

// NewChatbotService creates a new instance of ChatbotService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewChatbotService(t mockConstructorTestingTNewChatbotService) *ChatbotService {
	mock := &ChatbotService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
