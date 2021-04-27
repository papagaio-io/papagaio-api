// Code generated by MockGen. DO NOT EDIT.
// Source: api/git/github/githubApi.go

// Package mock_github is a generated GoMock package.
package mock_github

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	dto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	github "wecode.sorint.it/opensource/papagaio-api/api/git/github"
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
func (m *MockGithubInterface) CreateWebHook(gitSource *model.GitSource, gitOrgRef string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebHook", gitSource, gitOrgRef)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebHook indicates an expected call of CreateWebHook
func (mr *MockGithubInterfaceMockRecorder) CreateWebHook(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebHook", reflect.TypeOf((*MockGithubInterface)(nil).CreateWebHook), gitSource, gitOrgRef)
}

// DeleteWebHook mocks base method
func (m *MockGithubInterface) DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWebHook", gitSource, gitOrgRef, webHookID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWebHook indicates an expected call of DeleteWebHook
func (mr *MockGithubInterfaceMockRecorder) DeleteWebHook(gitSource, gitOrgRef, webHookID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebHook", reflect.TypeOf((*MockGithubInterface)(nil).DeleteWebHook), gitSource, gitOrgRef, webHookID)
}

// GetRepositories mocks base method
func (m *MockGithubInterface) GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositories", gitSource, gitOrgRef)
	ret0, _ := ret[0].(*[]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepositories indicates an expected call of GetRepositories
func (mr *MockGithubInterfaceMockRecorder) GetRepositories(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositories", reflect.TypeOf((*MockGithubInterface)(nil).GetRepositories), gitSource, gitOrgRef)
}

// CheckOrganizationExists mocks base method
func (m *MockGithubInterface) CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckOrganizationExists", gitSource, gitOrgRef)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckOrganizationExists indicates an expected call of CheckOrganizationExists
func (mr *MockGithubInterfaceMockRecorder) CheckOrganizationExists(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckOrganizationExists", reflect.TypeOf((*MockGithubInterface)(nil).CheckOrganizationExists), gitSource, gitOrgRef)
}

// GetRepositoryTeams mocks base method
func (m *MockGithubInterface) GetRepositoryTeams(gitSource *model.GitSource, gitOrgRef, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositoryTeams", gitSource, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(*[]dto.TeamResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepositoryTeams indicates an expected call of GetRepositoryTeams
func (mr *MockGithubInterfaceMockRecorder) GetRepositoryTeams(gitSource, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositoryTeams", reflect.TypeOf((*MockGithubInterface)(nil).GetRepositoryTeams), gitSource, gitOrgRef, repositoryRef)
}

// GetOrganizationTeams mocks base method
func (m *MockGithubInterface) GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationTeams", gitSource, gitOrgRef)
	ret0, _ := ret[0].(*[]dto.TeamResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationTeams indicates an expected call of GetOrganizationTeams
func (mr *MockGithubInterfaceMockRecorder) GetOrganizationTeams(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationTeams", reflect.TypeOf((*MockGithubInterface)(nil).GetOrganizationTeams), gitSource, gitOrgRef)
}

// GetTeamMembers mocks base method
func (m *MockGithubInterface) GetTeamMembers(gitSource *model.GitSource, organizationName string, teamId int) (*[]dto.UserTeamResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamMembers", gitSource, organizationName, teamId)
	ret0, _ := ret[0].(*[]dto.UserTeamResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeamMembers indicates an expected call of GetTeamMembers
func (mr *MockGithubInterfaceMockRecorder) GetTeamMembers(gitSource, organizationName, teamId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamMembers", reflect.TypeOf((*MockGithubInterface)(nil).GetTeamMembers), gitSource, organizationName, teamId)
}

// GetOrganizationMembers mocks base method
func (m *MockGithubInterface) GetOrganizationMembers(gitSource *model.GitSource, organizationName string) (*[]github.GitHubUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationMembers", gitSource, organizationName)
	ret0, _ := ret[0].(*[]github.GitHubUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationMembers indicates an expected call of GetOrganizationMembers
func (mr *MockGithubInterfaceMockRecorder) GetOrganizationMembers(gitSource, organizationName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationMembers", reflect.TypeOf((*MockGithubInterface)(nil).GetOrganizationMembers), gitSource, organizationName)
}

// GetRepositoryMembers mocks base method
func (m *MockGithubInterface) GetRepositoryMembers(gitSource *model.GitSource, organizationName, repositoryRef string) (*[]github.GitHubUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositoryMembers", gitSource, organizationName, repositoryRef)
	ret0, _ := ret[0].(*[]github.GitHubUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepositoryMembers indicates an expected call of GetRepositoryMembers
func (mr *MockGithubInterfaceMockRecorder) GetRepositoryMembers(gitSource, organizationName, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositoryMembers", reflect.TypeOf((*MockGithubInterface)(nil).GetRepositoryMembers), gitSource, organizationName, repositoryRef)
}

// GetBranches mocks base method
func (m *MockGithubInterface) GetBranches(gitSource *model.GitSource, gitOrgRef, repositoryRef string) map[string]bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranches", gitSource, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(map[string]bool)
	return ret0
}

// GetBranches indicates an expected call of GetBranches
func (mr *MockGithubInterfaceMockRecorder) GetBranches(gitSource, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranches", reflect.TypeOf((*MockGithubInterface)(nil).GetBranches), gitSource, gitOrgRef, repositoryRef)
}

// CheckRepositoryAgolaConfExists mocks base method
func (m *MockGithubInterface) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, gitOrgRef, repositoryRef string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRepositoryAgolaConfExists", gitSource, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRepositoryAgolaConfExists indicates an expected call of CheckRepositoryAgolaConfExists
func (mr *MockGithubInterfaceMockRecorder) CheckRepositoryAgolaConfExists(gitSource, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRepositoryAgolaConfExists", reflect.TypeOf((*MockGithubInterface)(nil).CheckRepositoryAgolaConfExists), gitSource, gitOrgRef, repositoryRef)
}

// GetCommitMetadata mocks base method
func (m *MockGithubInterface) GetCommitMetadata(gitSource *model.GitSource, gitOrgRef, repositoryRef, commitSha string) (*dto.CommitMetadataDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitMetadata", gitSource, gitOrgRef, repositoryRef, commitSha)
	ret0, _ := ret[0].(*dto.CommitMetadataDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitMetadata indicates an expected call of GetCommitMetadata
func (mr *MockGithubInterfaceMockRecorder) GetCommitMetadata(gitSource, gitOrgRef, repositoryRef, commitSha interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitMetadata", reflect.TypeOf((*MockGithubInterface)(nil).GetCommitMetadata), gitSource, gitOrgRef, repositoryRef, commitSha)
}

// GetOrganization mocks base method
func (m *MockGithubInterface) GetOrganization(gitSource *model.GitSource, gitOrgRef string) *dto.OrganizationDto {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganization", gitSource, gitOrgRef)
	ret0, _ := ret[0].(*dto.OrganizationDto)
	return ret0
}

// GetOrganization indicates an expected call of GetOrganization
func (mr *MockGithubInterfaceMockRecorder) GetOrganization(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganization", reflect.TypeOf((*MockGithubInterface)(nil).GetOrganization), gitSource, gitOrgRef)
}
