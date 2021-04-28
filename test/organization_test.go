package test

import (
	"encoding/json"
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
	"wecode.sorint.it/opensource/papagaio-api/service"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_agola"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
)

func TestGetOrganizationsOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organizationsMock := MakeOrganizationList()

	db := mock_repository.NewMockDatabase(ctl)
	db.EXPECT().GetOrganizations().Return(organizationsMock, nil)

	serviceOrganization := service.OrganizationService{
		Db: db,
	}

	ts := httptest.NewServer(http.HandlerFunc(serviceOrganization.GetOrganizations))
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL)

	assert.Equal(t, err, nil)

	var organizations *[]model.Organization
	parseBody(resp, organizations)

	//assertEventDtoEquals(t, eventDto, dto.CreateFullEventDto(eventMock))
}

func TestCreateOrganizationOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)

	organizationReqDto := dto.CreateOrganizationRequestDto{
		Name:          "Test",
		AgolaRef:      "Test",
		Visibility:    dto.Public,
		GitSourceName: "gitea",
		BehaviourType: dto.None,
	}

	gitSource := (*MakeGitSourceMap())[organizationReqDto.GitSourceName]
	db.EXPECT().GetGitSourceByName(gomock.Eq(organizationReqDto.GitSourceName)).Return(&gitSource, nil)
	giteaApi.EXPECT().CheckOrganizationExists(gomock.Any(), organizationReqDto.Name).Return(true)
	db.EXPECT().GetOrganizationByAgolaRef(organizationReqDto.AgolaRef).Return(nil, nil)
	giteaApi.EXPECT().CreateWebHook(gomock.Any(), organizationReqDto.Name, organizationReqDto.AgolaRef).Return(1, nil)
	agolaApi.EXPECT().CheckOrganizationExists(gomock.Any()).Return(false, "")
	agolaApi.EXPECT().CreateOrganization(gomock.Any(), organizationReqDto.Visibility).Return("123456", nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)

	setupSynkMembersUserTestMocks(agolaApi, giteaApi, organizationReqDto.Name)
	setupCheckoutAllGitRepositoryEmptyMocks(giteaApi, organizationReqDto.Name)

	serviceOrganization := service.OrganizationService{
		Db:         db,
		AgolaApi:   agolaApi,
		GitGateway: &git.GitGateway{GiteaApi: giteaApi},
	}

	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(controller.XAuthEmail, "testuser")
			h.ServeHTTP(w, r)
		})
	})
	router.HandleFunc("/", serviceOrganization.CreateOrganization)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(organizationReqDto)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/", "application/json", requestBody)

	assert.Equal(t, err, nil)

	var responseDto dto.CreateOrganizationResponseDto
	parseBody(resp, &responseDto)

	assert.Equal(t, responseDto.AgolaExists, false, "AgolaExists is not correct")
	assert.Check(t, strings.Contains(responseDto.OrganizationURL, "/org/"+organizationReqDto.AgolaRef), "OrganizationURL is not correct")
}

func setupSynkMembersUserTestMocks(agolaApi *mock_agola.MockAgolaApiInterface, giteaApi *mock_gitea.MockGiteaInterface, organizationName string) {
	gitTeams := []gitDto.TeamResponseDto{
		gitDto.TeamResponseDto{
			ID:         1,
			Name:       "Owners",
			Permission: "owner",
		},
	}
	giteaApi.EXPECT().GetOrganizationTeams(gomock.Any(), organizationName).Return(&gitTeams, nil)

	gitTeamMembers := []gitDto.UserTeamResponseDto{
		gitDto.UserTeamResponseDto{
			Username: "user.test",
			Email:    "user.test@email.com",
		},
	}
	giteaApi.EXPECT().GetTeamMembers(gomock.Any(), 1).Return(&gitTeamMembers, nil)

	agolaApi.EXPECT().GetOrganizationMembers(gomock.Any()).Return(&agola.OrganizationMembersResponseDto{}, nil)
	agolaApi.EXPECT().AddOrUpdateOrganizationMember(gomock.Any(), "usertest", "owner")
}

func setupCheckoutAllGitRepositoryEmptyMocks(giteaApi *mock_gitea.MockGiteaInterface, organizationName string) {
	repositoryList := make([]string, 0)
	giteaApi.EXPECT().GetRepositories(gomock.Any(), organizationName).AnyTimes().Return(&repositoryList, nil) //TODO remove AnyTimes after goroutin removing
}
