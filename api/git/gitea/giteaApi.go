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

func CreateWebHook(gitSource *model.GitSource, gitOrgRef string) (int, error) {
	fmt.Println("CreateWebHook gitOrgRef branchFilter:", gitOrgRef)

	client := &http.Client{}
	URLApi := getCreateWebHookUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)
	webHookConfigPath := controller.GetWebHookPath() + "/" + gitOrgRef

	webHookRequest := CreateWebHookRequestDto{
		Active:       true,
		BranchFilter: "*",
		Config:       WebHookConfigRequestDto{ContentType: "json", URL: config.Config.Server.LocalHostAddress + webHookConfigPath, HTTPMethod: "post"},
		Events:       []string{"repository", "push"},
		Type:         "gitea",
	}
	data, _ := json.Marshal(webHookRequest)
	fmt.Println("json data: ", string(data))

	reqBody := strings.NewReader(string(data))
	req, err := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	fmt.Println("CreateWebHook status response: ", resp.StatusCode, resp.Status)

	if err != nil {
		return -1, err
	}
	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return -1, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var webHookResponse CreateWebHookResponseDto
	json.Unmarshal(body, &webHookResponse)
	fmt.Println("webHookResponse: ", webHookResponse)

	return webHookResponse.ID, err
}

func DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error {
	client := &http.Client{}

	URLApi := getWehHookUrl(gitSource.GitAPIURL, gitOrgRef, fmt.Sprint(webHookID), gitSource.GitToken)

	req, _ := http.NewRequest("DELETE", URLApi, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	return err
}

func GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error) {
	client := &http.Client{}

	URLApi := getGetListRepositoryUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)

	req, err := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
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

func CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool {
	client := &http.Client{}

	URLApi := getOrganizationUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)
	fmt.Println("CheckOrganizationExists URLApi: ", URLApi)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	return api.IsResponseOK(resp.StatusCode)
}

func GetRepositoryTeams(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	client := &http.Client{}

	URLApi := getRepositoryTeamsListUrl(gitSource.GitAPIURL, gitOrgRef, repositoryRef, gitSource.GitToken)

	req, err := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
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

func GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	client := &http.Client{}

	URLApi := getOrganizationTeamsListUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)

	req, err := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
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

func GetTeamMembers(gitSource *model.GitSource, teamId int) (*[]dto.UserTeamResponseDto, error) {
	client := &http.Client{}

	URLApi := getTeamUsersListUrl(gitSource.GitAPIURL, fmt.Sprint(teamId), gitSource.GitToken)

	req, err := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
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

func getBranches(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (*[]BranchResponseDto, error) {
	client := &http.Client{}
	URLApi := getListBranchPath(gitSource.GitAPIURL, gitOrgRef, repositoryRef, gitSource.GitToken)

	req, err := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
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

func getRepositoryAgolaMetadata(gitSource *model.GitSource, gitOrgRef string, repositoryRef string, branchName string) (*[]MetadataResponseDto, error) {
	client := &http.Client{}
	URLApi := getListMetadataPath(gitSource.GitAPIURL, gitOrgRef, repositoryRef, ".agola", branchName, gitSource.GitToken)

	req, err := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
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

func CheckRepositoryAgolaConfExists(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (bool, error) {
	branchList, err := getBranches(gitSource, gitOrgRef, repositoryRef)
	if err != nil {
		return false, err
	}

	for _, branch := range *branchList {
		metadata, err := getRepositoryAgolaMetadata(gitSource, gitOrgRef, repositoryRef, branch.Name)
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

func GetCommitMetadata(gitSource *model.GitSource, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	client := &http.Client{}

	URLApi := getCommitMetadataPath(gitSource.GitAPIURL, gitOrgRef, repositoryRef, commitSha, gitSource.GitToken)

	req, err := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
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
