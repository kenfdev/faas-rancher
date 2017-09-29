// Code generated by mockery v1.0.0
package mocks

import client "github.com/rancher/go-rancher/v3"
import mock "github.com/stretchr/testify/mock"

// BridgeClient is an autogenerated mock type for the BridgeClient type
type BridgeClient struct {
	mock.Mock
}

// CreateService provides a mock function with given fields: spec
func (_m *BridgeClient) CreateService(spec *client.Service) (*client.Service, error) {
	ret := _m.Called(spec)

	var r0 *client.Service
	if rf, ok := ret.Get(0).(func(*client.Service) *client.Service); ok {
		r0 = rf(spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*client.Service) error); ok {
		r1 = rf(spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteService provides a mock function with given fields: spec
func (_m *BridgeClient) DeleteService(spec *client.Service) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*client.Service) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindServiceByName provides a mock function with given fields: name
func (_m *BridgeClient) FindServiceByName(name string) (*client.Service, error) {
	ret := _m.Called(name)

	var r0 *client.Service
	if rf, ok := ret.Get(0).(func(string) *client.Service); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListServices provides a mock function with given fields:
func (_m *BridgeClient) ListServices() ([]client.Service, error) {
	ret := _m.Called()

	var r0 []client.Service
	if rf, ok := ret.Get(0).(func() []client.Service); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]client.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateService provides a mock function with given fields: spec, updates
func (_m *BridgeClient) UpdateService(spec *client.Service, updates map[string]string) (*client.Service, error) {
	ret := _m.Called(spec, updates)

	var r0 *client.Service
	if rf, ok := ret.Get(0).(func(*client.Service, map[string]string) *client.Service); ok {
		r0 = rf(spec, updates)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*client.Service, map[string]string) error); ok {
		r1 = rf(spec, updates)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
