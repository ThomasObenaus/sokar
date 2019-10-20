package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	reflect "reflect"
	"sync"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/thomasobenaus/sokar/api"
)

const defaultTimeout = time.Second * 20

// MockHTTP is a mock of ScalingTarget interface
type MockHTTP struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPMockRecorder
	timeout  time.Duration
}

// MockHTTPMockRecorder is the mock recorder for MockHTTP
type MockHTTPMockRecorder struct {
	mock   *MockHTTP
	server *api.API
	wg     sync.WaitGroup

	registeredPOSTPaths map[string]struct{}
	registeredGETPaths  map[string]struct{}
}

// NewMockHTTP creates a new mock instance (timeout/ deadline is 20s)
// Pattern:
// mock := NewMockHTTP(t, 18000)
// defer mock.Finish()
// mock.EXPECT().GET("/path").Return(http.StatusOK, "Someting")
func NewMockHTTP(t *testing.T, port int) *MockHTTP {

	mockCtrl := gomock.NewController(t)

	server := api.New(port)
	server.Run()

	mock := &MockHTTP{ctrl: mockCtrl, timeout: defaultTimeout}
	mock.recorder = &MockHTTPMockRecorder{
		mock:                mock,
		server:              server,
		registeredPOSTPaths: make(map[string]struct{}, 0),
		registeredGETPaths:  make(map[string]struct{}, 0),
	}

	// Install a handler for all resources that are not expected to be called
	mock.recorder.server.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.GET(r.URL.Path)

		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Unexpected call to this resource")
	})
	// Disable the method not allowed handler to be able to catch all unexpected calls to resources
	mock.recorder.server.Router.HandleMethodNotAllowed = false

	return mock
}

// WithTimeout creates a new mock instance wit the specified timeout.
// This means if not all specified EXPECT calls are made within the defined timeout
// the test will fail.
// Pattern:
// mock := WithTimeout(t, 18000,time.Second*5)
// defer mock.Finish()
// mock.EXPECT().GET("/path").Return(http.StatusOK, "Someting")
func WithTimeout(t *testing.T, port int, timeout time.Duration) *MockHTTP {
	mock := NewMockHTTP(t, port)
	mock.timeout = timeout
	return mock
}

// Finish has to be called at the end to clean up and to check if all expected calls where made.
func (m *MockHTTP) Finish() {

	// wait until the wait group was released (all expected calls where made)
	// or until the deadline/ timeout was exceeded
	wait(&m.recorder.wg, m.timeout)

	// clean up
	m.recorder.server.Stop()
	m.ctrl.Finish()
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHTTP) EXPECT() *MockHTTPMockRecorder {
	return m.recorder
}

// POST mocks base method
func (m *MockHTTP) POST(path, data string) (int, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "POST", path, data)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// POST indicates an expected call of POST
func (mr *MockHTTPMockRecorder) POST(path, data string) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	mr.wg.Add(1)

	// Register the http handler, but only if it is not already registered for this path
	_, pathAlreadyRegistered := mr.registeredPOSTPaths[path]
	if !pathAlreadyRegistered {
		mr.registeredPOSTPaths[path] = struct{}{}

		mr.server.Router.HandlerFunc("POST", path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer mr.wg.Done()

			if r == nil {
				http.Error(w, "Request is nil", http.StatusInternalServerError)
				return
			}

			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			code, data := mr.mock.POST(path, string(body))
			w.WriteHeader(code)
			io.WriteString(w, data)
		}))
	}

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "POST", reflect.TypeOf((*MockHTTP)(nil).POST), path, data)
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
	mr.wg.Add(1)

	// Register the http handler, but only if it is not already registered for this path
	_, pathAlreadyRegistered := mr.registeredGETPaths[path]
	if !pathAlreadyRegistered {
		mr.registeredGETPaths[path] = struct{}{}

		mr.server.Router.HandlerFunc("GET", path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer mr.wg.Done()

			if r == nil {
				http.Error(w, "Request is nil", http.StatusInternalServerError)
				return
			}

			code, data := mr.mock.GET(path)
			w.WriteHeader(code)
			io.WriteString(w, data)
			fmt.Printf("Request: %v\n", r)
		}))
	}

	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GET", reflect.TypeOf((*MockHTTP)(nil).GET), path)
}
