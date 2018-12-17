// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import api "chainspace.io/prototype/pubsub/api"
import http "net/http"
import mock "github.com/stretchr/testify/mock"

// WSUpgrader is an autogenerated mock type for the WSUpgrader type
type WSUpgrader struct {
	mock.Mock
}

// Upgrade provides a mock function with given fields: w, r
func (_m *WSUpgrader) Upgrade(w http.ResponseWriter, r *http.Request) (api.WriteMessageCloser, error) {
	ret := _m.Called(w, r)

	var r0 api.WriteMessageCloser
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) api.WriteMessageCloser); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.WriteMessageCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) error); ok {
		r1 = rf(w, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}