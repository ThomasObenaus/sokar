// Code generated by MockGen. DO NOT EDIT.
// Source: nomad/nomadclient_IF.go

// Package mock_nomad is a generated GoMock package.
package mock_nomad

import (
	gomock "github.com/golang/mock/gomock"
	api "github.com/hashicorp/nomad/api"
	reflect "reflect"
)

// MockNomadJobs is a mock of NomadJobs interface
type MockNomadJobs struct {
	ctrl     *gomock.Controller
	recorder *MockNomadJobsMockRecorder
}

// MockNomadJobsMockRecorder is the mock recorder for MockNomadJobs
type MockNomadJobsMockRecorder struct {
	mock *MockNomadJobs
}

// NewMockNomadJobs creates a new mock instance
func NewMockNomadJobs(ctrl *gomock.Controller) *MockNomadJobs {
	mock := &MockNomadJobs{ctrl: ctrl}
	mock.recorder = &MockNomadJobsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNomadJobs) EXPECT() *MockNomadJobsMockRecorder {
	return m.recorder
}

// Info mocks base method
func (m *MockNomadJobs) Info(jobID string, q *api.QueryOptions) (*api.Job, *api.QueryMeta, error) {
	ret := m.ctrl.Call(m, "Info", jobID, q)
	ret0, _ := ret[0].(*api.Job)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Info indicates an expected call of Info
func (mr *MockNomadJobsMockRecorder) Info(jobID, q interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockNomadJobs)(nil).Info), jobID, q)
}

// Register mocks base method
func (m *MockNomadJobs) Register(job *api.Job, q *api.WriteOptions) (*api.JobRegisterResponse, *api.WriteMeta, error) {
	ret := m.ctrl.Call(m, "Register", job, q)
	ret0, _ := ret[0].(*api.JobRegisterResponse)
	ret1, _ := ret[1].(*api.WriteMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Register indicates an expected call of Register
func (mr *MockNomadJobsMockRecorder) Register(job, q interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockNomadJobs)(nil).Register), job, q)
}

// MockNomadDeployments is a mock of NomadDeployments interface
type MockNomadDeployments struct {
	ctrl     *gomock.Controller
	recorder *MockNomadDeploymentsMockRecorder
}

// MockNomadDeploymentsMockRecorder is the mock recorder for MockNomadDeployments
type MockNomadDeploymentsMockRecorder struct {
	mock *MockNomadDeployments
}

// NewMockNomadDeployments creates a new mock instance
func NewMockNomadDeployments(ctrl *gomock.Controller) *MockNomadDeployments {
	mock := &MockNomadDeployments{ctrl: ctrl}
	mock.recorder = &MockNomadDeploymentsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNomadDeployments) EXPECT() *MockNomadDeploymentsMockRecorder {
	return m.recorder
}

// Info mocks base method
func (m *MockNomadDeployments) Info(deploymentID string, q *api.QueryOptions) (*api.Deployment, *api.QueryMeta, error) {
	ret := m.ctrl.Call(m, "Info", deploymentID, q)
	ret0, _ := ret[0].(*api.Deployment)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Info indicates an expected call of Info
func (mr *MockNomadDeploymentsMockRecorder) Info(deploymentID, q interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockNomadDeployments)(nil).Info), deploymentID, q)
}

// MockNomadEvaluations is a mock of NomadEvaluations interface
type MockNomadEvaluations struct {
	ctrl     *gomock.Controller
	recorder *MockNomadEvaluationsMockRecorder
}

// MockNomadEvaluationsMockRecorder is the mock recorder for MockNomadEvaluations
type MockNomadEvaluationsMockRecorder struct {
	mock *MockNomadEvaluations
}

// NewMockNomadEvaluations creates a new mock instance
func NewMockNomadEvaluations(ctrl *gomock.Controller) *MockNomadEvaluations {
	mock := &MockNomadEvaluations{ctrl: ctrl}
	mock.recorder = &MockNomadEvaluationsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNomadEvaluations) EXPECT() *MockNomadEvaluationsMockRecorder {
	return m.recorder
}

// Info mocks base method
func (m *MockNomadEvaluations) Info(evalID string, q *api.QueryOptions) (*api.Evaluation, *api.QueryMeta, error) {
	ret := m.ctrl.Call(m, "Info", evalID, q)
	ret0, _ := ret[0].(*api.Evaluation)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Info indicates an expected call of Info
func (mr *MockNomadEvaluationsMockRecorder) Info(evalID, q interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockNomadEvaluations)(nil).Info), evalID, q)
}
