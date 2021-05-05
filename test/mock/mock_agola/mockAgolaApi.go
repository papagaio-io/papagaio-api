// Code generated by MockGen. DO NOT EDIT.
// Source: api/agola/agolaApi.go

// Package mock_agola is a generated GoMock package.
package mock_agola

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	agola "wecode.sorint.it/opensource/papagaio-api/api/agola"
	dto "wecode.sorint.it/opensource/papagaio-api/dto"
	model "wecode.sorint.it/opensource/papagaio-api/model"
)

// MockAgolaApiInterface is a mock of AgolaApiInterface interface
type MockAgolaApiInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAgolaApiInterfaceMockRecorder
}

// MockAgolaApiInterfaceMockRecorder is the mock recorder for MockAgolaApiInterface
type MockAgolaApiInterfaceMockRecorder struct {
	mock *MockAgolaApiInterface
}

// NewMockAgolaApiInterface creates a new mock instance
func NewMockAgolaApiInterface(ctrl *gomock.Controller) *MockAgolaApiInterface {
	mock := &MockAgolaApiInterface{ctrl: ctrl}
	mock.recorder = &MockAgolaApiInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAgolaApiInterface) EXPECT() *MockAgolaApiInterfaceMockRecorder {
	return m.recorder
}

// CheckOrganizationExists mocks base method
func (m *MockAgolaApiInterface) CheckOrganizationExists(organization *model.Organization) (bool, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckOrganizationExists", organization)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// CheckOrganizationExists indicates an expected call of CheckOrganizationExists
func (mr *MockAgolaApiInterfaceMockRecorder) CheckOrganizationExists(organization interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckOrganizationExists", reflect.TypeOf((*MockAgolaApiInterface)(nil).CheckOrganizationExists), organization)
}

// CheckProjectExists mocks base method
func (m *MockAgolaApiInterface) CheckProjectExists(organization *model.Organization, projectName string) (bool, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckProjectExists", organization, projectName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// CheckProjectExists indicates an expected call of CheckProjectExists
func (mr *MockAgolaApiInterfaceMockRecorder) CheckProjectExists(organization, projectName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckProjectExists", reflect.TypeOf((*MockAgolaApiInterface)(nil).CheckProjectExists), organization, projectName)
}

// CreateOrganization mocks base method
func (m *MockAgolaApiInterface) CreateOrganization(organization *model.Organization, visibility dto.VisibilityType) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrganization", organization, visibility)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrganization indicates an expected call of CreateOrganization
func (mr *MockAgolaApiInterfaceMockRecorder) CreateOrganization(organization, visibility interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrganization", reflect.TypeOf((*MockAgolaApiInterface)(nil).CreateOrganization), organization, visibility)
}

// DeleteOrganization mocks base method
func (m *MockAgolaApiInterface) DeleteOrganization(organization *model.Organization, agolaUserToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrganization", organization, agolaUserToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrganization indicates an expected call of DeleteOrganization
func (mr *MockAgolaApiInterfaceMockRecorder) DeleteOrganization(organization, agolaUserToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrganization", reflect.TypeOf((*MockAgolaApiInterface)(nil).DeleteOrganization), organization, agolaUserToken)
}

// CreateProject mocks base method
func (m *MockAgolaApiInterface) CreateProject(projectName, agolaProjectRef string, organization *model.Organization, remoteSourceName, agolaUserToken string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", projectName, agolaProjectRef, organization, remoteSourceName, agolaUserToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject
func (mr *MockAgolaApiInterfaceMockRecorder) CreateProject(projectName, agolaProjectRef, organization, remoteSourceName, agolaUserToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockAgolaApiInterface)(nil).CreateProject), projectName, agolaProjectRef, organization, remoteSourceName, agolaUserToken)
}

// DeleteProject mocks base method
func (m *MockAgolaApiInterface) DeleteProject(organization *model.Organization, agolaProjectRef, agolaUserToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", organization, agolaProjectRef, agolaUserToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject
func (mr *MockAgolaApiInterfaceMockRecorder) DeleteProject(organization, agolaProjectRef, agolaUserToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockAgolaApiInterface)(nil).DeleteProject), organization, agolaProjectRef, agolaUserToken)
}

// AddOrUpdateOrganizationMember mocks base method
func (m *MockAgolaApiInterface) AddOrUpdateOrganizationMember(organization *model.Organization, agolaUserRef, role string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrUpdateOrganizationMember", organization, agolaUserRef, role)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrUpdateOrganizationMember indicates an expected call of AddOrUpdateOrganizationMember
func (mr *MockAgolaApiInterfaceMockRecorder) AddOrUpdateOrganizationMember(organization, agolaUserRef, role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrUpdateOrganizationMember", reflect.TypeOf((*MockAgolaApiInterface)(nil).AddOrUpdateOrganizationMember), organization, agolaUserRef, role)
}

// RemoveOrganizationMember mocks base method
func (m *MockAgolaApiInterface) RemoveOrganizationMember(organization *model.Organization, agolaUserRef string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOrganizationMember", organization, agolaUserRef)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveOrganizationMember indicates an expected call of RemoveOrganizationMember
func (mr *MockAgolaApiInterfaceMockRecorder) RemoveOrganizationMember(organization, agolaUserRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOrganizationMember", reflect.TypeOf((*MockAgolaApiInterface)(nil).RemoveOrganizationMember), organization, agolaUserRef)
}

// GetOrganizationMembers mocks base method
func (m *MockAgolaApiInterface) GetOrganizationMembers(organization *model.Organization) (*agola.OrganizationMembersResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationMembers", organization)
	ret0, _ := ret[0].(*agola.OrganizationMembersResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationMembers indicates an expected call of GetOrganizationMembers
func (mr *MockAgolaApiInterfaceMockRecorder) GetOrganizationMembers(organization interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationMembers", reflect.TypeOf((*MockAgolaApiInterface)(nil).GetOrganizationMembers), organization)
}

// ArchiveProject mocks base method
func (m *MockAgolaApiInterface) ArchiveProject(organization *model.Organization, agolaProjectRef string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveProject", organization, agolaProjectRef)
	ret0, _ := ret[0].(error)
	return ret0
}

// ArchiveProject indicates an expected call of ArchiveProject
func (mr *MockAgolaApiInterfaceMockRecorder) ArchiveProject(organization, agolaProjectRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveProject", reflect.TypeOf((*MockAgolaApiInterface)(nil).ArchiveProject), organization, agolaProjectRef)
}

// UnarchiveProject mocks base method
func (m *MockAgolaApiInterface) UnarchiveProject(organization *model.Organization, agolaProjectRef string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnarchiveProject", organization, agolaProjectRef)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnarchiveProject indicates an expected call of UnarchiveProject
func (mr *MockAgolaApiInterfaceMockRecorder) UnarchiveProject(organization, agolaProjectRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnarchiveProject", reflect.TypeOf((*MockAgolaApiInterface)(nil).UnarchiveProject), organization, agolaProjectRef)
}

// GetRuns mocks base method
func (m *MockAgolaApiInterface) GetRuns(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) (*[]agola.RunDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRuns", projectRef, lastRun, phase, startRunID, limit, asc)
	ret0, _ := ret[0].(*[]agola.RunDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRuns indicates an expected call of GetRuns
func (mr *MockAgolaApiInterfaceMockRecorder) GetRuns(projectRef, lastRun, phase, startRunID, limit, asc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRuns", reflect.TypeOf((*MockAgolaApiInterface)(nil).GetRuns), projectRef, lastRun, phase, startRunID, limit, asc)
}

// GetRun mocks base method
func (m *MockAgolaApiInterface) GetRun(runID string) (*agola.RunDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRun", runID)
	ret0, _ := ret[0].(*agola.RunDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRun indicates an expected call of GetRun
func (mr *MockAgolaApiInterfaceMockRecorder) GetRun(runID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRun", reflect.TypeOf((*MockAgolaApiInterface)(nil).GetRun), runID)
}

// GetTask mocks base method
func (m *MockAgolaApiInterface) GetTask(runID, taskID string) (*agola.TaskDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", runID, taskID)
	ret0, _ := ret[0].(*agola.TaskDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask
func (mr *MockAgolaApiInterfaceMockRecorder) GetTask(runID, taskID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockAgolaApiInterface)(nil).GetTask), runID, taskID)
}

// GetLogs mocks base method
func (m *MockAgolaApiInterface) GetLogs(runID, taskID string, step int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs", runID, taskID, step)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs
func (mr *MockAgolaApiInterfaceMockRecorder) GetLogs(runID, taskID, step interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockAgolaApiInterface)(nil).GetLogs), runID, taskID, step)
}

// GetRemoteSource mocks base method
func (m *MockAgolaApiInterface) GetRemoteSource(agolaRemoteSource string) (*agola.RemoteSourceDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemoteSource", agolaRemoteSource)
	ret0, _ := ret[0].(*agola.RemoteSourceDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRemoteSource indicates an expected call of GetRemoteSource
func (mr *MockAgolaApiInterfaceMockRecorder) GetRemoteSource(agolaRemoteSource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemoteSource", reflect.TypeOf((*MockAgolaApiInterface)(nil).GetRemoteSource), agolaRemoteSource)
}

// GetUsers mocks base method
func (m *MockAgolaApiInterface) GetUsers() (*[]agola.UserDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers")
	ret0, _ := ret[0].(*[]agola.UserDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers
func (mr *MockAgolaApiInterfaceMockRecorder) GetUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockAgolaApiInterface)(nil).GetUsers))
}
