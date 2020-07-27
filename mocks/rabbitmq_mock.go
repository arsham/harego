// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	harego "github.com/blokur/harego"
	mock "github.com/stretchr/testify/mock"
)

// RabbitMQ is an autogenerated mock type for the RabbitMQ type
type RabbitMQ struct {
	mock.Mock
}

// Channel provides a mock function with given fields:
func (_m *RabbitMQ) Channel() (harego.Channel, error) {
	ret := _m.Called()

	var r0 harego.Channel
	if rf, ok := ret.Get(0).(func() harego.Channel); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(harego.Channel)
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

// Close provides a mock function with given fields:
func (_m *RabbitMQ) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
