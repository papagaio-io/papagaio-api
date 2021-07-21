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

func newString(s string) *string {
	return &s
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
		GitType:               "gitea",
		GitAPIURL:             newString("http://test"),
		GitClientID:           "test",
		GitClientSecret:       "test",
		AgolaRemoteSourceName: newString("test"),
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

func TestRemoveGitsourcesOK(t *testing.T) {
	setupGitsourceMock(t)

	gitSource := (*test.MakeGitSourceMap())["gitea"]

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	db.EXPECT().GetOrganizationsByGitSource(gitSource.Name).Return(nil, nil)
	db.EXPECT().GetUsersIDByGitSourceName(gitSource.Name).Return(make([]uint64, 0), nil)
	agolaApiInt.EXPECT().DeleteRemotesource(gitSource.AgolaRemoteSource).Return(nil)
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

func TestUpdateGitsourcesOK(t *testing.T) {
	setupGitsourceMock(t)

	giteaType := types.Gitea

	gitSource := (*test.MakeGitSourceMap())["gitea"]
	reqDto := dto.UpdateGitSourceRequestDto{
		GitType:           &giteaType,
		GitAPIURL:         newString("http://test"),
		GitClientID:       newString("test"),
		GitClientSecret:   newString("test"),
		AgolaRemoteSource: newString("test"),
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
