package gitea

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

type GiteaInterface interface {
	CreateWebHook(gitSource *model.GitSource, gitOrgRef string, organizationRef string) (int, error)
	DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error
	GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error)
	GetOrganization(gitSource *model.GitSource, gitOrgRef string) *dto.OrganizationDto
	CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool
	GetRepositoryTeams(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error)
	GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]dto.TeamResponseDto, error)
	GetTeamMembers(gitSource *model.GitSource, teamId int) (*[]dto.UserTeamResponseDto, error)
	GetBranches(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) map[string]bool
	CheckRepositoryAgolaConfExists(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (bool, error)
	GetCommitMetadata(gitSource *model.GitSource, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error)
	GetOrganizations(gitSource *model.GitSource) (*[]string, error)
}

type GiteaApi struct{}

func (giteaApi *GiteaApi) CreateWebHook(gitSource *model.GitSource, gitOrgRef string, organizationRef string) (int, error) {
	client := &http.Client{}
	URLApi := getCreateWebHookUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)
	webHookConfigPath := controller.GetWebHookPath() + "/" + organizationRef

	webHookRequest := CreateWebHookRequestDto{
		Active:       true,
		BranchFilter: "*",
		Config:       WebHookConfigRequestDto{ContentType: "json", URL: config.Config.Server.LocalHostAddress + webHookConfigPath, HTTPMethod: "post"},
		Events:       []string{"repository", "push", "create", "delete"},
		Type:         "gitea",
	}
	data, _ := json.Marshal(webHookRequest)

	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return -1, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var webHookResponse CreateWebHookResponseDto
	json.Unmarshal(body, &webHookResponse)

	return webHookResponse.ID, err
}

func (giteaApi *GiteaApi) DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error {
	client := &http.Client{}

	URLApi := getWehHookUrl(gitSource.GitAPIURL, gitOrgRef, fmt.Sprint(webHookID), gitSource.GitToken)

	req, _ := http.NewRequest("DELETE", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	return err
}

func (giteaApi *GiteaApi) GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error) {
	client := &http.Client{}

	URLApi := getGetListRepositoryUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var repositoryesResponse []RepositoryDto
	json.Unmarshal(body, &repositoryesResponse)

	retVal := make([]string, 0)
	for _, repo := range repositoryesResponse {
		retVal = append(retVal, repo.Name)
	}

	return &retVal, err
}

func (giteaApi *GiteaApi) GetOrganization(gitSource *model.GitSource, gitOrgRef string) *dto.OrganizationDto {
	client := &http.Client{}

	URLApi := getOrganizationUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var data OrganizationResponseDto
		json.Unmarshal(body, &data)

		organization := dto.OrganizationDto{Name: data.Name, ID: data.ID, AvatarURL: data.AvatarURL}
		return &organization
	}

	return nil
}

func (giteaApi *GiteaApi) CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool {
	client := &http.Client{}

	URLApi := getOrganizationUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return api.IsResponseOK(resp.StatusCode)
}

func (giteaApi *GiteaApi) GetRepositoryTeams(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	client := &http.Client{}

	URLApi := getRepositoryTeamsListUrl(gitSource.GitAPIURL, gitOrgRef, repositoryRef, gitSource.GitToken)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var teamsResponse []dto.TeamResponseDto
	json.Unmarshal(body, &teamsResponse)

	return &teamsResponse, err
}

func (giteaApi *GiteaApi) GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	client := &http.Client{}

	URLApi := getOrganizationTeamsListUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var teamsResponse []dto.TeamResponseDto
	json.Unmarshal(body, &teamsResponse)

	return &teamsResponse, err
}

func (giteaApi *GiteaApi) GetTeamMembers(gitSource *model.GitSource, teamId int) (*[]dto.UserTeamResponseDto, error) {
	client := &http.Client{}

	URLApi := getTeamUsersListUrl(gitSource.GitAPIURL, fmt.Sprint(teamId), gitSource.GitToken)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var usersResponse []dto.UserTeamResponseDto
	json.Unmarshal(body, &usersResponse)

	return &usersResponse, err
}

func (giteaApi *GiteaApi) GetBranches(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) map[string]bool {
	branchList, _ := giteaApi.getBranches(gitSource, gitOrgRef, repositoryRef)
	retVal := make(map[string]bool)

	for _, branche := range *branchList {
		retVal[branche.Name] = true
	}

	return retVal
}

func (giteaApi *GiteaApi) getBranches(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (*[]BranchResponseDto, error) {
	client := &http.Client{}
	URLApi := getListBranchPath(gitSource.GitAPIURL, gitOrgRef, repositoryRef, gitSource.GitToken)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var branchesResponse []BranchResponseDto
	json.Unmarshal(body, &branchesResponse)

	return &branchesResponse, err
}

func (giteaApi *GiteaApi) getRepositoryAgolaMetadata(gitSource *model.GitSource, gitOrgRef string, repositoryRef string, branchName string) (*[]MetadataResponseDto, error) {
	client := &http.Client{}
	URLApi := getListMetadataPath(gitSource.GitAPIURL, gitOrgRef, repositoryRef, ".agola", branchName, gitSource.GitToken)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var metadataResponse []MetadataResponseDto
	json.Unmarshal(body, &metadataResponse)

	return &metadataResponse, err
}

func (giteaApi *GiteaApi) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (bool, error) {
	branchList, err := giteaApi.getBranches(gitSource, gitOrgRef, repositoryRef)
	if err != nil {
		return false, err
	}

	for _, branch := range *branchList {
		metadata, err := giteaApi.getRepositoryAgolaMetadata(gitSource, gitOrgRef, repositoryRef, branch.Name)
		if err != nil {
			return false, err
		}

		for _, file := range *metadata {
			if strings.Compare(file.Type, "file") == 0 && (strings.Compare(file.Name, "config.jsonnet") == 0 || strings.Compare(file.Name, "config.yml") == 0 || strings.Compare(file.Name, "config.json") == 0) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (giteaApi *GiteaApi) GetCommitMetadata(gitSource *model.GitSource, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	client := &http.Client{}

	URLApi := getCommitMetadataPath(gitSource.GitAPIURL, gitOrgRef, repositoryRef, commitSha, gitSource.GitToken)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var commitMetadataResponse []dto.CommitMetadataDto
	json.Unmarshal(body, &commitMetadataResponse)

	if len(commitMetadataResponse) == 1 {
		return &commitMetadataResponse[0], err
	} else {
		return nil, err
	}
}

func (giteaApi *GiteaApi) GetOrganizations(gitSource *model.GitSource) (*[]string, error) {
	client := &http.Client{}

	URLApi := getOrganizationsPath(gitSource.GitAPIURL, gitSource.GitToken)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var organizations []OrganizationResponseDto
		json.Unmarshal(body, &organizations)

		retVal := make([]string, 0)
		for _, org := range organizations {
			retVal = append(retVal, org.Username)
		}

		return &retVal, nil
	}

	return nil, err
}
