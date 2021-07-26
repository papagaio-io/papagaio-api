package gitea

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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
	GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*dto.OrganizationDto, error)
	GetEmailsRepositoryUsersOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]string, error)
	GetRepositoryTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error)
	GetOrganizationTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]dto.TeamResponseDto, error)
	GetTeamMembers(gitSource *model.GitSource, user *model.User, teamId int) (*[]dto.UserTeamResponseDto, error)
	GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) map[string]bool
	CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error)
	GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error)
	GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error)
	IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error)

	GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error)
	GetUserByLogin(gitSource *model.GitSource, login string) (*dto.UserInfoDto, error)
	CreateAgolaApp(gitSource *model.GitSource, user *model.User) (*CreateOauth2AppResponseDto, error)

	GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error)
	RefreshToken(gitSource *model.GitSource, refreshToken string) (*common.Token, error)
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
	req.Header.Set("content-type", "application/json")
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
	err = json.Unmarshal(body, &webHookResponse)
	if err != nil {
		return -1, err
	}

	return webHookResponse.ID, nil
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
	err = json.Unmarshal(body, &repositoryesResponse)
	if err != nil {
		return nil, err
	}

	retVal := make([]string, 0)
	for _, repo := range repositoryesResponse {
		retVal = append(retVal, repo.Name)
	}

	return &retVal, nil
}

func (giteaApi *GiteaApi) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*dto.OrganizationDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getOrganizationUrl(gitSource.GitAPIURL, gitOrgRef)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var data OrganizationResponseDto
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}

		organization := dto.OrganizationDto{Name: data.Name, Path: data.Name, ID: data.ID, AvatarURL: data.AvatarURL}
		return &organization, nil
	}

	return nil, errors.New("internal error")
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
	err = json.Unmarshal(body, &teamsResponse)
	if err != nil {
		return nil, err
	}

	return &teamsResponse, err
}

func (giteaApi *GiteaApi) GetEmailsRepositoryUsersOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]string, error) {
	retVal := make([]string, 0)

	teams, err := giteaApi.GetRepositoryTeams(gitSource, user, gitOrgRef, repositoryRef)
	if err != nil {
		return nil, err
	}

	for _, team := range *teams {
		if strings.Compare(team.Permission, "owner") != 0 {
			continue
		}

		users, err := giteaApi.GetTeamMembers(gitSource, user, team.ID)
		if err != nil {
			continue
		}

		for _, user := range *users {
			retVal = append(retVal, user.Email)
		}
	}

	return &retVal, nil
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
	err = json.Unmarshal(body, &teamsResponse)
	if err != nil {
		return nil, err
	}

	return &teamsResponse, nil
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
	err = json.Unmarshal(body, &usersResponse)
	if err != nil {
		return nil, err
	}

	return &usersResponse, nil
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
	err = json.Unmarshal(body, &branchesResponse)
	if err != nil {
		return nil, err
	}

	return &branchesResponse, nil
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
	err = json.Unmarshal(body, &metadataResponse)
	if err != nil {
		return nil, err
	}

	return &metadataResponse, nil
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
	err = json.Unmarshal(body, &commitMetadataResponse)
	if err != nil {
		return nil, err
	}

	if len(commitMetadataResponse) == 1 {
		return &commitMetadataResponse[0], err
	} else {
		return nil, err
	}
}

func (giteaApi *GiteaApi) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error) {
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
		err = json.Unmarshal(body, &organizations)
		if err != nil {
			return nil, err
		}

		retVal := make([]dto.OrganizationDto, 0)
		for _, org := range organizations {
			isUserOwner, _ := giteaApi.IsUserOwner(gitSource, user, org.Username)
			if isUserOwner {
				orgDto := dto.OrganizationDto{Path: org.Username, ID: org.ID, AvatarURL: org.AvatarURL}
				if len(org.Name) > 0 {
					orgDto.Name = org.Name
				} else {
					orgDto.Name = org.Username
				}
				retVal = append(retVal, orgDto)
			}
		}

		return &retVal, nil
	}

	return nil, err
}

func (giteaApi *GiteaApi) IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error) {
	teams, err := giteaApi.GetOrganizationTeams(gitSource, user, gitOrgRef)

	if err != nil || teams == nil {
		log.Println("IsUserOwner error in GetOrganizationTeams:", err)
		return false, err
	}

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
	log.Println("GetUserInfo start")

	client, _ := giteaApi.getClient(gitSource, user)

	URLApi := getLoggedUserInfoUrl(gitSource.GitAPIURL)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response dto.UserInfoDto
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		response.UserPageURL = gitSource.GitAPIURL + "/" + response.Login

		log.Println("GetUserInfo end")

		return &response, nil
	}

	return nil, err
}

func (giteaApi *GiteaApi) GetUserByLogin(gitSource *model.GitSource, login string) (*dto.UserInfoDto, error) {
	log.Println("GetUserByLogin start")

	client := &http.Client{}

	URLApi := getUserInfoUrl(gitSource.GitAPIURL, login)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response dto.UserInfoDto
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		log.Println("GetUserByLogin end")

		return &response, nil
	}

	return nil, err
}

func (giteaApi *GiteaApi) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error) {
	log.Println("GetOauth2AccessToken start")

	client := &http.Client{}

	URLApi := getOauth2AccessTokenUrl(gitSource.GitAPIURL)
	accessTokenRequest := dto.AccessTokenRequestDto{ClientID: gitSource.GitClientID, ClientSecret: gitSource.GitSecret, GrantType: "authorization_code", Code: code, RedirectURL: controller.GetRedirectUrl()}
	data, _ := json.Marshal(accessTokenRequest)
	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", URLApi, reqBody)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response common.Token
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		if response.Expiry > 0 {
			response.ExpiryAt = time.Now().Add(time.Second * time.Duration(response.Expiry))
		}

		log.Println("GetOauth2AccessToken end")

		return &response, nil
	} else {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}
}

func (giteaApi *GiteaApi) RefreshToken(gitSource *model.GitSource, refreshToken string) (*common.Token, error) {
	client := &http.Client{}

	URLApi := getOauth2AccessTokenUrl(gitSource.GitAPIURL)
	accessTokenRequest := dto.AccessTokenRequestDto{ClientID: gitSource.GitClientID, ClientSecret: gitSource.GitSecret, GrantType: "refresh_token", RedirectURL: controller.GetRedirectUrl(), RefreshToken: refreshToken}
	data, _ := json.Marshal(accessTokenRequest)
	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response common.Token
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		response.ExpiryAt = time.Now().Add(time.Second * time.Duration(response.Expiry))

		return &response, nil
	} else {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}
}

type Extra struct {
	Expiry int `json:"expires_in,omitempty"`
}

func (giteaApi *GiteaApi) CreateAgolaApp(gitSource *model.GitSource, user *model.User) (*CreateOauth2AppResponseDto, error) {
	client, _ := giteaApi.getClient(gitSource, user)
	URLApi := getCreateOauth2AppUrl(gitSource.GitAPIURL)

	createOauth2AppRequest := CreateOauth2AppRequestDto{
		Name:         "Agola",
		RedirectUris: []string{config.Config.Agola.AgolaAddr + "/oauth2/callback"},
	}
	data, _ := json.Marshal(createOauth2AppRequest)

	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Set("content-type", "application/json")
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

	var response CreateOauth2AppResponseDto
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

///////////////

func (giteaApi *GiteaApi) getClient(gitSource *model.GitSource, user *model.User) (*httpClient, error) {
	if common.IsAccessTokenExpired(user.Oauth2AccessTokenExpiresAt) {
		if user.UserID == nil {
			log.Println("Can not refresh token for user nil")
			return nil, errors.New("can not refresh token for user nil")
		}

		log.Println("token is to refresh")
		token, err := giteaApi.RefreshToken(gitSource, user.Oauth2RefreshToken)

		if err != nil {
			log.Println("error during refresh token")
			return nil, err
		}

		user.Oauth2AccessToken = token.AccessToken
		user.Oauth2RefreshToken = token.RefreshToken
		user.Oauth2AccessTokenExpiresAt = time.Now().Add(time.Second * time.Duration(token.Expiry))

		err = giteaApi.Db.SaveUser(user)
		if err != nil {
			log.Println("error in SaveUser:", err)
			return nil, err
		}
	}

	client := &httpClient{c: &http.Client{}, accessToken: user.Oauth2AccessToken}

	return client, nil
}

type httpClient struct {
	c           *http.Client
	accessToken string
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+c.accessToken)
	return c.c.Do(req)

}
