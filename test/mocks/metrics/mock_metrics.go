// Code generated by MockGen. DO NOT EDIT.
// Source: metrics/metrics.go

// Package mock_metrics is a generated GoMock package.
package mock_metrics

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	metrics "github.com/thomasobenaus/sokar/metrics"
)

// MockCounter is a mock of Counter interface.
type MockCounter struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMockRecorder
}

// MockCounterMockRecorder is the mock recorder for MockCounter.
type MockCounterMockRecorder struct {
	mock *MockCounter
}

// NewMockCounter creates a new mock instance.
func NewMockCounter(ctrl *gomock.Controller) *MockCounter {
	mock := &MockCounter{ctrl: ctrl}
	mock.recorder = &MockCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCounter) EXPECT() *MockCounterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCounter) Add(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add.
func (mr *MockCounterMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCounter)(nil).Add), arg0)
}

// Inc mocks base method.
func (m *MockCounter) Inc() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc.
func (mr *MockCounterMockRecorder) Inc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockCounter)(nil).Inc))
}

// MockGauge is a mock of Gauge interface.
type MockGauge struct {
	ctrl     *gomock.Controller
	recorder *MockGaugeMockRecorder
}

// MockGaugeMockRecorder is the mock recorder for MockGauge.
type MockGaugeMockRecorder struct {
	mock *MockGauge
}

// NewMockGauge creates a new mock instance.
func NewMockGauge(ctrl *gomock.Controller) *MockGauge {
	mock := &MockGauge{ctrl: ctrl}
	mock.recorder = &MockGaugeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGauge) EXPECT() *MockGaugeMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockGauge) Add(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add.
func (mr *MockGaugeMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockGauge)(nil).Add), arg0)
}

// Set mocks base method.
func (m *MockGauge) Set(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0)
}

// Set indicates an expected call of Set.
func (mr *MockGaugeMockRecorder) Set(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockGauge)(nil).Set), arg0)
}

// MockGaugeVec is a mock of GaugeVec interface.
type MockGaugeVec struct {
	ctrl     *gomock.Controller
	recorder *MockGaugeVecMockRecorder
}

// MockGaugeVecMockRecorder is the mock recorder for MockGaugeVec.
type MockGaugeVecMockRecorder struct {
	mock *MockGaugeVec
}

// NewMockGaugeVec creates a new mock instance.
func NewMockGaugeVec(ctrl *gomock.Controller) *MockGaugeVec {
	mock := &MockGaugeVec{ctrl: ctrl}
	mock.recorder = &MockGaugeVecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGaugeVec) EXPECT() *MockGaugeVecMockRecorder {
	return m.recorder
}

// WithLabelValues mocks base method.
func (m *MockGaugeVec) WithLabelValues(lvs ...string) metrics.Gauge {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range lvs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WithLabelValues", varargs...)
	ret0, _ := ret[0].(metrics.Gauge)
	return ret0
}

// WithLabelValues indicates an expected call of WithLabelValues.
func (mr *MockGaugeVecMockRecorder) WithLabelValues(lvs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithLabelValues", reflect.TypeOf((*MockGaugeVec)(nil).WithLabelValues), lvs...)
}

// MockCounterVec is a mock of CounterVec interface.
type MockCounterVec struct {
	ctrl     *gomock.Controller
	recorder *MockCounterVecMockRecorder
}

// MockCounterVecMockRecorder is the mock recorder for MockCounterVec.
type MockCounterVecMockRecorder struct {
	mock *MockCounterVec
}

// NewMockCounterVec creates a new mock instance.
func NewMockCounterVec(ctrl *gomock.Controller) *MockCounterVec {
	mock := &MockCounterVec{ctrl: ctrl}
	mock.recorder = &MockCounterVecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCounterVec) EXPECT() *MockCounterVecMockRecorder {
	return m.recorder
}

// WithLabelValues mocks base method.
func (m *MockCounterVec) WithLabelValues(lvs ...string) metrics.Counter {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range lvs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WithLabelValues", varargs...)
	ret0, _ := ret[0].(metrics.Counter)
	return ret0
}

// WithLabelValues indicates an expected call of WithLabelValues.
func (mr *MockCounterVecMockRecorder) WithLabelValues(lvs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithLabelValues", reflect.TypeOf((*MockCounterVec)(nil).WithLabelValues), lvs...)
}

// MockHistogram is a mock of Histogram interface.
type MockHistogram struct {
	ctrl     *gomock.Controller
	recorder *MockHistogramMockRecorder
}

// MockHistogramMockRecorder is the mock recorder for MockHistogram.
type MockHistogramMockRecorder struct {
	mock *MockHistogram
}

// NewMockHistogram creates a new mock instance.
func NewMockHistogram(ctrl *gomock.Controller) *MockHistogram {
	mock := &MockHistogram{ctrl: ctrl}
	mock.recorder = &MockHistogramMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHistogram) EXPECT() *MockHistogramMockRecorder {
	return m.recorder
}

// Observe mocks base method.
func (m *MockHistogram) Observe(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Observe", arg0)
}

// Observe indicates an expected call of Observe.
func (mr *MockHistogramMockRecorder) Observe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Observe", reflect.TypeOf((*MockHistogram)(nil).Observe), arg0)
}
