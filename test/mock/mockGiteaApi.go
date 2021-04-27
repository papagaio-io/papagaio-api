// Code generated by MockGen. DO NOT EDIT.
// Source: api/git/gitea/giteaApi.go

// Package mock_gitea is a generated GoMock package.
package mock_gitea

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	dto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	gitea "wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	model "wecode.sorint.it/opensource/papagaio-api/model"
)

// MockGiteaInterface is a mock of GiteaInterface interface
type MockGiteaInterface struct {
	ctrl     *gomock.Controller
	recorder *MockGiteaInterfaceMockRecorder
}

// MockGiteaInterfaceMockRecorder is the mock recorder for MockGiteaInterface
type MockGiteaInterfaceMockRecorder struct {
	mock *MockGiteaInterface
}

// NewMockGiteaInterface creates a new mock instance
func NewMockGiteaInterface(ctrl *gomock.Controller) *MockGiteaInterface {
	mock := &MockGiteaInterface{ctrl: ctrl}
	mock.recorder = &MockGiteaInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGiteaInterface) EXPECT() *MockGiteaInterfaceMockRecorder {
	return m.recorder
}

// CreateWebHook mocks base method
func (m *MockGiteaInterface) CreateWebHook(gitSource *model.GitSource, gitOrgRef string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebHook", gitSource, gitOrgRef)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebHook indicates an expected call of CreateWebHook
func (mr *MockGiteaInterfaceMockRecorder) CreateWebHook(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebHook", reflect.TypeOf((*MockGiteaInterface)(nil).CreateWebHook), gitSource, gitOrgRef)
}

// DeleteWebHook mocks base method
func (m *MockGiteaInterface) DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWebHook", gitSource, gitOrgRef, webHookID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWebHook indicates an expected call of DeleteWebHook
func (mr *MockGiteaInterfaceMockRecorder) DeleteWebHook(gitSource, gitOrgRef, webHookID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebHook", reflect.TypeOf((*MockGiteaInterface)(nil).DeleteWebHook), gitSource, gitOrgRef, webHookID)
}

// GetRepositories mocks base method
func (m *MockGiteaInterface) GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositories", gitSource, gitOrgRef)
	ret0, _ := ret[0].(*[]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepositories indicates an expected call of GetRepositories
func (mr *MockGiteaInterfaceMockRecorder) GetRepositories(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositories", reflect.TypeOf((*MockGiteaInterface)(nil).GetRepositories), gitSource, gitOrgRef)
}

// GetOrganization mocks base method
func (m *MockGiteaInterface) GetOrganization(gitSource *model.GitSource, gitOrgRef string) *dto.OrganizationDto {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganization", gitSource, gitOrgRef)
	ret0, _ := ret[0].(*dto.OrganizationDto)
	return ret0
}

// GetOrganization indicates an expected call of GetOrganization
func (mr *MockGiteaInterfaceMockRecorder) GetOrganization(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganization", reflect.TypeOf((*MockGiteaInterface)(nil).GetOrganization), gitSource, gitOrgRef)
}

// CheckOrganizationExists mocks base method
func (m *MockGiteaInterface) CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckOrganizationExists", gitSource, gitOrgRef)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckOrganizationExists indicates an expected call of CheckOrganizationExists
func (mr *MockGiteaInterfaceMockRecorder) CheckOrganizationExists(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckOrganizationExists", reflect.TypeOf((*MockGiteaInterface)(nil).CheckOrganizationExists), gitSource, gitOrgRef)
}

// GetRepositoryTeams mocks base method
func (m *MockGiteaInterface) GetRepositoryTeams(gitSource *model.GitSource, gitOrgRef, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositoryTeams", gitSource, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(*[]dto.TeamResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepositoryTeams indicates an expected call of GetRepositoryTeams
func (mr *MockGiteaInterfaceMockRecorder) GetRepositoryTeams(gitSource, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositoryTeams", reflect.TypeOf((*MockGiteaInterface)(nil).GetRepositoryTeams), gitSource, gitOrgRef, repositoryRef)
}

// GetOrganizationTeams mocks base method
func (m *MockGiteaInterface) GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganizationTeams", gitSource, gitOrgRef)
	ret0, _ := ret[0].(*[]dto.TeamResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganizationTeams indicates an expected call of GetOrganizationTeams
func (mr *MockGiteaInterfaceMockRecorder) GetOrganizationTeams(gitSource, gitOrgRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganizationTeams", reflect.TypeOf((*MockGiteaInterface)(nil).GetOrganizationTeams), gitSource, gitOrgRef)
}

// GetTeamMembers mocks base method
func (m *MockGiteaInterface) GetTeamMembers(gitSource *model.GitSource, teamId int) (*[]dto.UserTeamResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamMembers", gitSource, teamId)
	ret0, _ := ret[0].(*[]dto.UserTeamResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeamMembers indicates an expected call of GetTeamMembers
func (mr *MockGiteaInterfaceMockRecorder) GetTeamMembers(gitSource, teamId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamMembers", reflect.TypeOf((*MockGiteaInterface)(nil).GetTeamMembers), gitSource, teamId)
}

// GetBranches mocks base method
func (m *MockGiteaInterface) GetBranches(gitSource *model.GitSource, gitOrgRef, repositoryRef string) map[string]bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranches", gitSource, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(map[string]bool)
	return ret0
}

// GetBranches indicates an expected call of GetBranches
func (mr *MockGiteaInterfaceMockRecorder) GetBranches(gitSource, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranches", reflect.TypeOf((*MockGiteaInterface)(nil).GetBranches), gitSource, gitOrgRef, repositoryRef)
}

// getBranches mocks base method
func (m *MockGiteaInterface) getBranches(gitSource *model.GitSource, gitOrgRef, repositoryRef string) (*[]gitea.BranchResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getBranches", gitSource, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(*[]gitea.BranchResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getBranches indicates an expected call of getBranches
func (mr *MockGiteaInterfaceMockRecorder) getBranches(gitSource, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getBranches", reflect.TypeOf((*MockGiteaInterface)(nil).getBranches), gitSource, gitOrgRef, repositoryRef)
}

// getRepositoryAgolaMetadata mocks base method
func (m *MockGiteaInterface) getRepositoryAgolaMetadata(gitSource *model.GitSource, gitOrgRef, repositoryRef, branchName string) (*[]gitea.MetadataResponseDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getRepositoryAgolaMetadata", gitSource, gitOrgRef, repositoryRef, branchName)
	ret0, _ := ret[0].(*[]gitea.MetadataResponseDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getRepositoryAgolaMetadata indicates an expected call of getRepositoryAgolaMetadata
func (mr *MockGiteaInterfaceMockRecorder) getRepositoryAgolaMetadata(gitSource, gitOrgRef, repositoryRef, branchName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getRepositoryAgolaMetadata", reflect.TypeOf((*MockGiteaInterface)(nil).getRepositoryAgolaMetadata), gitSource, gitOrgRef, repositoryRef, branchName)
}

// CheckRepositoryAgolaConfExists mocks base method
func (m *MockGiteaInterface) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, gitOrgRef, repositoryRef string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRepositoryAgolaConfExists", gitSource, gitOrgRef, repositoryRef)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRepositoryAgolaConfExists indicates an expected call of CheckRepositoryAgolaConfExists
func (mr *MockGiteaInterfaceMockRecorder) CheckRepositoryAgolaConfExists(gitSource, gitOrgRef, repositoryRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRepositoryAgolaConfExists", reflect.TypeOf((*MockGiteaInterface)(nil).CheckRepositoryAgolaConfExists), gitSource, gitOrgRef, repositoryRef)
}

// GetCommitMetadata mocks base method
func (m *MockGiteaInterface) GetCommitMetadata(gitSource *model.GitSource, gitOrgRef, repositoryRef, commitSha string) (*dto.CommitMetadataDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitMetadata", gitSource, gitOrgRef, repositoryRef, commitSha)
	ret0, _ := ret[0].(*dto.CommitMetadataDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitMetadata indicates an expected call of GetCommitMetadata
func (mr *MockGiteaInterfaceMockRecorder) GetCommitMetadata(gitSource, gitOrgRef, repositoryRef, commitSha interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitMetadata", reflect.TypeOf((*MockGiteaInterface)(nil).GetCommitMetadata), gitSource, gitOrgRef, repositoryRef, commitSha)
}
