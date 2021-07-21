package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/test"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_agola"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	"wecode.sorint.it/opensource/papagaio-api/types"
	"wecode.sorint.it/opensource/papagaio-api/utils"

	agolaDto "wecode.sorint.it/opensource/papagaio-api/api/agola"
	gitDto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
)

var serviceGitsource GitSourceService

func setupGitsourceMock(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db = mock_repository.NewMockDatabase(ctl)
	agolaApiInt = mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi = mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex = utils.NewEventMutex()

	serviceGitsource = GitSourceService{
		Db:         db,
		AgolaApi:   agolaApiInt,
		GitGateway: &git.GitGateway{GiteaApi: giteaApi},
	}
}

func TestGetGitsourcesOK(t *testing.T) {
	setupGitsourceMock(t)

	gitSource := (*test.MakeGitSourceMap())["gitea"]
	gitSources := make([]model.GitSource, 0)
	gitSources = append(gitSources, gitSource)

	db.EXPECT().GetGitSources().Return(&gitSources, nil)

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/gitsources", serviceGitsource.GetGitSources)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/gitsources")

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestAddGitsourcesOK(t *testing.T) {
	setupGitsourceMock(t)

	reqDto := dto.CreateGitSourceRequestDto{
		Name:                  "test",
		GitType:               "github",
		GitClientID:           "test",
		GitClientSecret:       "test",
		AgolaRemoteSourceName: utils.NewString("test"),
	}

	db.EXPECT().GetGitSourceByName(reqDto.Name).Return(nil, nil)
	db.EXPECT().SaveGitSource(gomock.Any()).Return(nil)

	data, _ := json.Marshal(reqDto)
	requestBody := strings.NewReader(string(data))

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/gitsource", serviceGitsource.AddGitSource)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Post(ts.URL+"/gitsource", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestAddGitsourcesWithCreateRemotesourceOK(t *testing.T) {
	setupGitsourceMock(t)

	reqDto := dto.CreateGitSourceRequestDto{
		Name:              "test",
		GitType:           "github",
		GitClientID:       "test",
		GitClientSecret:   "test",
		AgolaClientID:     utils.NewString("test"),
		AgolaClientSecret: utils.NewString("test"),
	}

	remoteSources := make([]agolaDto.RemoteSourceDto, 0)
	remoteSources = append(remoteSources, agolaDto.RemoteSourceDto{Name: reqDto.Name})

	db.EXPECT().GetGitSourceByName(reqDto.Name).Return(nil, nil)
	agolaApiInt.EXPECT().GetRemoteSources().Return(&remoteSources, nil)
	agolaApiInt.EXPECT().CreateRemoteSource(reqDto.Name+"0", string(reqDto.GitType), "https://api.github.com", *reqDto.AgolaClientID, *reqDto.AgolaClientSecret).Return(nil)
	db.EXPECT().SaveGitSource(gomock.Any()).Return(nil)

	data, _ := json.Marshal(reqDto)
	requestBody := strings.NewReader(string(data))

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/gitsource", serviceGitsource.AddGitSource)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Post(ts.URL+"/gitsource", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestAddGitsourcesNotValid(t *testing.T) {
	setupGitsourceMock(t)

	reqDto := dto.CreateGitSourceRequestDto{
		Name:            "test",
		GitType:         "gitea",
		GitClientID:     "test",
		GitClientSecret: "test",
	}

	db.EXPECT().GetGitSourceByName(reqDto.Name).Return(nil, nil)
	db.EXPECT().SaveGitSource(gomock.Any()).Return(nil)

	data, _ := json.Marshal(reqDto)
	requestBody := strings.NewReader(string(data))

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/gitsource", serviceGitsource.AddGitSource)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Post(ts.URL+"/gitsource", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity, "http StatusCode is not correct")
}

func TestRemoveGitsourcesOK(t *testing.T) {
	setupGitsourceMock(t)

	gitSource := (*test.MakeGitSourceMap())["gitea"]

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	db.EXPECT().GetOrganizationsByGitSource(gitSource.Name).Return(nil, nil)
	db.EXPECT().GetUsersIDByGitSourceName(gitSource.Name).Return(make([]uint64, 0), nil)
	db.EXPECT().DeleteGitSource(gitSource.Name).Return(nil)

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/gitsource/{gitSourceName}", serviceGitsource.RemoveGitSource)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/gitsource/" + gitSource.Name)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestRemoveGitsourcesWithDeleteRemotesource(t *testing.T) {
	setupGitsourceMock(t)

	gitSource := (*test.MakeGitSourceMap())["gitea"]

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	db.EXPECT().GetOrganizationsByGitSource(gitSource.Name).Return(nil, nil)
	db.EXPECT().GetUsersIDByGitSourceName(gitSource.Name).Return([]uint64{1}, nil)
	db.EXPECT().DeleteUser(uint64(1)).Return(nil)
	agolaApiInt.EXPECT().DeleteRemotesource(gitSource.AgolaRemoteSource).Return(nil)
	db.EXPECT().DeleteGitSource(gitSource.Name).Return(nil)

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/gitsource/{gitSourceName}", serviceGitsource.RemoveGitSource)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/gitsource/" + gitSource.Name + "?deleteremotesource")

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestUpdateGitsourcesOK(t *testing.T) {
	setupGitsourceMock(t)

	giteaType := types.Gitea

	gitSource := (*test.MakeGitSourceMap())["gitea"]
	reqDto := dto.UpdateGitSourceRequestDto{
		GitType:           &giteaType,
		GitAPIURL:         utils.NewString("http://test"),
		GitClientID:       utils.NewString("test"),
		GitClientSecret:   utils.NewString("test"),
		AgolaRemoteSource: utils.NewString("test"),
	}

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	db.EXPECT().SaveGitSource(gomock.Any()).Return(nil)

	data, _ := json.Marshal(reqDto)
	requestBody := strings.NewReader(string(data))

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/gitsource/{gitSourceName}", serviceGitsource.UpdateGitSource)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Post(ts.URL+"/gitsource/"+gitSource.Name, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestGetGitOrganizationsOK(t *testing.T) {
	setupGitsourceMock(t)

	user := test.MakeUser()
	gitSource := (*test.MakeGitSourceMap())[user.GitSourceName]

	organizations := make([]gitDto.OrganizationDto, 0)
	organizations = append(organizations, gitDto.OrganizationDto{
		ID:   1,
		Name: "test1",
	})
	organizations = append(organizations, gitDto.OrganizationDto{
		ID:   2,
		Name: "test2",
	})

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	giteaApi.EXPECT().GetOrganizations(gomock.Any(), gomock.Any()).Return(&organizations, nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/gitorganizations", serviceGitsource.GetGitOrganizations)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/gitorganizations")

	var responseDto []gitDto.OrganizationDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
	assert.Check(t, len(responseDto) == 2)
}
