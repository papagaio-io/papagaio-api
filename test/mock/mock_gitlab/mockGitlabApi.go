// Code generated by MockGen. DO NOT EDIT.
// Source: api/git/gitlab/gitlabApi.go

// Package mock_gitlab is a generated GoMock package.
package mock_gitlab

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	dto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	gitlab "wecode.sorint.it/opensource/papagaio-api/api/git/gitlab"
	common "wecode.sorint.it/opensource/papagaio-api/common"
	model "wecode.sorint.it/opensource/papagaio-api/model"
)

// MockGitlabInterface is a mock of GitlabInterface interface
type MockGitlabInterface struct {
	ctrl     *gomock.Controller
	recorder *MockGitlabInterfaceMockRecorder
}

// MockGitlabInterfaceMockRecorder is the mock recorder for MockGitlabInterface
type MockGitlabInterfaceMockRecorder struct {
	mock *MockGitlabInterface
}

// NewMockGitlabInterface creates a new mock instance
func NewMockGitlabInterface(ctrl *gomock.Controller) *MockGitlabInterface {
	mock := &MockGitlabInterface{ctrl: ctrl}
	mock.recorder = &MockGitlabInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGitlabInterface) EXPECT() *MockGitlabInterfaceMockRecorder {
	return m.recorder
}

// CreateWebHook mocks base method
func (m *MockGitlabInterface) CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef, organizationRef string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebHook", gitSource, user, gitOrgRef, organizationRef)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebHook indicates an expected call of CreateWebHook
func (mr *MockGitlabInterfaceMockRecorder) CreateWebHook(gitSource, user, gitOrgRef, organizationRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebHook", reflect.TypeOf((*MockGitlabInterface)(nil).CreateWebHook), gitSource, user, gitOrgRef, organizationRef)
}

// DeleteWebHook mocks base method
func (m *MockGitlabInterface) DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWebHook", gitSource, user, gitOrgRef, webHookID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWebHook indicates an expected call of DeleteWebHook
func (mr *MockGitlabInterfaceMockRecorder) DeleteWebHook(gitSource, user, gitOrgRef, webHookID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebHook", reflect.TypeOf((*MockGitlabInterface)(nil).DeleteWebHook), gitSource, user, gitOrgRef, webHookID)
}

// GetRepositories mocks base method
func (m *MockGitlabInterface) GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositories", gitSource, user, gitOrgRef)
	ret0, _ := ret[0].(*[]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepositories indicates an expected call of GetRepositories
func (mr *MockGitlabInterfaceMockRecorder) GetRepositories(gitSource, user, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositories", reflect.TypeOf((*MockGitlabInterface)(nil).GetRepositories), gitSource, user, gitOrgRef)
}

// GetEmailsRepositoryUsersOwner mocks base method
func (m *MockGitlabInterface) GetEmailsRepositoryUsersOwner(gitSource *model.GitSource, user *model.User, gitOrgRef, repositoryRef string) (*[]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailsRepositoryUsersOwner", gitSource, user, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(*[]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmailsRepositoryUsersOwner indicates an expected call of GetEmailsRepositoryUsersOwner
func (mr *MockGitlabInterfaceMockRecorder) GetEmailsRepositoryUsersOwner(gitSource, user, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailsRepositoryUsersOwner", reflect.TypeOf((*MockGitlabInterface)(nil).GetEmailsRepositoryUsersOwner), gitSource, user, gitOrgRef, repositoryRef)
}

// GetOrganizationMembers mocks base method
func (m *MockGitlabInterface) GetOrganizationMembers(gitSource *model.GitSource, user *model.User, organizationName string) (*[]gitlab.GitlabUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationMembers", gitSource, user, organizationName)
	ret0, _ := ret[0].(*[]gitlab.GitlabUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationMembers indicates an expected call of GetOrganizationMembers
func (mr *MockGitlabInterfaceMockRecorder) GetOrganizationMembers(gitSource, user, organizationName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationMembers", reflect.TypeOf((*MockGitlabInterface)(nil).GetOrganizationMembers), gitSource, user, organizationName)
}

// GetBranches mocks base method
func (m *MockGitlabInterface) GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef, repositoryRef string) (map[string]bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranches", gitSource, user, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(map[string]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBranches indicates an expected call of GetBranches
func (mr *MockGitlabInterfaceMockRecorder) GetBranches(gitSource, user, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranches", reflect.TypeOf((*MockGitlabInterface)(nil).GetBranches), gitSource, user, gitOrgRef, repositoryRef)
}

// CheckRepositoryAgolaConfExists mocks base method
func (m *MockGitlabInterface) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef, repositoryRef string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRepositoryAgolaConfExists", gitSource, user, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRepositoryAgolaConfExists indicates an expected call of CheckRepositoryAgolaConfExists
func (mr *MockGitlabInterfaceMockRecorder) CheckRepositoryAgolaConfExists(gitSource, user, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRepositoryAgolaConfExists", reflect.TypeOf((*MockGitlabInterface)(nil).CheckRepositoryAgolaConfExists), gitSource, user, gitOrgRef, repositoryRef)
}

// GetCommitMetadata mocks base method
func (m *MockGitlabInterface) GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef, repositoryRef, commitSha string) (*dto.CommitMetadataDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitMetadata", gitSource, user, gitOrgRef, repositoryRef, commitSha)
	ret0, _ := ret[0].(*dto.CommitMetadataDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitMetadata indicates an expected call of GetCommitMetadata
func (mr *MockGitlabInterfaceMockRecorder) GetCommitMetadata(gitSource, user, gitOrgRef, repositoryRef, commitSha interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitMetadata", reflect.TypeOf((*MockGitlabInterface)(nil).GetCommitMetadata), gitSource, user, gitOrgRef, repositoryRef, commitSha)
}

// GetOrganization mocks base method
func (m *MockGitlabInterface) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*dto.OrganizationDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganization", gitSource, user, gitOrgRef)
	ret0, _ := ret[0].(*dto.OrganizationDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganization indicates an expected call of GetOrganization
func (mr *MockGitlabInterfaceMockRecorder) GetOrganization(gitSource, user, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganization", reflect.TypeOf((*MockGitlabInterface)(nil).GetOrganization), gitSource, user, gitOrgRef)
}

// GetOrganizations mocks base method
func (m *MockGitlabInterface) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizations", gitSource, user)
	ret0, _ := ret[0].(*[]dto.OrganizationDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizations indicates an expected call of GetOrganizations
func (mr *MockGitlabInterfaceMockRecorder) GetOrganizations(gitSource, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizations", reflect.TypeOf((*MockGitlabInterface)(nil).GetOrganizations), gitSource, user)
}

// IsUserOwner mocks base method
func (m *MockGitlabInterface) IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserOwner", gitSource, user, gitOrgRef)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserOwner indicates an expected call of IsUserOwner
func (mr *MockGitlabInterfaceMockRecorder) IsUserOwner(gitSource, user, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserOwner", reflect.TypeOf((*MockGitlabInterface)(nil).IsUserOwner), gitSource, user, gitOrgRef)
}

// GetUserInfo mocks base method
func (m *MockGitlabInterface) GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfo", gitSource, user)
	ret0, _ := ret[0].(*dto.UserInfoDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfo indicates an expected call of GetUserInfo
func (mr *MockGitlabInterfaceMockRecorder) GetUserInfo(gitSource, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfo", reflect.TypeOf((*MockGitlabInterface)(nil).GetUserInfo), gitSource, user)
}

// GetUserByLogin mocks base method
func (m *MockGitlabInterface) GetUserByLogin(gitSource *model.GitSource, id int) (*dto.UserInfoDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", gitSource, id)
	ret0, _ := ret[0].(*dto.UserInfoDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin
func (mr *MockGitlabInterfaceMockRecorder) GetUserByLogin(gitSource, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockGitlabInterface)(nil).GetUserByLogin), gitSource, id)
}

// GetOauth2AccessToken mocks base method
func (m *MockGitlabInterface) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOauth2AccessToken", gitSource, code)
	ret0, _ := ret[0].(*common.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOauth2AccessToken indicates an expected call of GetOauth2AccessToken
func (mr *MockGitlabInterfaceMockRecorder) GetOauth2AccessToken(gitSource, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOauth2AccessToken", reflect.TypeOf((*MockGitlabInterface)(nil).GetOauth2AccessToken), gitSource, code)
}

// RefreshToken mocks base method
func (m *MockGitlabInterface) RefreshToken(gitSource *model.GitSource, refreshToken string) (*common.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", gitSource, refreshToken)
	ret0, _ := ret[0].(*common.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken
func (mr *MockGitlabInterfaceMockRecorder) RefreshToken(gitSource, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockGitlabInterface)(nil).RefreshToken), gitSource, refreshToken)
}
