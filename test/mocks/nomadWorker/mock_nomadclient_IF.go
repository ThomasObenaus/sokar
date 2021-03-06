// Code generated by MockGen. DO NOT EDIT.
// Source: nomadWorker/nomadclient_IF.go

// Package mock_nomadWorker is a generated GoMock package.
package mock_nomadWorker

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	api "github.com/hashicorp/nomad/api"
	reflect "reflect"
)

// MockNodes is a mock of Nodes interface
type MockNodes struct {
	ctrl     *gomock.Controller
	recorder *MockNodesMockRecorder
}

// MockNodesMockRecorder is the mock recorder for MockNodes
type MockNodesMockRecorder struct {
	mock *MockNodes
}

// NewMockNodes creates a new mock instance
func NewMockNodes(ctrl *gomock.Controller) *MockNodes {
	mock := &MockNodes{ctrl: ctrl}
	mock.recorder = &MockNodesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNodes) EXPECT() *MockNodesMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockNodes) List(q *api.QueryOptions) ([]*api.NodeListStub, *api.QueryMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", q)
	ret0, _ := ret[0].([]*api.NodeListStub)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List
func (mr *MockNodesMockRecorder) List(q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNodes)(nil).List), q)
}

// ToggleEligibility mocks base method
func (m *MockNodes) ToggleEligibility(nodeID string, eligible bool, q *api.WriteOptions) (*api.NodeEligibilityUpdateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToggleEligibility", nodeID, eligible, q)
	ret0, _ := ret[0].(*api.NodeEligibilityUpdateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToggleEligibility indicates an expected call of ToggleEligibility
func (mr *MockNodesMockRecorder) ToggleEligibility(nodeID, eligible, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleEligibility", reflect.TypeOf((*MockNodes)(nil).ToggleEligibility), nodeID, eligible, q)
}

// UpdateDrain mocks base method
func (m *MockNodes) UpdateDrain(nodeID string, spec *api.DrainSpec, markEligible bool, q *api.WriteOptions) (*api.NodeDrainUpdateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDrain", nodeID, spec, markEligible, q)
	ret0, _ := ret[0].(*api.NodeDrainUpdateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateDrain indicates an expected call of UpdateDrain
func (mr *MockNodesMockRecorder) UpdateDrain(nodeID, spec, markEligible, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDrain", reflect.TypeOf((*MockNodes)(nil).UpdateDrain), nodeID, spec, markEligible, q)
}

// MonitorDrain mocks base method
func (m *MockNodes) MonitorDrain(ctx context.Context, nodeID string, index uint64, ignoreSys bool) <-chan *api.MonitorMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MonitorDrain", ctx, nodeID, index, ignoreSys)
	ret0, _ := ret[0].(<-chan *api.MonitorMessage)
	return ret0
}

// MonitorDrain indicates an expected call of MonitorDrain
func (mr *MockNodesMockRecorder) MonitorDrain(ctx, nodeID, index, ignoreSys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MonitorDrain", reflect.TypeOf((*MockNodes)(nil).MonitorDrain), ctx, nodeID, index, ignoreSys)
}

// Allocations mocks base method
func (m *MockNodes) Allocations(nodeID string, q *api.QueryOptions) ([]*api.Allocation, *api.QueryMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Allocations", nodeID, q)
	ret0, _ := ret[0].([]*api.Allocation)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Allocations indicates an expected call of Allocations
func (mr *MockNodesMockRecorder) Allocations(nodeID, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Allocations", reflect.TypeOf((*MockNodes)(nil).Allocations), nodeID, q)
}
