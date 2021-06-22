package service

import (
	"encoding/json"
	"errors"
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
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func TestGetOrganizationsOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organizationsMock := test.MakeOrganizationList()

	db := mock_repository.NewMockDatabase(ctl)
	db.EXPECT().GetOrganizations().Return(organizationsMock, nil)

	serviceOrganization := OrganizationService{
		Db: db,
	}

	ts := httptest.NewServer(http.HandlerFunc(serviceOrganization.GetOrganizations))
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL)

	assert.Equal(t, err, nil)

	var organizations *[]model.Organization
	test.ParseBody(resp, organizations)
}

func TestAddExternalUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	commonMutex := utils.NewEventMutex()
	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}
	user := test.MakeUser()

	org := (*test.MakeOrganizationList())[0]
	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.Name).Return(true, nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(user)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

}
func TestAddExternalUserWhenOrganizationNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	commonMutex := utils.NewEventMutex()
	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}
	user := test.MakeUser()

	organizationRefTest := "testnotfound"
	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(nil, nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(user)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organizationRefTest, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "http StatusCode is not OK")
}

func TestRemoveExternalUserOk(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	commonMutex := utils.NewEventMutex()
	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}
	mail := "user@email.com"
	user := test.MakeUser()

	org := (*test.MakeOrganizationList())[0]
	org.ExternalUsers = make(map[string]bool)
	org.ExternalUsers[mail] = true

	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.Name).Return(true, nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationName}", serviceOrganization.RemoveExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mail})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode OK")
	exist := org.ExternalUsers[mail]
	assert.Check(t, !exist, "")
}

func TestRemoveExternalUserWhenAgolaRefNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	commonMutex := utils.NewEventMutex()
	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)

	serviceOrganization := OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}
	mail := "user@email.com"
	user := test.MakeUser()

	org := (*test.MakeOrganizationList())[0]
	org.ExternalUsers = make(map[string]bool)
	org.ExternalUsers[mail] = true

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, errors.New(string("someError")))

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.RemoveExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mail})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "http StatusCode is not OK")
}
