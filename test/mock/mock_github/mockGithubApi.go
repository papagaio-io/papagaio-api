// Code generated by MockGen. DO NOT EDIT.
// Source: api/git/github/githubApi.go

// Package mock_github is a generated GoMock package.
package mock_github

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	dto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	github "wecode.sorint.it/opensource/papagaio-api/api/git/github"
	common "wecode.sorint.it/opensource/papagaio-api/common"
	model "wecode.sorint.it/opensource/papagaio-api/model"
)

// MockGithubInterface is a mock of GithubInterface interface
type MockGithubInterface struct {
	ctrl     *gomock.Controller
	recorder *MockGithubInterfaceMockRecorder
}

// MockGithubInterfaceMockRecorder is the mock recorder for MockGithubInterface
type MockGithubInterfaceMockRecorder struct {
	mock *MockGithubInterface
}

// NewMockGithubInterface creates a new mock instance
func NewMockGithubInterface(ctrl *gomock.Controller) *MockGithubInterface {
	mock := &MockGithubInterface{ctrl: ctrl}
	mock.recorder = &MockGithubInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGithubInterface) EXPECT() *MockGithubInterfaceMockRecorder {
	return m.recorder
}

// CreateWebHook mocks base method
func (m *MockGithubInterface) CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef, organizationRef string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebHook", gitSource, user, gitOrgRef, organizationRef)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebHook indicates an expected call of CreateWebHook
func (mr *MockGithubInterfaceMockRecorder) CreateWebHook(gitSource, user, gitOrgRef, organizationRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebHook", reflect.TypeOf((*MockGithubInterface)(nil).CreateWebHook), gitSource, user, gitOrgRef, organizationRef)
}

// DeleteWebHook mocks base method
func (m *MockGithubInterface) DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWebHook", gitSource, user, gitOrgRef, webHookID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWebHook indicates an expected call of DeleteWebHook
func (mr *MockGithubInterfaceMockRecorder) DeleteWebHook(gitSource, user, gitOrgRef, webHookID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebHook", reflect.TypeOf((*MockGithubInterface)(nil).DeleteWebHook), gitSource, user, gitOrgRef, webHookID)
}

// GetRepositories mocks base method
func (m *MockGithubInterface) GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositories", gitSource, user, gitOrgRef)
	ret0, _ := ret[0].(*[]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepositories indicates an expected call of GetRepositories
func (mr *MockGithubInterfaceMockRecorder) GetRepositories(gitSource, user, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositories", reflect.TypeOf((*MockGithubInterface)(nil).GetRepositories), gitSource, user, gitOrgRef)
}

// GetEmailsRepositoryUsersOwner mocks base method
func (m *MockGithubInterface) GetEmailsRepositoryUsersOwner(gitSource *model.GitSource, user *model.User, gitOrgRef, repositoryRef string) (*[]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailsRepositoryUsersOwner", gitSource, user, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(*[]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmailsRepositoryUsersOwner indicates an expected call of GetEmailsRepositoryUsersOwner
func (mr *MockGithubInterfaceMockRecorder) GetEmailsRepositoryUsersOwner(gitSource, user, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailsRepositoryUsersOwner", reflect.TypeOf((*MockGithubInterface)(nil).GetEmailsRepositoryUsersOwner), gitSource, user, gitOrgRef, repositoryRef)
}

// GetOrganizationMembers mocks base method
func (m *MockGithubInterface) GetOrganizationMembers(gitSource *model.GitSource, user *model.User, organizationName string) (*[]github.GitHubUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationMembers", gitSource, user, organizationName)
	ret0, _ := ret[0].(*[]github.GitHubUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationMembers indicates an expected call of GetOrganizationMembers
func (mr *MockGithubInterfaceMockRecorder) GetOrganizationMembers(gitSource, user, organizationName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationMembers", reflect.TypeOf((*MockGithubInterface)(nil).GetOrganizationMembers), gitSource, user, organizationName)
}

// GetBranches mocks base method
func (m *MockGithubInterface) GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef, repositoryRef string) (map[string]bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranches", gitSource, user, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(map[string]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBranches indicates an expected call of GetBranches
func (mr *MockGithubInterfaceMockRecorder) GetBranches(gitSource, user, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranches", reflect.TypeOf((*MockGithubInterface)(nil).GetBranches), gitSource, user, gitOrgRef, repositoryRef)
}

// CheckRepositoryAgolaConfExists mocks base method
func (m *MockGithubInterface) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef, repositoryRef string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRepositoryAgolaConfExists", gitSource, user, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRepositoryAgolaConfExists indicates an expected call of CheckRepositoryAgolaConfExists
func (mr *MockGithubInterfaceMockRecorder) CheckRepositoryAgolaConfExists(gitSource, user, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRepositoryAgolaConfExists", reflect.TypeOf((*MockGithubInterface)(nil).CheckRepositoryAgolaConfExists), gitSource, user, gitOrgRef, repositoryRef)
}

// GetCommitMetadata mocks base method
func (m *MockGithubInterface) GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef, repositoryRef, commitSha string) (*dto.CommitMetadataDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitMetadata", gitSource, user, gitOrgRef, repositoryRef, commitSha)
	ret0, _ := ret[0].(*dto.CommitMetadataDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitMetadata indicates an expected call of GetCommitMetadata
func (mr *MockGithubInterfaceMockRecorder) GetCommitMetadata(gitSource, user, gitOrgRef, repositoryRef, commitSha interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitMetadata", reflect.TypeOf((*MockGithubInterface)(nil).GetCommitMetadata), gitSource, user, gitOrgRef, repositoryRef, commitSha)
}

// GetOrganization mocks base method
func (m *MockGithubInterface) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*dto.OrganizationDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganization", gitSource, user, gitOrgRef)
	ret0, _ := ret[0].(*dto.OrganizationDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganization indicates an expected call of GetOrganization
func (mr *MockGithubInterfaceMockRecorder) GetOrganization(gitSource, user, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganization", reflect.TypeOf((*MockGithubInterface)(nil).GetOrganization), gitSource, user, gitOrgRef)
}

// GetOrganizations mocks base method
func (m *MockGithubInterface) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizations", gitSource, user)
	ret0, _ := ret[0].(*[]dto.OrganizationDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizations indicates an expected call of GetOrganizations
func (mr *MockGithubInterfaceMockRecorder) GetOrganizations(gitSource, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizations", reflect.TypeOf((*MockGithubInterface)(nil).GetOrganizations), gitSource, user)
}

// IsUserOwner mocks base method
func (m *MockGithubInterface) IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserOwner", gitSource, user, gitOrgRef)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserOwner indicates an expected call of IsUserOwner
func (mr *MockGithubInterfaceMockRecorder) IsUserOwner(gitSource, user, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserOwner", reflect.TypeOf((*MockGithubInterface)(nil).IsUserOwner), gitSource, user, gitOrgRef)
}

// GetUserInfo mocks base method
func (m *MockGithubInterface) GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfo", gitSource, user)
	ret0, _ := ret[0].(*dto.UserInfoDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfo indicates an expected call of GetUserInfo
func (mr *MockGithubInterfaceMockRecorder) GetUserInfo(gitSource, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfo", reflect.TypeOf((*MockGithubInterface)(nil).GetUserInfo), gitSource, user)
}

// GetUserByLogin mocks base method
func (m *MockGithubInterface) GetUserByLogin(gitSource *model.GitSource, login string) (*dto.UserInfoDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", gitSource, login)
	ret0, _ := ret[0].(*dto.UserInfoDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin
func (mr *MockGithubInterfaceMockRecorder) GetUserByLogin(gitSource, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockGithubInterface)(nil).GetUserByLogin), gitSource, login)
}

// GetOauth2AccessToken mocks base method
func (m *MockGithubInterface) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOauth2AccessToken", gitSource, code)
	ret0, _ := ret[0].(*common.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOauth2AccessToken indicates an expected call of GetOauth2AccessToken
func (mr *MockGithubInterfaceMockRecorder) GetOauth2AccessToken(gitSource, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOauth2AccessToken", reflect.TypeOf((*MockGithubInterface)(nil).GetOauth2AccessToken), gitSource, code)
}

// RefreshToken mocks base method
func (m *MockGithubInterface) RefreshToken(gitSource *model.GitSource, refreshToken string) (*common.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", gitSource, refreshToken)
	ret0, _ := ret[0].(*common.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken
func (mr *MockGithubInterfaceMockRecorder) RefreshToken(gitSource, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockGithubInterface)(nil).RefreshToken), gitSource, refreshToken)
}
