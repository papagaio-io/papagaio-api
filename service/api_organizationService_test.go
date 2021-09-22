package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/agola"
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
	mailTest := "emailtest@sorint.it"
	org.ExternalUsers = make(map[string]bool)

	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(true, nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mailTest})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

}

func TestAddExternalUserWhenUserNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	commonMutex := utils.NewEventMutex()
	db := mock_repository.NewMockDatabase(ctl)

	serviceOrganization := OrganizationService{
		Db:          db,
		CommonMutex: &commonMutex,
	}
	user := test.MakeUser()
	organizationRefTest := "anyOrganization"
	db.EXPECT().GetUserByUserId(gomock.Any()).Return(nil, nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(user)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+organizationRefTest, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not OK")
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

func TestAddExternalUserWhenGitSourceNotFound(t *testing.T) {
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
	mailTest := "emailtest@sorint.it"
	org.ExternalUsers = make(map[string]bool)

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(nil, nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mailTest})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not OK")
}

func TestAddExternalUserWhenUserIsNotOwner(t *testing.T) {
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
	mailTest := "emailtest@sorint.it"
	org.ExternalUsers = make(map[string]bool)

	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(false, nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mailTest})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)

	var responseDto dto.ExternalUsersDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
	assert.Equal(t, responseDto.ErrorCode, dto.UserNotOwnerError, "ErrorCode is not correct")
}

func TestAddExternalUserWhenEmailIsNotValid(t *testing.T) {
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
	mailTest := "emailtestsorint.it"
	org.ExternalUsers = make(map[string]bool)

	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(true, nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mailTest})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)

	var responseDto dto.ExternalUsersDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
	assert.Equal(t, responseDto.ErrorCode, dto.InvalidEmail, "ErrorCode is not correct")
}

func TestAddExternalUserWhenEmailAlreadyExists(t *testing.T) {
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
	mailTest := "emailtest@sorint.it"
	org.ExternalUsers = make(map[string]bool)
	org.ExternalUsers[mailTest] = true

	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(true, nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mailTest})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)

	var responseDto dto.ExternalUsersDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
	assert.Equal(t, responseDto.ErrorCode, dto.EmailAlreadyExists, "ErrorCode is not correct")
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
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(true, nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationName}", serviceOrganization.RemoveExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mail})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode not correct")
	exist := org.ExternalUsers[mail]
	assert.Check(t, !exist, "")
}

func TestRemoveExternalUserWhenUserNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	commonMutex := utils.NewEventMutex()
	db := mock_repository.NewMockDatabase(ctl)

	serviceOrganization := OrganizationService{
		Db:          db,
		CommonMutex: &commonMutex,
	}
	mail := "user@email.com"
	user := test.MakeUser()

	org := (*test.MakeOrganizationList())[0]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(nil, nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationName}", serviceOrganization.RemoveExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mail})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "ErrorCode is not correct")
}

func TestRemoveExternalUserWhenGitSourceNotFound(t *testing.T) {
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

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(nil, nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationName}", serviceOrganization.RemoveExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mail})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)

	exist := org.ExternalUsers[mail]
	assert.Check(t, exist, "")

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode not correct")
}

func TestRemoveExternalUserNotOwner(t *testing.T) {
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
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(false, nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationName}", serviceOrganization.RemoveExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mail})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)

	exist := org.ExternalUsers[mail]
	assert.Check(t, exist, "")

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode not correct")
	assert.Equal(t, responseDto.ErrorCode, dto.UserNotOwnerError, "ErrorCode is not correct")
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
func TestRemoveExternalUserWhenEmailNotFound(t *testing.T) {
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

	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(true, nil)

	router := test.SetupBaseRouter(user)
	router.HandleFunc("/{organizationName}", serviceOrganization.RemoveExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(dto.ExternalUserDto{Email: mail})
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)

	var responseDto dto.ExternalUsersDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode not correct")
	assert.Equal(t, responseDto.ErrorCode, dto.EmailNotFound, "ErrorCode is not correct")
}
func TestGetAgolaOrganizationsOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	agolaOrganizations := []agola.OrganizationDto{
		{
			ID: "1", Name: "test", Visibility: "public",
		},
	}

	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	db := mock_repository.NewMockDatabase(ctl)

	agolaApi.EXPECT().GetOrganizations().Return(&agolaOrganizations, nil)
	db.EXPECT().GetOrganizationByAgolaRef(agolaOrganizations[0].Name).Return(nil, nil)

	serviceOrganization := OrganizationService{
		AgolaApi: agolaApi,
		Db:       db,
	}

	ts := httptest.NewServer(http.HandlerFunc(serviceOrganization.GetAgolaOrganizations))
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL)

	assert.Equal(t, err, nil)

	var organizations []string
	test.ParseBody(resp, &organizations)

	fmt.Println("result:", organizations)

	assert.Check(t, len(organizations) == 1, "organizations is empty")
	assert.Equal(t, organizations[0], "test")
}

func TestGetExternalUsersOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)

	serviceOrganization := OrganizationService{
		Db:         db,
		AgolaApi:   agolaApi,
		GitGateway: &git.GitGateway{GiteaApi: giteaApi},
	}
	user := test.MakeUser()

	org := (*test.MakeOrganizationList())[0]
	mailTest := "email1@sorint.it"
	org.ExternalUsers = make(map[string]bool)
	org.ExternalUsers[mailTest] = true

	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(true, nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.GetExternalUsers)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + org.AgolaOrganizationRef)

	var dtoResponse = dto.ExternalUsersDto{}

	test.ParseBody(resp, &dtoResponse)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
	assert.Equal(t, len(*dtoResponse.EmailList), len(org.ExternalUsers))
	assert.Equal(t, mailTest, (*dtoResponse.EmailList)[0])
	assert.Equal(t, dtoResponse.ErrorCode, dto.NoError)
}
func TestGetExternalUsersWhenUserNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)

	serviceOrganization := OrganizationService{
		Db: db,
	}
	org := (*test.MakeOrganizationList())[0]
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(gomock.Any()).Return(nil, nil)
	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.GetExternalUsers)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + org.AgolaOrganizationRef)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not OK")
}

func TestGetExternalUserWhenAgolarefNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)

	serviceOrganization := OrganizationService{
		Db:       db,
		AgolaApi: agolaApi,
	}
	org := (*test.MakeOrganizationList())[0]
	org.AgolaOrganizationRef = "orgNotFound"
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(gomock.Any()).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(nil, nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.GetExternalUsers)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + org.AgolaOrganizationRef)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "http StatusCode is not OK")
}
func TestGetExternalUserWhenGitSourceNotFound(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)

	serviceOrganization := OrganizationService{
		Db:       db,
		AgolaApi: agolaApi,
	}
	org := (*test.MakeOrganizationList())[0]
	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(gomock.Any()).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Any()).Return(nil, nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.GetExternalUsers)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + org.AgolaOrganizationRef)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not OK")
}

func TestGetExternalUserWhenIsnotOwner(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)

	serviceOrganization := OrganizationService{
		Db:         db,
		AgolaApi:   agolaApi,
		GitGateway: &git.GitGateway{GiteaApi: giteaApi},
	}
	org := (*test.MakeOrganizationList())[0]
	user := test.MakeUser()
	gitSource := (*test.MakeGitSourceMap())[org.GitSourceName]

	db.EXPECT().GetUserByUserId(gomock.Any()).Return(user, nil)
	db.EXPECT().GetOrganizationByAgolaRef(gomock.Any()).Return(&org, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(org.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), org.GitPath).Return(false, nil)

	router := test.SetupBaseRouter(user)

	router.HandleFunc("/{organizationName}", serviceOrganization.GetExternalUsers)
	ts := httptest.NewServer(router)

	client := ts.Client()

	resp, err := client.Get(ts.URL + "/" + org.AgolaOrganizationRef)

	var dtoResponse = dto.ExternalUsersDto{}
	test.ParseBody(resp, &dtoResponse)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode not correct")
	assert.Equal(t, dtoResponse.ErrorCode, dto.UserNotOwnerError)
}
