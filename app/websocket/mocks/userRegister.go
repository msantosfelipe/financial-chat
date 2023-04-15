// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	websocket "github.com/gorilla/websocket"
	mock "github.com/stretchr/testify/mock"
)

// userRegister is an autogenerated mock type for the userRegister type
type userRegister struct {
	mock.Mock
}

// AddUserToRoom provides a mock function with given fields: wsConn, room
func (_m *userRegister) AddUserToRoom(wsConn *websocket.Conn, room string) error {
	ret := _m.Called(wsConn, room)

	var r0 error
	if rf, ok := ret.Get(0).(func(*websocket.Conn, string) error); ok {
		r0 = rf(wsConn, room)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewUserRegister interface {
	mock.TestingT
	Cleanup(func())
}

// newUserRegister creates a new instance of userRegister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newUserRegister(t mockConstructorTestingTnewUserRegister) *userRegister {
	mock := &userRegister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
