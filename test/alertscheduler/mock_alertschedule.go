// Code generated by MockGen. DO NOT EDIT.
// Source: alertscheduler/alertschedule_IF.go

// Package mock_alertscheduler is a generated GoMock package.
package mock_alertscheduler

import (
	gomock "github.com/golang/mock/gomock"
	helper "github.com/thomasobenaus/sokar/helper"
	reflect "reflect"
	time "time"
)

// MockAlertSchedule is a mock of AlertSchedule interface
type MockAlertSchedule struct {
	ctrl     *gomock.Controller
	recorder *MockAlertScheduleMockRecorder
}

// MockAlertScheduleMockRecorder is the mock recorder for MockAlertSchedule
type MockAlertScheduleMockRecorder struct {
	mock *MockAlertSchedule
}

// NewMockAlertSchedule creates a new mock instance
func NewMockAlertSchedule(ctrl *gomock.Controller) *MockAlertSchedule {
	mock := &MockAlertSchedule{ctrl: ctrl}
	mock.recorder = &MockAlertScheduleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAlertSchedule) EXPECT() *MockAlertScheduleMockRecorder {
	return m.recorder
}

// IsActiveAt mocks base method
func (m *MockAlertSchedule) IsActiveAt(day time.Weekday, at helper.SimpleTime) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsActiveAt", day, at)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsActiveAt indicates an expected call of IsActiveAt
func (mr *MockAlertScheduleMockRecorder) IsActiveAt(day, at interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsActiveAt", reflect.TypeOf((*MockAlertSchedule)(nil).IsActiveAt), day, at)
}
