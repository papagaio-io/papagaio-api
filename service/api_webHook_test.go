package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/test"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_agola"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func TestRepositoryCreatedWithAgolaConfigOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	repositoryRef := "repositoryTest"

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "created",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), organization.Name, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gitSource.AgolaToken).Return("projectTestID", nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	serviceWebHook := WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceWebHook.WebHookOrganization)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	project, exists := organization.Projects[repositoryRef]
	assert.Check(t, exists)

	assert.Equal(t, project.AgolaProjectRef, repositoryRef)
	assert.Check(t, !project.Archivied)
}

func TestRepositoryCreatedWithoutAgolaConfigOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	repositoryRef := "repositoryTest"

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "created",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), organization.Name, webHookMessage.Repository.Name).Return(false, nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	serviceWebHook := WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceWebHook.WebHookOrganization)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	project, exists := organization.Projects[repositoryRef]
	assert.Check(t, exists)

	assert.Equal(t, project.AgolaProjectRef, repositoryRef)
	assert.Check(t, project.Archivied)
}

func TestRepositoryDeletedOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	repositoryRef := "repositoryTest"
	organization.Projects = make(map[string]model.Project)
	organization.Projects[repositoryRef] = model.Project{AgolaProjectRef: repositoryRef}

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "deleted",
	}

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	agolaApi.EXPECT().DeleteProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gitSource.AgolaToken).Return(nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	serviceWebHook := WebHookService{
		Db:          db,
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceWebHook.WebHookOrganization)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	_, exists := organization.Projects[repositoryRef]
	assert.Check(t, !exists)
}

func TestRepositoryPushWithAgolaConfAndProjectNotExists(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	repositoryRef := "repositoryTest"

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), organization.Name, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gitSource.AgolaToken).Return("projectTestID", nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupBranchSynckMock(db, giteaApi, organization.Name, repositoryRef)

	serviceWebHook := WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceWebHook.WebHookOrganization)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	project, exists := organization.Projects[repositoryRef]
	assert.Check(t, exists)

	assert.Equal(t, project.AgolaProjectRef, repositoryRef)
	assert.Check(t, !project.Archivied)
	assert.Equal(t, project.Branchs["master"].Name, "master")
}

func TestRepositoryPushWithAgolaConfAndProjectArchivied(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	repositoryRef := "repositoryTest"
	organization.Projects = make(map[string]model.Project)
	organization.Projects[repositoryRef] = model.Project{AgolaProjectRef: repositoryRef, Archivied: true}

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), organization.Name, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().UnarchiveProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)).Return(nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupBranchSynckMock(db, giteaApi, organization.Name, repositoryRef)

	serviceWebHook := WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceWebHook.WebHookOrganization)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	project, exists := organization.Projects[repositoryRef]
	assert.Check(t, exists)

	assert.Equal(t, project.AgolaProjectRef, repositoryRef)
	assert.Check(t, !project.Archivied)
	assert.Equal(t, project.Branchs["master"].Name, "master")
}

func TestRepositoryPushWithAgolaConfAndProjectNotArchivied(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	repositoryRef := "repositoryTest"
	organization.Projects = make(map[string]model.Project)
	organization.Projects[repositoryRef] = model.Project{AgolaProjectRef: repositoryRef, Archivied: false}

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), organization.Name, webHookMessage.Repository.Name).Return(true, nil)

	setupBranchSynckMock(db, giteaApi, organization.Name, repositoryRef)

	serviceWebHook := WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceWebHook.WebHookOrganization)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	project, exists := organization.Projects[repositoryRef]
	assert.Check(t, exists)

	assert.Equal(t, project.AgolaProjectRef, repositoryRef)
	assert.Check(t, !project.Archivied)
	assert.Equal(t, project.Branchs["master"].Name, "master")
}

func TestRepositoryPushWithoutAgolaConfAndProjectArchivied(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	repositoryRef := "repositoryTest"
	organization.Projects = make(map[string]model.Project)
	organization.Projects[repositoryRef] = model.Project{AgolaProjectRef: repositoryRef, Archivied: true}

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), organization.Name, webHookMessage.Repository.Name).Return(false, nil)

	setupBranchSynckMock(db, giteaApi, organization.Name, repositoryRef)

	serviceWebHook := WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceWebHook.WebHookOrganization)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	project, exists := organization.Projects[repositoryRef]
	assert.Check(t, exists)

	assert.Equal(t, project.AgolaProjectRef, repositoryRef)
	assert.Check(t, project.Archivied)
	assert.Equal(t, project.Branchs["master"].Name, "master")
}

func TestRepositoryPushWithoutAgolaConfAndProjectNotExists(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	repositoryRef := "repositoryTest"
	organization.Projects = make(map[string]model.Project)
	organization.Projects[repositoryRef] = model.Project{AgolaProjectRef: repositoryRef, Archivied: false}

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), organization.Name, webHookMessage.Repository.Name).Return(false, nil)
	agolaApi.EXPECT().ArchiveProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)).Return(nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupBranchSynckMock(db, giteaApi, organization.Name, repositoryRef)

	serviceWebHook := WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceWebHook.WebHookOrganization)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	project, exists := organization.Projects[repositoryRef]
	assert.Check(t, exists)

	assert.Equal(t, project.AgolaProjectRef, repositoryRef)
	assert.Check(t, project.Archivied)
	assert.Equal(t, project.Branchs["master"].Name, "master")
}

func setupBranchSynckMock(db *mock_repository.MockDatabase, giteaApi *mock_gitea.MockGiteaInterface, organizationName string, repositoryName string) {
	branches := make(map[string]bool)
	branches["master"] = true

	giteaApi.EXPECT().GetBranches(gomock.Any(), organizationName, repositoryName).Return(branches)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)
}