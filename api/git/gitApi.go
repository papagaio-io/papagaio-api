package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"wecode.sorint.it/opensource/papagaio-be/config"
	"wecode.sorint.it/opensource/papagaio-be/controller"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func CreateWebHook(gitSource *model.GitSource, gitOrgRef string, branchFilter string) (int, error) {
	client := &http.Client{}

	URLApi := getCreateWebHookUrl(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)
	fmt.Println("CreateWebHook URLApi: ", URLApi)

	webHookConfigPath := fmt.Sprintf(controller.WebHookPath+controller.WenHookPathParam, gitOrgRef)
	webHookRequest := CreateWebHookRequestDto{
		Active:       true,
		BranchFilter: branchFilter,
		Config:       WebHookConfigRequestDto{ContentType: "json", URL: config.Config.Server.LocalHostAddress + webHookConfigPath, HTTPMethod: "post"},
		Events:       []string{"create", "delete", "repository"},
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
	if resp.StatusCode == 400 {
		return -1, errors.New(resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var webHookResponse CreateWebHookResponseDto
	json.Unmarshal(body, &webHookResponse)
	fmt.Println("webHookResponse: ", webHookResponse)

	return webHookResponse.ID, err
}

//TODO
func DeleteWebHook(gitSource *model.GitSource, webHookID int) error {
	var err error
	return err
}

func GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]RepositoryDto, error) {
	client := &http.Client{}

	URLApi := getGetListRepositoryPath(gitSource.GitAPIURL, gitOrgRef, gitSource.GitToken)

	req, err := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 400 {
		return nil, errors.New(resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var repositoryesResponse []RepositoryDto
	json.Unmarshal(body, &repositoryesResponse)

	return &repositoryesResponse, err
}

//TODO
/*func GetGitOrganizations(gitSource *model.GitSource) ([]string, error) {
	var organizations []string
	var err error
	return organizations, err
}*/
