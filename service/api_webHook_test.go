package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/xanzy/go-gitlab"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/test"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_agola"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitlab"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	"wecode.sorint.it/opensource/papagaio-api/types"
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

	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gomock.Any()).Return("projectTestID", nil)
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
	user := test.MakeUser()

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
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(false, nil)
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
	user := test.MakeUser()

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
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	agolaApi.EXPECT().DeleteProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any()).Return(nil)
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
	user := test.MakeUser()

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
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gomock.Any()).Return("projectTestID", nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupBranchSynckMock(db, giteaApi, organization.GitPath, repositoryRef)

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

func TestRepositoryPushWithAgolaConfAndProjectNotExistsWithErrors(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	repositoryRef := "repositoryTest"

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

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

	// agola CreateProject error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gomock.Any()).Return("", errors.New("test error"))

	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	// SaveOrganization errpr

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gomock.Any()).Return("projectTestID", nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(errors.New("test error"))

	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestRepositoryPushWithAgolaConfAndProjectArchivied(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	repositoryRef := "repositoryTest"
	organization.Projects = make(map[string]model.Project)
	organization.Projects[repositoryRef] = model.Project{AgolaProjectRef: repositoryRef, Archivied: true, AgolaProjectID: "test"}

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
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().UnarchiveProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)).Return(nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupBranchSynckMock(db, giteaApi, organization.GitPath, repositoryRef)

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

func TestRepositoryPushWithAgolaConfAndProjectArchiviedWithErrors(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	repositoryRef := "repositoryTest"
	organization.Projects = make(map[string]model.Project)
	organization.Projects[repositoryRef] = model.Project{AgolaProjectRef: repositoryRef, Archivied: true, AgolaProjectID: "test"}

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: repositoryRef},
		Action:     "",
	}

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

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

	// agola UnarchiveProject error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().UnarchiveProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)).Return(errors.New("test error"))

	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	// SaveOrganization error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().UnarchiveProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)).Return(nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(errors.New("test error"))

	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestRepositoryPushWithAgolaConfAndProjectNotArchivied(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	repositoryRef := "repositoryTest"
	organization.Projects = make(map[string]model.Project)
	organization.Projects[repositoryRef] = model.Project{AgolaProjectRef: repositoryRef, Archivied: false, AgolaProjectID: "test"}

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
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)

	setupBranchSynckMock(db, giteaApi, organization.GitPath, repositoryRef)

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
	user := test.MakeUser()

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
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(false, nil)

	setupBranchSynckMock(db, giteaApi, organization.GitPath, repositoryRef)

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

func TestRepositoryPushWithoutAgolaConfAndProjectNotArchivied(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

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
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(false, nil)
	agolaApi.EXPECT().ArchiveProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)).Return(nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupBranchSynckMock(db, giteaApi, organization.GitPath, repositoryRef)

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

func TestRepositoryPushWithoutAgolaConfAndProjectNotArchiviedWithErrors(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

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

	// agola ArchiveProject error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(false, nil)
	agolaApi.EXPECT().ArchiveProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)).Return(errors.New("test error"))

	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	// SaveOrganization error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(false, nil)
	agolaApi.EXPECT().ArchiveProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)).Return(nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(errors.New("test error"))

	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestRepositoryGitlabPushWithAgolaConfAndProjectNotExists(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	organization.GitSourceName = "gitlab"
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()
	user.GitSourceName = "gitlab"

	repositoryRef := "repositoryTest"

	webHookMessage := gitlab.PushEvent{
		ProjectID:   1,
		Repository:  &gitlab.Repository{Name: repositoryRef},
		CheckoutSHA: "test",
	}

	db := mock_repository.NewMockDatabase(ctl)
	gitlabApi := mock_gitlab.NewMockGitlabInterface(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	commonMutex := utils.NewEventMutex()

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	gitlabApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gomock.Any()).Return("projectTestID", nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupBranchSynckGitlabMock(db, gitlabApi, organization.GitPath, repositoryRef)

	serviceWebHook := WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GitlabApi: gitlabApi},
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

func TestRepositoryCreatedWithAgolaConfigWithErrors(t *testing.T) {
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

	user := test.MakeUser()

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

	// organization not found

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(nil, nil)

	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	// gitSource not found

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(nil, nil)

	data, _ = json.Marshal(webHookMessage)
	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	// repository behaviour exclude

	organizationTest := organization
	organizationTest.BehaviourType = types.Regex
	organizationTest.BehaviourExclude = repositoryRef

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organizationTest, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)

	data, _ = json.Marshal(webHookMessage)
	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity)

	// user not found

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(organization.UserIDConnected).Return(nil, nil)

	data, _ = json.Marshal(webHookMessage)
	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	// agola CreateProject error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(organization.UserIDConnected).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gomock.Any()).Return("", errors.New("error test"))

	data, _ = json.Marshal(webHookMessage)
	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	// SaveOrganization error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(organization.UserIDConnected).Return(user, nil)
	giteaApi.EXPECT().CheckRepositoryAgolaConfExists(gomock.Any(), gomock.Any(), organization.GitPath, webHookMessage.Repository.Name).Return(true, nil)
	agolaApi.EXPECT().CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any(), gitSource.AgolaRemoteSource, gomock.Any()).Return("projectTestID", nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(errors.New("test error"))

	data, _ = json.Marshal(webHookMessage)
	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestRepositoryDeletedWithErrors(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

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

	// repository not found

	organizationTest := organization
	organizationTest.Projects = make(map[string]model.Project)

	db.EXPECT().GetOrganizationByAgolaRef(organizationTest.AgolaOrganizationRef).Return(&organizationTest, nil)
	db.EXPECT().GetGitSourceByName(organizationTest.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)

	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organizationTest.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity)

	// agola DeleteProject error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	agolaApi.EXPECT().DeleteProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any()).Return(errors.New("test error"))

	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	// SaveOrganization error

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	agolaApi.EXPECT().DeleteProject(gomock.Any(), utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gomock.Any()).Return(nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(errors.New("test error"))

	requestBody = strings.NewReader(string(data))
	resp, err = client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func setupBranchSynckMock(db *mock_repository.MockDatabase, giteaApi *mock_gitea.MockGiteaInterface, organizationName string, repositoryName string) {
	branches := make(map[string]bool)
	branches["master"] = true

	giteaApi.EXPECT().GetBranches(gomock.Any(), gomock.Any(), organizationName, repositoryName).Return(branches, nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)
}

func setupBranchSynckGitlabMock(db *mock_repository.MockDatabase, gitlabApi *mock_gitlab.MockGitlabInterface, organizationName string, repositoryName string) {
	branches := make(map[string]bool)
	branches["master"] = true

	gitlabApi.EXPECT().GetBranches(gomock.Any(), gomock.Any(), organizationName, repositoryName).Return(branches, nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)
}
