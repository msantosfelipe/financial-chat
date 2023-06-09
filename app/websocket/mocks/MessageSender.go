// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	websocket "github.com/gorilla/websocket"
	mock "github.com/stretchr/testify/mock"
)

// MessageSender is an autogenerated mock type for the MessageSender type
type MessageSender struct {
	mock.Mock
}

// ListenAndSendMessage provides a mock function with given fields: wsConn, room
func (_m *MessageSender) ListenAndSendMessage(wsConn *websocket.Conn, room string) {
	_m.Called(wsConn, room)
}

// PublishMessageToQueue provides a mock function with given fields: msg, queue
func (_m *MessageSender) PublishMessageToQueue(msg []byte, queue string) error {
	ret := _m.Called(msg, queue)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, string) error); ok {
		r0 = rf(msg, queue)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendBotMessage provides a mock function with given fields: room, text
func (_m *MessageSender) SendBotMessage(room string, text string) {
	_m.Called(room, text)
}

// SendMessage provides a mock function with given fields: user, room, text
func (_m *MessageSender) SendMessage(user string, room string, text string) {
	_m.Called(user, room, text)
}

// SendPreviousCachedMessages provides a mock function with given fields: wsConn, room
func (_m *MessageSender) SendPreviousCachedMessages(wsConn *websocket.Conn, room string) {
	_m.Called(wsConn, room)
}

type mockConstructorTestingTNewMessageSender interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageSender creates a new instance of MessageSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageSender(t mockConstructorTestingTNewMessageSender) *MessageSender {
	mock := &MessageSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
