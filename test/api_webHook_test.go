package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/service"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_agola"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func TestRepositoryCreatedWithAgolaConfigOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organization := (*MakeOrganizationList())[0]
	gitSource := (*MakeGitSourceMap())[organization.GitSourceName]

	webHookMessage := dto.WebHookDto{
		Repository: dto.RepositoryDto{ID: 1, Name: "repositoryTest"},
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

	serviceWebHook := service.WebHookService{
		Db:          db,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}

	ts := httptest.NewServer(http.HandlerFunc(serviceWebHook.WebHookOrganization))
	defer ts.Close()

	fmt.Println("--------------------url:", ts.URL+"/"+organization.AgolaOrganizationRef)

	client := ts.Client()
	data, _ := json.Marshal(webHookMessage)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organization.AgolaOrganizationRef, "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
