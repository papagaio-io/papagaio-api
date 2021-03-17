package gitea

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/api"
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
		Events:       []string{"repository"},
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

	URLApi := getDeleteWehHookUrl(gitSource.GitAPIURL, gitOrgRef, string(webHookID), gitSource.GitToken)

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

func GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]TeamResponseDto, error) {
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

	var teamsResponse []TeamResponseDto
	json.Unmarshal(body, &teamsResponse)

	return &teamsResponse, err
}

func GetTeamMembers(gitSource *model.GitSource, teamId int) (*[]UserTeamResponseDto, error) {
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

	var usersResponse []UserTeamResponseDto
	json.Unmarshal(body, &usersResponse)

	return &usersResponse, err
}
