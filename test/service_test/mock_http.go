package main

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHTTP is a mock of ScalingTarget interface
type MockHTTP struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPMockRecorder
}

// MockHTTPMockRecorder is the mock recorder for MockHTTP
type MockHTTPMockRecorder struct {
	mock *MockHTTP
}

// NewMockHTTP creates a new mock instance
func NewMockHTTP(ctrl *gomock.Controller) *MockHTTP {
	mock := &MockHTTP{ctrl: ctrl}
	mock.recorder = &MockHTTPMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHTTP) EXPECT() *MockHTTPMockRecorder {
	return m.recorder
}

// POST mocks base method
func (m *MockHTTP) POST(data string) (int, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "POST", data)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// POST indicates an expected call of POST
func (mr *MockHTTPMockRecorder) POST(data string) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "POST", reflect.TypeOf((*MockHTTP)(nil).POST), data)
}
