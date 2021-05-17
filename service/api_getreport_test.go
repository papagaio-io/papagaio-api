package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	gitDto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/test"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

func TestGetReportOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)

	organization := (*test.MakeOrganizationList())[0]
	insertRunsData(&organization)
	organizationList := make([]model.Organization, 0)
	organizationList = append(organizationList, organization)

	gitSource := (*test.MakeGitSourceMap())[organization.GitSourceName]

	db.EXPECT().GetOrganizations().Return(&organizationList, nil)
	db.EXPECT().GetGitSourceByName(organization.GitSourceName).Return(&gitSource, nil)
	giteaApi.EXPECT().GetOrganization(gomock.Any(), organization.Name).Return(&gitDto.OrganizationDto{})

	serviceOrganization := OrganizationService{
		Db:         db,
		GitGateway: &git.GitGateway{GiteaApi: giteaApi},
	}

	ts := httptest.NewServer(http.HandlerFunc(serviceOrganization.GetReport))
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

	var organizationsDto []dto.OrganizationDto
	test.ParseBody(resp, &organizationsDto)

	assert.Check(t, len(organizationsDto) == 1)
	assertOrganizationDto(t, &organization, &organizationsDto[0])
}

func TestGetProjectReportOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db := mock_repository.NewMockDatabase(ctl)
	giteaApi := mock_gitea.NewMockGiteaInterface(ctl)

	organization := (*test.MakeOrganizationList())[0]
	insertRunsData(&organization)

	db.EXPECT().GetOrganizationByAgolaRef(organization.AgolaOrganizationRef).Return(&organization, nil)

	serviceOrganization := OrganizationService{
		Db:         db,
		GitGateway: &git.GitGateway{GiteaApi: giteaApi},
	}

	router := mux.NewRouter()
	router.HandleFunc("/{organizationRef}/{projectName}", serviceOrganization.GetProjectReport)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/" + organization.AgolaOrganizationRef + "/test1")

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

	var projectDto dto.ProjectDto
	test.ParseBody(resp, &projectDto)

	assertProjectDto(t, &projectDto)
}

func assertOrganizationDto(t *testing.T, organization *model.Organization, organizationDto *dto.OrganizationDto) {
	organizationDto.Projects = test.SortProjectsDto(organizationDto.Projects)

	assert.Equal(t, organization.AgolaOrganizationRef, organizationDto.AgolaRef, "AgolaRef is not correct")
	assert.Equal(t, organization.Name, organizationDto.Name, "Name is not correct")
	assert.Equal(t, len(organization.Projects), len(organizationDto.Projects), "There are not all the projects")

	projectA := organizationDto.Projects[0]
	projectA.Branchs = test.SortBranchesDto(projectA.Branchs)
	assert.Equal(t, projectA.Name, "test1")
	assert.Equal(t, projectA.Branchs[0].Name, "master")
	assert.Equal(t, projectA.Branchs[0].State, types.RunStateSuccess)
	assert.Equal(t, projectA.Branchs[1].Name, "test")
	assert.Equal(t, projectA.Branchs[1].State, types.RunStateFailed)
	assert.Equal(t, projectA.WorstReport.BranchName, "test")

	projectB := organizationDto.Projects[1]
	assert.Equal(t, projectB.Name, "test2")
	assert.Equal(t, projectB.Branchs[0].Name, "master")

	projectC := organizationDto.Projects[2]
	assert.Equal(t, projectC.Name, "test3")
	assert.Equal(t, projectC.Branchs[0].Name, "master")
	assert.Equal(t, projectC.Branchs[0].State, types.RunStateFailed)
	assert.Equal(t, projectC.WorstReport.BranchName, "master")

	projectD := organizationDto.Projects[3]
	assert.Equal(t, projectD.Name, "test4")

	projectE := organizationDto.Projects[4]
	assert.Equal(t, projectE.Name, "test5")
	assert.Equal(t, projectE.Branchs[0].Name, "master")
	assert.Equal(t, projectE.Branchs[0].State, types.RunStateNone)
}

func assertProjectDto(t *testing.T, projectDto *dto.ProjectDto) {
	projectDto.Branchs = test.SortBranchesDto(projectDto.Branchs)
	assert.Equal(t, projectDto.Name, "test1")
	assert.Equal(t, projectDto.Branchs[0].Name, "master")
	assert.Equal(t, projectDto.Branchs[0].State, types.RunStateSuccess)
	assert.Equal(t, projectDto.Branchs[1].Name, "test")
	assert.Equal(t, projectDto.Branchs[1].State, types.RunStateFailed)
	assert.Equal(t, projectDto.WorstReport.BranchName, "test")

}

func insertRunsData(organization *model.Organization) {
	projects := make(map[string]model.Project)

	now := time.Now()

	projectA := model.Project{
		GitRepoPath:     "test1",
		AgolaProjectRef: "test1",
		AgolaProjectID:  "test1_123456",
		Archivied:       false,
	}
	projectA.PushNewRun(model.RunInfo{
		ID:           "1",
		Branch:       "master",
		Phase:        types.RunPhaseFinished,
		Result:       types.RunResultSuccess,
		RunStartDate: now.AddDate(0, 0, -1),
		RunEndDate:   now,
	})
	projectA.PushNewRun(model.RunInfo{
		ID:           "2",
		Branch:       "test",
		Phase:        types.RunPhaseFinished,
		Result:       types.RunResultFailed,
		RunStartDate: now.AddDate(0, 0, -2),
		RunEndDate:   now,
	})
	projects[projectA.GitRepoPath] = projectA
	//
	projectB := model.Project{
		GitRepoPath:     "test2",
		AgolaProjectRef: "test2",
		AgolaProjectID:  "test2_123456",
		Archivied:       false,
	}
	projectB.PushNewRun(model.RunInfo{
		ID:           "1",
		Branch:       "master",
		Phase:        types.RunPhaseFinished,
		Result:       types.RunResultSuccess,
		RunStartDate: now.AddDate(0, 0, -1),
		RunEndDate:   now,
	})
	projects[projectB.GitRepoPath] = projectB
	//
	projectC := model.Project{
		GitRepoPath:     "test3",
		AgolaProjectRef: "test3",
		AgolaProjectID:  "test3_123456",
		Archivied:       false,
	}
	projectC.PushNewRun(model.RunInfo{
		ID:           "1",
		Branch:       "master",
		Phase:        types.RunPhaseFinished,
		Result:       types.RunResultFailed,
		RunStartDate: now.AddDate(0, 0, -1),
		RunEndDate:   now,
	})
	projects[projectC.GitRepoPath] = projectC
	//Empty project
	projectD := model.Project{
		GitRepoPath:     "test4",
		AgolaProjectRef: "test4",
		AgolaProjectID:  "test4_123456",
		Archivied:       true,
	}
	projects[projectD.GitRepoPath] = projectD
	//Empty branch
	projectE := model.Project{
		GitRepoPath:     "test5",
		AgolaProjectRef: "test5",
		AgolaProjectID:  "test5_123456",
		Archivied:       true,
	}
	projectE.Branchs = make(map[string]model.Branch)
	projectE.Branchs["master"] = model.Branch{Name: "master"}
	projects[projectE.GitRepoPath] = projectE

	organization.Projects = projects
}
