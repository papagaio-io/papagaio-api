package gitea

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type GiteaInterface interface {
	CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int, error)
	DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int) error
	GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error)
	GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) *dto.OrganizationDto
	CheckOrganizationExists(gitSource *model.GitSource, user *model.User, gitOrgRef string) bool
	GetRepositoryTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error)
	GetOrganizationTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]dto.TeamResponseDto, error)
	GetTeamMembers(gitSource *model.GitSource, user *model.User, teamId int) (*[]dto.UserTeamResponseDto, error)
	GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) map[string]bool
	CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error)
	GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error)
	GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]string, error)
	IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error)

	GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error)

	GetOauth2AccessToken(gitSource *model.GitSource, code string) (*oauth2.Token, error)
	RefreshToken(gitSource *model.GitSource, refreshToken string) (*oauth2.Token, error)
}

type GiteaApi struct {
	Db repository.Database
}

func (giteaApi *GiteaApi) CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int, error) {
	client, _ := giteaApi.getClient(gitSource, user)
	URLApi := getCreateWebHookUrl(gitSource.GitAPIURL, gitOrgRef)
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

func (giteaApi *GiteaApi) DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int) error {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getWehHookUrl(gitSource.GitAPIURL, gitOrgRef, fmt.Sprint(webHookID))

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

func (giteaApi *GiteaApi) GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error) {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getGetListRepositoryUrl(gitSource.GitAPIURL, gitOrgRef)

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

func (giteaApi *GiteaApi) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) *dto.OrganizationDto {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getOrganizationUrl(gitSource.GitAPIURL, gitOrgRef)
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

func (giteaApi *GiteaApi) CheckOrganizationExists(gitSource *model.GitSource, user *model.User, gitOrgRef string) bool {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getOrganizationUrl(gitSource.GitAPIURL, gitOrgRef)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return api.IsResponseOK(resp.StatusCode)
}

func (giteaApi *GiteaApi) GetRepositoryTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getRepositoryTeamsListUrl(gitSource.GitAPIURL, gitOrgRef, repositoryRef)

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

func (giteaApi *GiteaApi) GetOrganizationTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getOrganizationTeamsListUrl(gitSource.GitAPIURL, gitOrgRef)

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

func (giteaApi *GiteaApi) GetTeamMembers(gitSource *model.GitSource, user *model.User, teamId int) (*[]dto.UserTeamResponseDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getTeamUsersListUrl(gitSource.GitAPIURL, fmt.Sprint(teamId))

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

func (giteaApi *GiteaApi) GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) map[string]bool {
	branchList, _ := giteaApi.getBranches(gitSource, user, gitOrgRef, repositoryRef)
	retVal := make(map[string]bool)

	for _, branche := range *branchList {
		retVal[branche.Name] = true
	}

	return retVal
}

func (giteaApi *GiteaApi) getBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]BranchResponseDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)
	URLApi := getListBranchUrl(gitSource.GitAPIURL, gitOrgRef, repositoryRef)

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

func (giteaApi *GiteaApi) getRepositoryAgolaMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, branchName string) (*[]MetadataResponseDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)
	URLApi := getListMetadataUrl(gitSource.GitAPIURL, gitOrgRef, repositoryRef, ".agola", branchName)

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

func (giteaApi *GiteaApi) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error) {
	branchList, err := giteaApi.getBranches(gitSource, user, gitOrgRef, repositoryRef)
	if err != nil {
		return false, err
	}

	for _, branch := range *branchList {
		metadata, err := giteaApi.getRepositoryAgolaMetadata(gitSource, user, gitOrgRef, repositoryRef, branch.Name)
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

func (giteaApi *GiteaApi) GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getCommitMetadataPath(gitSource.GitAPIURL, gitOrgRef, repositoryRef, commitSha)

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

func (giteaApi *GiteaApi) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]string, error) {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getOrganizationsUrl(gitSource.GitAPIURL)
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
			isUserOwner, _ := giteaApi.IsUserOwner(gitSource, user, org.Username)
			if isUserOwner {
				retVal = append(retVal, org.Username)
			}
		}

		return &retVal, nil
	}

	return nil, err
}

func (giteaApi *GiteaApi) IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error) {
	teams, _ := giteaApi.GetOrganizationTeams(gitSource, user, gitOrgRef)
	for _, team := range *teams {
		if team.HasOwnerPermission() {
			members, _ := giteaApi.GetTeamMembers(gitSource, user, team.ID)
			for _, member := range *members {
				if member.ID == int(user.ID) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (giteaApi *GiteaApi) GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getUserInfoUrl(gitSource.GitAPIURL)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response dto.UserInfoDto
		json.Unmarshal(body, &response)

		return &response, nil
	}

	return nil, err
}

func (giteaApi *GiteaApi) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*oauth2.Token, error) {
	client := &http.Client{}

	URLApi := getOauth2AccessTokenUrl(gitSource.GitAPIURL)
	accessTokenRequest := dto.AccessTokenRequestDto{ClientID: gitSource.GitClientID, ClientSecret: gitSource.GitSecret, GrantType: "authorization_code", Code: code, RedirectURL: controller.GetRedirectUrl()}
	data, _ := json.Marshal(accessTokenRequest)
	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", URLApi, reqBody)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response oauth2.Token
		json.Unmarshal(body, &response)

		return &response, nil
	}

	return nil, err
}

func (giteaApi *GiteaApi) RefreshToken(gitSource *model.GitSource, refreshToken string) (*oauth2.Token, error) {
	client := &http.Client{}

	URLApi := getOauth2AccessTokenUrl(gitSource.GitAPIURL)
	accessTokenRequest := dto.AccessTokenRequestDto{ClientID: gitSource.GitClientID, ClientSecret: gitSource.GitSecret, GrantType: "refresh_token", RedirectURL: controller.GetRedirectUrl()}
	data, _ := json.Marshal(accessTokenRequest)
	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", URLApi, reqBody)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response oauth2.Token
		json.Unmarshal(body, &response)

		return &response, nil
	}

	return nil, err
}

///////////////

func (giteaApi *GiteaApi) getClient(gitSource *model.GitSource, user *model.User) (*httpClient, error) {
	if common.IsAccessTokenExpired(user.Oauth2AccessTokenExpiresAt) {
		if user.UserID == nil {
			log.Println("Can not refresh token for user nil")
			return nil, errors.New("Can not refresh token for user nil")
		}

		token, err := giteaApi.RefreshToken(gitSource, user.Oauth2RefreshToken)

		if err != nil {
			log.Println("error during refresh token")
			return nil, err
		}

		user.Oauth2AccessToken = token.AccessToken
		user.Oauth2RefreshToken = token.RefreshToken
		user.Oauth2AccessTokenExpiresAt = token.Expiry

		giteaApi.Db.SaveUser(user)
	}

	client := &httpClient{c: &http.Client{}, accessToken: user.Oauth2AccessToken}

	return client, nil
}

type httpClient struct {
	c           *http.Client
	accessToken string
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "bearer "+c.accessToken)
	return c.c.Do(req)

}
