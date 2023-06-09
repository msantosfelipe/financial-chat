// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ServiceCleaner is an autogenerated mock type for the ServiceCleaner type
type ServiceCleaner struct {
	mock.Mock
}

// Clean provides a mock function with given fields:
func (_m *ServiceCleaner) Clean() {
	_m.Called()
}

type mockConstructorTestingTNewServiceCleaner interface {
	mock.TestingT
	Cleanup(func())
}

// NewServiceCleaner creates a new instance of ServiceCleaner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewServiceCleaner(t mockConstructorTestingTNewServiceCleaner) *ServiceCleaner {
	mock := &ServiceCleaner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
