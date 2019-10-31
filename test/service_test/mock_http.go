package main

import (
	"fmt"
	"io"
	"net/http"
	reflect "reflect"

	"github.com/thomasobenaus/sokar/api"
)

// MockHTTP is a mock of ScalingTarget interface
type MockHTTP struct {
	ctrl     *Controller
	recorder *MockHTTPMockRecorder

	// failOnUnexpectedCalls set if the test should fail in case a end-point is called which is not covered by an EXPECT call.
	failOnUnexpectedCalls bool
}

// MockHTTPMockRecorder is the mock recorder for MockHTTP
type MockHTTPMockRecorder struct {
	mock   *MockHTTP
	server *api.API

	registeredPOSTPaths map[string]struct{}
	registeredGETPaths  map[string]struct{}
}

// Option represents an option for the MockHTTP
type Option func(m *MockHTTP)

// FailOnUnexpectedCalls set if the test should fail in case a end-point is called which is not covered by an EXPECT call.
func FailOnUnexpectedCalls(fail bool) Option {
	return func(m *MockHTTP) {
		m.failOnUnexpectedCalls = fail
	}
}

// NewMockHTTP creates a new mock instance (timeout/ deadline is 20s)
// Pattern:
// mock := NewMockHTTP(t, 18000)
// defer mock.Finish()
// mock.EXPECT().GET("/path").Return(http.StatusOK, "Someting")
func NewMockHTTP(c *Controller, port int, options ...Option) *MockHTTP {

	server := api.New(port)
	server.Run()

	mock := &MockHTTP{
		ctrl:                  c,
		failOnUnexpectedCalls: false,
	}
	mock.recorder = &MockHTTPMockRecorder{
		mock:                mock,
		server:              server,
		registeredPOSTPaths: make(map[string]struct{}, 0),
		registeredGETPaths:  make(map[string]struct{}, 0),
	}

	// apply the options
	for _, opt := range options {
		opt(mock)
	}

	if mock.failOnUnexpectedCalls {
		// Install a handler for all resources that are not expected to be called
		mock.recorder.server.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Unexpected call to this resource")

			mock.ctrl.releaseAllCallLocks()

			if r == nil {
				mock.ctrl.gmckCtrl.T.Fatalf("Received request object is nil.")
			}

			path := r.URL.Path

			switch r.Method {
			case http.MethodGet:
				mock.GET(path)
				break
			default:
				mock.ctrl.gmckCtrl.T.Fatalf("HTTP method %s is not yet supported (Request was %s: %s%s).", r.Method, r.Method, r.Host, r.URL)
				break
			}
		})
		// Disable the method not allowed handler to be able to catch all unexpected calls to resources
		mock.recorder.server.Router.HandleMethodNotAllowed = false
	}

	// register the stop method of the router to be called after the controller has been finished
	c.addCleanupCallback(mock.recorder.server.Stop)

	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHTTP) EXPECT() *MockHTTPMockRecorder {
	return m.recorder
}

//// POST mocks base method
//func (m *MockHTTP) POST(path, data string) (int, string) {
//	m.ctrl.T.Helper()
//	ret := m.ctrl.Call(m, "POST", path, data)
//	ret0, _ := ret[0].(int)
//	ret1, _ := ret[1].(string)
//	return ret0, ret1
//}
//
//// POST indicates an expected call of POST
//func (mr *MockHTTPMockRecorder) POST(path, data string) Call {
//	mr.mock.ctrl.T.Helper()
//	mr.wg.Add(1)
//
//	// Register the http handler, but only if it is not already registered for this path
//	_, pathAlreadyRegistered := mr.registeredPOSTPaths[path]
//	if !pathAlreadyRegistered {
//		mr.registeredPOSTPaths[path] = struct{}{}
//
//		mr.server.Router.HandlerFunc("POST", path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			defer mr.wg.Done()
//
//			if r == nil {
//				http.Error(w, "Request is nil", http.StatusInternalServerError)
//				return
//			}
//
//			defer r.Body.Close()
//			body, err := ioutil.ReadAll(r.Body)
//			if err != nil {
//				http.Error(w, err.Error(), http.StatusBadRequest)
//				return
//			}
//
//			code, data := mr.mock.POST(path, string(body))
//			w.WriteHeader(code)
//			io.WriteString(w, data)
//		}))
//	}
//
//	call := callImpl{
//		gomockCall: mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "POST", reflect.TypeOf((*MockHTTP)(nil).POST), path, data),
//	}
//
//	return &call
//}

// GET mocks base method
func (m *MockHTTP) GET(path string) Response {
	m.ctrl.gmckCtrl.T.Helper()
	ret := m.ctrl.gmckCtrl.Call(m, "GET", path)
	ret0, _ := ret[0].(Response)
	return ret0
}

// GET indicates an expected call of GET
func (mr *MockHTTPMockRecorder) GET(path string) Call {
	mr.mock.ctrl.gmckCtrl.T.Helper()

	// Register the http handler, but only if it is not already registered for this path
	_, pathAlreadyRegistered := mr.registeredGETPaths[path]
	if !pathAlreadyRegistered {
		mr.registeredGETPaths[path] = struct{}{}

		mr.server.Router.HandlerFunc("GET", path, mr.mock.handleRequest)
	}

	gomockCall := mr.mock.ctrl.gmckCtrl.RecordCallWithMethodType(mr.mock, "GET", reflect.TypeOf((*MockHTTP)(nil).GET), path)
	return mr.mock.ctrl.CreateAndRegisterGETCall(gomockCall)
}

func (m *MockHTTP) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		http.Error(w, "Request is nil", http.StatusInternalServerError)
		return
	}
	if r.URL == nil {
		http.Error(w, "Request.URL is nil", http.StatusInternalServerError)
		return
	}

	path := r.URL.Path

	var response Response
	if r.Method == http.MethodGet {
		response = m.GET(path)
	} else {
		panic(fmt.Sprintf("HTTP Method '%s' not implemented yet.", r.Method))
	}

	// fill the response (header, data and status code)
	for key, valueList := range response.header {
		for _, value := range valueList {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(response.statusCode)
	w.Write(response.data)
}
