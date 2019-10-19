package main

import (
	"io"
	"net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"github.com/thomasobenaus/sokar/api"
)

// MockHTTP is a mock of ScalingTarget interface
type MockHTTP struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPMockRecorder
}

// MockHTTPMockRecorder is the mock recorder for MockHTTP
type MockHTTPMockRecorder struct {
	mock     *MockHTTP
	receiver *api.API
}

// NewMockHTTP creates a new mock instance
func NewMockHTTP(ctrl *gomock.Controller, port int) *MockHTTP {
	receiver := api.New(port)
	receiver.Run()
	mock := &MockHTTP{ctrl: ctrl}
	mock.recorder = &MockHTTPMockRecorder{mock: mock, receiver: receiver}

	// TODO: How/ when to stop the api server
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

// GET mocks base method
func (m *MockHTTP) GET(path string) (int, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GET", path)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// GET indicates an expected call of GET
func (mr *MockHTTPMockRecorder) GET(path string) *gomock.Call {
	mr.mock.ctrl.T.Helper()

	mr.receiver.Router.HandlerFunc("GET", path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, data := mr.mock.GET(path)
		w.WriteHeader(code)
		io.WriteString(w, data)
	}))

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GET", reflect.TypeOf((*MockHTTP)(nil).GET), path)
}
