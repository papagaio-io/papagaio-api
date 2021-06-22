package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	gitDto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/test"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_agola"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	"wecode.sorint.it/opensource/papagaio-api/types"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

var organizationReqDto dto.CreateOrganizationRequestDto
var commonMutex utils.CommonMutex
var giteaApi *mock_gitea.MockGiteaInterface
var agolaApiInt *mock_agola.MockAgolaApiInterface
var db *mock_repository.MockDatabase
var gitSource model.GitSource
var serviceOrganization OrganizationService
var organizationList []model.Organization

func setupMock(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db = mock_repository.NewMockDatabase(ctl)
	agolaApiInt = mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi = mock_gitea.NewMockGiteaInterface(ctl)
	commonMutex = utils.NewEventMutex()

	organizationReqDto = dto.CreateOrganizationRequestDto{
		Name:          "Test",
		AgolaRef:      "Test",
		Visibility:    types.Public,
		BehaviourType: types.None,
	}

	user := test.MakeUser()

	gitSource = (*test.MakeGitSourceMap())[user.GitSourceName]

	serviceOrganization = OrganizationService{
		Db:          db,
		AgolaApi:    agolaApiInt,
		GitGateway:  &git.GitGateway{GiteaApi: giteaApi},
		CommonMutex: &commonMutex,
	}
	organization := (*test.MakeOrganizationList())[0]
	insertRunsData(&organization)
	organizationList = make([]model.Organization, 0)
	organizationList = append(organizationList, organization)
}

func setupRouter(user *model.User) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", serviceOrganization.CreateOrganization)
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if user == nil {
				ctx = context.WithValue(ctx, "admin", true)
			} else {
				ctx = context.WithValue(ctx, "admin", false)
				ctx = context.WithValue(ctx, controller.XAuthUserId, user.ID)
			}

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	return router
}

func TestCreateOrganizationOK(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()

	db.EXPECT().GetUserByUserId(user.ID).Return(user, nil)
	db.EXPECT().GetOrganizationsByGitSource(user.GitSourceName).Return(&organizationList, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(user.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true)
	giteaApi.EXPECT().IsUserOwner(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true, nil)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(nil, nil)
	giteaApi.EXPECT().CreateWebHook(gomock.Any(), gomock.Any(), organizationReqDto.Name, organizationReqDto.AgolaRef).Return(1, nil)
	agolaApiInt.EXPECT().CheckOrganizationExists(gomock.Any()).Return(false, "")
	agolaApiInt.EXPECT().CreateOrganization(gomock.Any(), organizationReqDto.Visibility).Return("123456", nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupSynkMembersUserTestMocks(agolaApiInt, giteaApi, organizationReqDto.Name, gitSource.AgolaRemoteSource)
	setupCheckoutAllGitRepositoryEmptyMocks(giteaApi, organizationReqDto.Name)

	ts := httptest.NewServer(setupRouter(user))

	client := ts.Client()

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, responseDto.ErrorCode, dto.NoError, "ErrorCode is not correct")
	assert.Check(t, strings.Contains(responseDto.OrganizationURL, "/org/"+organizationReqDto.AgolaRef), "OrganizationURL is not correct")
}

func TestCreateOrganizationJustExistsInPapagaio(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()

	organizationModel := model.Organization{Name: organizationReqDto.Name, AgolaOrganizationRef: organizationReqDto.AgolaRef}

	db.EXPECT().GetGitSourceByName(gomock.Eq(user.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(&organizationModel, nil)

	ts := httptest.NewServer(setupRouter(user))

	client := ts.Client()

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, responseDto.ErrorCode, dto.PapagaioOrganizationExistsError, "ErrorCode is not correct")
}

func TestCreateOrganizationJustExistsInAgola(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()

	db.EXPECT().GetOrganizationsByGitSource(user.GitSourceName).Return(&organizationList, nil)
	db.EXPECT().GetGitSourceByName(gomock.Eq(user.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(nil, nil)
	giteaApi.EXPECT().CreateWebHook(gomock.Any(), organizationReqDto.Name, gomock.Any(), organizationReqDto.AgolaRef).Return(1, nil)
	agolaApiInt.EXPECT().CheckOrganizationExists(gomock.Any()).Return(true, "test123456")
	giteaApi.EXPECT().DeleteWebHook(gomock.Any(), gomock.Any(), organizationReqDto.Name, 1).Return(nil)

	ts := httptest.NewServer(setupRouter(user))

	client := ts.Client()

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, responseDto.ErrorCode, dto.AgolaOrganizationExistsError, "ErrorCode is not correct")
}

func TestCreateOrganizationGitOrganizationNotFound(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()

	db.EXPECT().GetGitSourceByName(gomock.Eq(user.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(false)

	ts := httptest.NewServer(setupRouter(user))

	client := ts.Client()

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, responseDto.ErrorCode, dto.GitOrganizationNotFoundError, "ErrorCode is not correct")
}

func TestCreateOrganizationJustExistsInAgolaForce(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()

	db.EXPECT().GetOrganizationsByGitSource(user.GitSourceName).Return(&organizationList, nil)
	db.EXPECT().GetGitSourceByName(user.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(nil, nil)
	giteaApi.EXPECT().CreateWebHook(gomock.Any(), organizationReqDto.Name, gomock.Any(), organizationReqDto.AgolaRef).Return(1, nil)
	agolaApiInt.EXPECT().CheckOrganizationExists(gomock.Any()).Return(true, "test123456")
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupSynkMembersUserTestMocks(agolaApiInt, giteaApi, organizationReqDto.Name, gitSource.AgolaRemoteSource)
	setupCheckoutAllGitRepositoryEmptyMocks(giteaApi, organizationReqDto.Name)

	ts := httptest.NewServer(setupRouter(user))

	client := ts.Client()

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/?force", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, responseDto.ErrorCode, dto.NoError, "ErrorCode is not correct")
}
func TestCreateOrganizationForceInvalidParam(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	user := test.MakeUser()

	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/?force=invalid", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity, "http StatusCode is not OK")
}

func TestCreateOrganizationInvalidAgolaRef(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()

	organizationReqDto.AgolaRef = "invalidOrg."
	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	expected := dto.CreateOrganizationResponseDto{ErrorCode: dto.AgolaRefNotValid}
	assert.Equal(t, responseDto.ErrorCode, expected.ErrorCode, "AgolaRef is not valid")
}
func TestCreateOrganizationParametersVisibilityInvalid(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()

	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()
	organizationReqDto.Visibility = "invalid"
	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity, "http StatusCode is not OK")
}

func TestCreateOrganizationParametersBehaviourTypeInvalid(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()

	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()
	organizationReqDto.BehaviourType = "invalid"
	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity, "http StatusCode is not OK")
}

func TestCreateOrganizationGitSourceInvalid(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()
	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()

	db.EXPECT().GetGitSourceByName(user.GitSourceName).Return(nil, nil)

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity, "http StatusCode is not OK")
}

func TestCreateOrganizationGitSourceAlreadyExists(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()
	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()

	db.EXPECT().GetGitSourceByName(user.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(false)

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)
	assert.Equal(t, err, nil)

	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	expected := dto.CreateOrganizationResponseDto{ErrorCode: dto.GitOrganizationNotFoundError}
	assert.Equal(t, responseDto.ErrorCode, expected.ErrorCode, "Git Organization not found")
}

func TestCreateOrganizationWhenCreateWebhookFailed(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()
	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()

	db.EXPECT().GetOrganizationsByGitSource(user.GitSourceName).Return(&organizationList, nil)
	db.EXPECT().GetGitSourceByName(user.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(nil, nil)
	giteaApi.EXPECT().CreateWebHook(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, errors.New(string("someError")))

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not OK")
}

func TestCreateOrganizationWhenFailsToCreateInAgola(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()
	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()

	db.EXPECT().GetOrganizationsByGitSource(user.GitSourceName).Return(&organizationList, nil)
	db.EXPECT().GetGitSourceByName(user.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(nil, nil)
	giteaApi.EXPECT().CreateWebHook(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
	agolaApiInt.EXPECT().CheckOrganizationExists(gomock.Any()).Return(false, "")
	agolaApiInt.EXPECT().CreateOrganization(gomock.Any(), organizationReqDto.Visibility).Return("123456", errors.New(string("someError")))
	giteaApi.EXPECT().DeleteWebHook(gomock.Any(), gomock.Any(), organizationReqDto.Name, 1).Return(nil)

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not OK")
}

func TestCreateOrganizationWhenFailedToSaveOrgInDB(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()
	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()

	organization := (*test.MakeOrganizationList())[0]
	insertRunsData(&organization)
	organizationList := make([]model.Organization, 0)
	organizationList = append(organizationList, organization)

	db.EXPECT().GetGitSourceByName(user.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(nil, nil)
	db.EXPECT().GetOrganizationsByGitSource(user.GitSourceName).Return(&organizationList, nil)
	giteaApi.EXPECT().CreateWebHook(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
	agolaApiInt.EXPECT().CheckOrganizationExists(gomock.Any()).Return(false, "")
	db.EXPECT().SaveOrganization(gomock.Any()).Return(errors.New(string("someError")))
	agolaApiInt.EXPECT().CreateOrganization(gomock.Any(), organizationReqDto.Visibility).Return("123456", nil)

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "failed to save organization in db")
}

func TestCreateOrganizationWhenOrgNameAlreadyExistsInPapagaio(t *testing.T) {
	setupMock(t)

	user := test.MakeUser()
	ts := httptest.NewServer(setupRouter(user))
	client := ts.Client()

	organization := (*test.MakeOrganizationList())[0]
	organizationReqDto.Name = organization.Name
	insertRunsData(&organization)
	organizationList := make([]model.Organization, 0)
	organizationList = append(organizationList, organization)

	db.EXPECT().GetGitSourceByName(user.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), gomock.Any(), organizationReqDto.Name).Return(true)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(nil, nil)
	db.EXPECT().GetOrganizationsByGitSource(user.GitSourceName).Return(&organizationList, nil)
	giteaApi.EXPECT().CreateWebHook(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
	agolaApiInt.EXPECT().CheckOrganizationExists(gomock.Any()).Return(false, "")
	db.EXPECT().SaveOrganization(gomock.Any()).Return(errors.New(string("someError")))
	agolaApiInt.EXPECT().CreateOrganization(gomock.Any(), organizationReqDto.Visibility).Return("123456", nil)

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)
	var responseDto dto.CreateOrganizationResponseDto
	test.ParseBody(resp, &responseDto)

	expected := dto.CreateOrganizationResponseDto{ErrorCode: dto.PapagaioOrganizationExistsError}
	assert.Equal(t, responseDto.ErrorCode, expected.ErrorCode, "Organization exists in Papagaio")
}

func setupSynkMembersUserTestMocks(agolaApiInt *mock_agola.MockAgolaApiInterface, giteaApi *mock_gitea.MockGiteaInterface, organizationName string, remoteSource string) {
	gitTeams := []gitDto.TeamResponseDto{
		{
			ID:         1,
			Name:       "Owners",
			Permission: "owner",
		},
	}
	giteaApi.EXPECT().GetOrganizationTeams(gomock.Any(), gomock.Any(), organizationName).Return(&gitTeams, nil)

	gitTeamMembers := []gitDto.UserTeamResponseDto{
		{
			Username: "user.test",
			Email:    "user.test@email.com",
		},
	}
	giteaApi.EXPECT().GetTeamMembers(gomock.Any(), gomock.Any(), 1).Return(&gitTeamMembers, nil)

	remoteSourceDto := agola.RemoteSourceDto{ID: "123456"}
	agolaApiInt.EXPECT().GetRemoteSource("gitea").Return(&remoteSourceDto, nil)

	users := []agola.UserDto{
		{
			Username: "usertest",
			LinkedAccounts: []agola.LinkedAccountDto{
				{RemoteUserName: "user.test", RemoteSourceID: remoteSourceDto.ID},
			},
		},
	}

	agolaApiInt.EXPECT().GetUsers().Return(&users, nil)

	agolaApiInt.EXPECT().GetOrganizationMembers(gomock.Any()).Return(&agola.OrganizationMembersResponseDto{}, nil)
	agolaApiInt.EXPECT().AddOrUpdateOrganizationMember(gomock.Any(), "usertest", "owner")
}

func setupCheckoutAllGitRepositoryEmptyMocks(giteaApi *mock_gitea.MockGiteaInterface, organizationName string) {
	repositoryList := make([]string, 0)
	giteaApi.EXPECT().GetRepositories(gomock.Any(), gomock.Any(), organizationName).Return(&repositoryList, nil)
}
