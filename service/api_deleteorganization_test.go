package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/test"
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

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organization.AgolaOrganizationRef)).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(organization.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), organization.GitPath).Return(true, nil)
	db.EXPECT().GetUserByUserId(organization.UserIDConnected).Return(user, nil)
	agolaApi.EXPECT().DeleteOrganization(gomock.Any(), gomock.Any()).Return(nil)
	giteaApi.EXPECT().DeleteWebHook(gomock.Any(), gomock.Any(), gomock.Eq(organization.GitPath), gomock.Eq(organization.WebHookID)).Return(nil)
	db.EXPECT().DeleteOrganization(gomock.Eq(organization.AgolaOrganizationRef)).Return(nil)

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := test.SetupBaseRouter(user)
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

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organization.AgolaOrganizationRef)).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(organization.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), organization.GitPath).Return(true, nil)
	db.EXPECT().GetUserByUserId(organization.UserIDConnected).Return(user, nil)
	giteaApi.EXPECT().DeleteWebHook(gomock.Any(), gomock.Any(), gomock.Eq(organization.GitPath), gomock.Eq(organization.WebHookID)).Return(nil)
	db.EXPECT().DeleteOrganization(gomock.Eq(organization.AgolaOrganizationRef)).Return(nil)

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := test.SetupBaseRouter(user)
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
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organizationRefTest)).Return(nil, nil)

	serviceOrganization := OrganizationService{
		Db:          db,
		CommonMutex: &commonMutex,
	}

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organizationRefTest)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "http StatusCode is not correct")
}
func TestDeleteOrganizationInternalonlyInvalidParam(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex := utils.NewEventMutex()

	organization := (*test.MakeOrganizationList())[0]
	user := test.MakeUser()

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organization.AgolaOrganizationRef + "?internalonly=invalid")

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity, "http StatusCode is not OK")
}

func TestDeleteOrganizationWhenAgolaRefNotFoundOnDB(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex := utils.NewEventMutex()

	organization := (*test.MakeOrganizationList())[0]
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organization.AgolaOrganizationRef)).Return(&organization, errors.New(string("someError")))

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organization.AgolaOrganizationRef)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "Organization not found on db")
}

func TestDeleteOrganizationWhenGitSourceNotFoundOnDB(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex := utils.NewEventMutex()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organization.AgolaOrganizationRef)).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(gomock.Any()).Return(&gitSource, errors.New(string("someError")))

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organization.AgolaOrganizationRef)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "Organization not found on db")
}

func TestDeleteOrganizationWhenDeletingOrganizationInAgola(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex := utils.NewEventMutex()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organization.AgolaOrganizationRef)).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(gomock.Any()).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), organization.GitPath).Return(true, nil)
	db.EXPECT().GetUserByUserId(organization.UserIDConnected).Return(user, nil)
	agolaApi.EXPECT().DeleteOrganization(gomock.Any(), gomock.Any()).Return(errors.New(string("someError")))

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organization.AgolaOrganizationRef)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "Organization not found in Agola")
}

func TestDeleteOrganizationWhenInternalOnlyDeletingOrganizationError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex := utils.NewEventMutex()

	organization := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Eq(organization.AgolaOrganizationRef)).Return(&organization, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(organization.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), organization.GitPath).Return(true, nil)
	db.EXPECT().GetUserByUserId(organization.UserIDConnected).Return(user, nil)
	agolaApi.EXPECT().DeleteOrganization(gomock.Any(), gomock.Any()).Return(nil)
	giteaApi.EXPECT().DeleteWebHook(gomock.Any(), gomock.Any(), gomock.Eq(organization.GitPath), gomock.Eq(organization.WebHookID)).Return(nil)
	db.EXPECT().DeleteOrganization(gomock.Eq(organization.AgolaOrganizationRef)).Return(errors.New(string("someError")))

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationRef}", serviceOrganization.DeleteOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + organization.AgolaOrganizationRef)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "Organization not found in Agola")
}
