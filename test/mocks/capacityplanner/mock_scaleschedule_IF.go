// Code generated by MockGen. DO NOT EDIT.
// Source: capacityplanner/scaleschedule_IF.go

// Package mock_capacityplanner is a generated GoMock package.
package mock_capacityplanner

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	helper "github.com/thomasobenaus/sokar/helper"
)

// MockScaleSchedule is a mock of ScaleSchedule interface.
type MockScaleSchedule struct {
	ctrl     *gomock.Controller
	recorder *MockScaleScheduleMockRecorder
}

// MockScaleScheduleMockRecorder is the mock recorder for MockScaleSchedule.
type MockScaleScheduleMockRecorder struct {
	mock *MockScaleSchedule
}

// NewMockScaleSchedule creates a new mock instance.
func NewMockScaleSchedule(ctrl *gomock.Controller) *MockScaleSchedule {
	mock := &MockScaleSchedule{ctrl: ctrl}
	mock.recorder = &MockScaleScheduleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScaleSchedule) EXPECT() *MockScaleScheduleMockRecorder {
	return m.recorder
}

// ScaleRangeAt mocks base method.
func (m *MockScaleSchedule) ScaleRangeAt(day time.Weekday, at helper.SimpleTime) (uint, uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScaleRangeAt", day, at)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(uint)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ScaleRangeAt indicates an expected call of ScaleRangeAt.
func (mr *MockScaleScheduleMockRecorder) ScaleRangeAt(day, at interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScaleRangeAt", reflect.TypeOf((*MockScaleSchedule)(nil).ScaleRangeAt), day, at)
}
