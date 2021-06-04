// Code generated by MockGen. DO NOT EDIT.
// Source: repository/database.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	model "wecode.sorint.it/opensource/papagaio-api/model"
)

// MockDatabase is a mock of Database interface
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// GetOrganizationsRef mocks base method
func (m *MockDatabase) GetOrganizationsRef() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationsRef")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationsRef indicates an expected call of GetOrganizationsRef
func (mr *MockDatabaseMockRecorder) GetOrganizationsRef() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationsRef", reflect.TypeOf((*MockDatabase)(nil).GetOrganizationsRef))
}

// GetOrganizations mocks base method
func (m *MockDatabase) GetOrganizations() (*[]model.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizations")
	ret0, _ := ret[0].(*[]model.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizations indicates an expected call of GetOrganizations
func (mr *MockDatabaseMockRecorder) GetOrganizations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizations", reflect.TypeOf((*MockDatabase)(nil).GetOrganizations))
}

// SaveOrganization mocks base method
func (m *MockDatabase) SaveOrganization(organization *model.Organization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveOrganization", organization)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveOrganization indicates an expected call of SaveOrganization
func (mr *MockDatabaseMockRecorder) SaveOrganization(organization interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveOrganization", reflect.TypeOf((*MockDatabase)(nil).SaveOrganization), organization)
}

// GetOrganizationByAgolaRef mocks base method
func (m *MockDatabase) GetOrganizationByAgolaRef(organizationName string) (*model.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationByAgolaRef", organizationName)
	ret0, _ := ret[0].(*model.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationByAgolaRef indicates an expected call of GetOrganizationByAgolaRef
func (mr *MockDatabaseMockRecorder) GetOrganizationByAgolaRef(organizationName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationByAgolaRef", reflect.TypeOf((*MockDatabase)(nil).GetOrganizationByAgolaRef), organizationName)
}

// GetOrganizationById mocks base method
func (m *MockDatabase) GetOrganizationById(organizationID string) (*model.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationById", organizationID)
	ret0, _ := ret[0].(*model.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationById indicates an expected call of GetOrganizationById
func (mr *MockDatabaseMockRecorder) GetOrganizationById(organizationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationById", reflect.TypeOf((*MockDatabase)(nil).GetOrganizationById), organizationID)
}

// DeleteOrganization mocks base method
func (m *MockDatabase) DeleteOrganization(organizationID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrganization", organizationID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrganization indicates an expected call of DeleteOrganization
func (mr *MockDatabaseMockRecorder) DeleteOrganization(organizationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrganization", reflect.TypeOf((*MockDatabase)(nil).DeleteOrganization), organizationID)
}

// GetOrganizationsByGitSource mocks base method
func (m *MockDatabase) GetOrganizationsByGitSource(gitSource string) (*[]model.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationsByGitSource", gitSource)
	ret0, _ := ret[0].(*[]model.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationsByGitSource indicates an expected call of GetOrganizationsByGitSource
func (mr *MockDatabaseMockRecorder) GetOrganizationsByGitSource(gitSource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationsByGitSource", reflect.TypeOf((*MockDatabase)(nil).GetOrganizationsByGitSource), gitSource)
}

// GetGitSources mocks base method
func (m *MockDatabase) GetGitSources() (*[]model.GitSource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitSources")
	ret0, _ := ret[0].(*[]model.GitSource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitSources indicates an expected call of GetGitSources
func (mr *MockDatabaseMockRecorder) GetGitSources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitSources", reflect.TypeOf((*MockDatabase)(nil).GetGitSources))
}

// SaveGitSource mocks base method
func (m *MockDatabase) SaveGitSource(gitSource *model.GitSource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveGitSource", gitSource)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveGitSource indicates an expected call of SaveGitSource
func (mr *MockDatabaseMockRecorder) SaveGitSource(gitSource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveGitSource", reflect.TypeOf((*MockDatabase)(nil).SaveGitSource), gitSource)
}

// GetGitSourceById mocks base method
func (m *MockDatabase) GetGitSourceById(id string) (*model.GitSource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitSourceById", id)
	ret0, _ := ret[0].(*model.GitSource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitSourceById indicates an expected call of GetGitSourceById
func (mr *MockDatabaseMockRecorder) GetGitSourceById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitSourceById", reflect.TypeOf((*MockDatabase)(nil).GetGitSourceById), id)
}

// GetGitSourceByName mocks base method
func (m *MockDatabase) GetGitSourceByName(name string) (*model.GitSource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitSourceByName", name)
	ret0, _ := ret[0].(*model.GitSource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitSourceByName indicates an expected call of GetGitSourceByName
func (mr *MockDatabaseMockRecorder) GetGitSourceByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitSourceByName", reflect.TypeOf((*MockDatabase)(nil).GetGitSourceByName), name)
}

// DeleteGitSource mocks base method
func (m *MockDatabase) DeleteGitSource(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGitSource", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGitSource indicates an expected call of DeleteGitSource
func (mr *MockDatabaseMockRecorder) DeleteGitSource(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGitSource", reflect.TypeOf((*MockDatabase)(nil).DeleteGitSource), id)
}

// GetOrganizationsTriggerTime mocks base method
func (m *MockDatabase) GetOrganizationsTriggerTime() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationsTriggerTime")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetOrganizationsTriggerTime indicates an expected call of GetOrganizationsTriggerTime
func (mr *MockDatabaseMockRecorder) GetOrganizationsTriggerTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationsTriggerTime", reflect.TypeOf((*MockDatabase)(nil).GetOrganizationsTriggerTime))
}

// SaveOrganizationsTriggerTime mocks base method
func (m *MockDatabase) SaveOrganizationsTriggerTime(val int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveOrganizationsTriggerTime", val)
}

// SaveOrganizationsTriggerTime indicates an expected call of SaveOrganizationsTriggerTime
func (mr *MockDatabaseMockRecorder) SaveOrganizationsTriggerTime(val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveOrganizationsTriggerTime", reflect.TypeOf((*MockDatabase)(nil).SaveOrganizationsTriggerTime), val)
}

// GetRunFailedTriggerTime mocks base method
func (m *MockDatabase) GetRunFailedTriggerTime() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunFailedTriggerTime")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetRunFailedTriggerTime indicates an expected call of GetRunFailedTriggerTime
func (mr *MockDatabaseMockRecorder) GetRunFailedTriggerTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunFailedTriggerTime", reflect.TypeOf((*MockDatabase)(nil).GetRunFailedTriggerTime))
}

// SaveRunFailedTriggerTime mocks base method
func (m *MockDatabase) SaveRunFailedTriggerTime(val int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveRunFailedTriggerTime", val)
}

// SaveRunFailedTriggerTime indicates an expected call of SaveRunFailedTriggerTime
func (mr *MockDatabaseMockRecorder) SaveRunFailedTriggerTime(val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRunFailedTriggerTime", reflect.TypeOf((*MockDatabase)(nil).SaveRunFailedTriggerTime), val)
}

// GetUserByUserId mocks base method
func (m *MockDatabase) GetUserByUserId(userId uint) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUserId", userId)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUserId indicates an expected call of GetUserByUserId
func (mr *MockDatabaseMockRecorder) GetUserByUserId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUserId", reflect.TypeOf((*MockDatabase)(nil).GetUserByUserId), userId)
}

// GetUserByGitSourceNameAndID mocks base method
func (m *MockDatabase) GetUserByGitSourceNameAndID(gitSourceName string, id uint64) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByGitSourceNameAndID", gitSourceName, id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByGitSourceNameAndID indicates an expected call of GetUserByGitSourceNameAndID
func (mr *MockDatabaseMockRecorder) GetUserByGitSourceNameAndID(gitSourceName, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByGitSourceNameAndID", reflect.TypeOf((*MockDatabase)(nil).GetUserByGitSourceNameAndID), gitSourceName, id)
}

// SaveUser mocks base method
func (m *MockDatabase) SaveUser(user *model.User) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", user)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUser indicates an expected call of SaveUser
func (mr *MockDatabaseMockRecorder) SaveUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*MockDatabase)(nil).SaveUser), user)
}

// DeleteUser mocks base method
func (m *MockDatabase) DeleteUser(userId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockDatabaseMockRecorder) DeleteUser(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockDatabase)(nil).DeleteUser), userId)
}
