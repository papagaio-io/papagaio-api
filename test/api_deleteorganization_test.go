package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/service"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_agola"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func TestDeleteOrganizationOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex := utils.NewEventMutex()

	organization := (*MakeOrganizationList())[0]
	gitSource := (*MakeGitSourceMap())[organization.GitSourceName]

	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organization.AgolaOrganizationRef)).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(organization.GitSourceName)).Return(&gitSource, nil)
	agolaApi.EXPECT().DeleteOrganization(gomock.Any(), gomock.Any()).Return(nil)
	giteaApi.EXPECT().DeleteWebHook(gomock.Any(), gomock.Eq(organization.Name), gomock.Eq(organization.WebHookID)).Return(nil)
	db.EXPECT().DeleteOrganization(gomock.Eq(organization.AgolaOrganizationRef)).Return(nil)

	serviceOrganization := service.OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organization.AgolaOrganizationRef)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestDeleteOrganizationInternalOnly(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex := utils.NewEventMutex()

	organization := (*MakeOrganizationList())[0]
	gitSource := (*MakeGitSourceMap())[organization.GitSourceName]

	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organization.AgolaOrganizationRef)).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(organization.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().DeleteWebHook(gomock.Any(), gomock.Eq(organization.Name), gomock.Eq(organization.WebHookID)).Return(nil)
	db.EXPECT().DeleteOrganization(gomock.Eq(organization.AgolaOrganizationRef)).Return(nil)

	serviceOrganization := service.OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organization.AgolaOrganizationRef + "?internalonly")

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestDeleteOrganizationNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	commonMutex := utils.NewEventMutex()

	organizationRefTest := "testnotfound"

	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organizationRefTest)).Return(nil, nil)

	serviceOrganization := service.OrganizationService{
		Db:          db,
		CommonMutex: &commonMutex,
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organizationRefTest)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "http StatusCode is not correct")
}
