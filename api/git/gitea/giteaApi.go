package gitea

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"code.gitea.io/sdk/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type GiteaInterface interface {
	CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int64, error)
	DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int64) error
	GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error)
	GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*dto.OrganizationDto, error)
	GetEmailsRepositoryUsersOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]string, error)
	GetRepositoryTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error)
	GetOrganizationTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]dto.TeamResponseDto, error)
	GetTeamMembers(gitSource *model.GitSource, user *model.User, teamId int64) (*[]dto.UserTeamResponseDto, error)
	GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (map[string]bool, error)
	CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error)
	GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error)
	GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error)
	IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error)

	GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error)
	GetUserByLogin(gitSource *model.GitSource, login string) (*dto.UserInfoDto, error)

	GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error)
	RefreshToken(gitSource *model.GitSource, refreshToken string) (*common.Token, error)
}

type GiteaApi struct {
	Db repository.Database
}

const oauth2AuthorizePath string = "%s/login/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=&state=%s"
const oauth2AccessTokenPath string = "%s/login/oauth/access_token"

func (giteaApi *GiteaApi) CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int64, error) {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return -1, err
	}

	optConf := map[string]string{
		"content_type": "json",
		"url":          config.Config.Server.LocalHostAddress + controller.GetWebHookPath() + "/" + organizationRef,
		"http_method":  "post",
	}

	opt := gitea.CreateHookOption{
		Type:         "gitea",
		Config:       optConf,
		Events:       []string{"repository", "push", "create", "delete"},
		Active:       true,
		BranchFilter: "*",
	}

	hook, _, err := client.CreateOrgHook(gitOrgRef, opt)
	if err != nil {
		return -1, err
	}

	return hook.ID, nil
}

func (giteaApi *GiteaApi) DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int64) error {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return err
	}

	_, err = client.DeleteOrgHook(gitOrgRef, webHookID)
	return err
}

func (giteaApi *GiteaApi) GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error) {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	repoList, _, err := client.ListOrgRepos(gitOrgRef, gitea.ListOrgReposOptions{})
	if err != nil {
		return nil, err
	}

	retVal := make([]string, 0)
	for _, repo := range repoList {
		retVal = append(retVal, repo.Name)
	}

	return &retVal, nil
}

func (giteaApi *GiteaApi) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*dto.OrganizationDto, error) {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	org, _, err := client.GetOrg(gitOrgRef)
	if err != nil {
		return nil, err
	}

	if org != nil {
		retVal := dto.OrganizationDto{
			Path:      org.UserName,
			Name:      org.FullName,
			AvatarURL: org.AvatarURL,
			ID:        org.ID,
		}
		return &retVal, nil
	}

	return nil, nil
}

func (giteaApi *GiteaApi) GetRepositoryTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	teams, _, err := client.ListOrgTeams(gitOrgRef+"/"+repositoryRef, gitea.ListTeamsOptions{})
	if err != nil {
		return nil, err
	}

	teamsResponse := make([]dto.TeamResponseDto, 0)
	for _, team := range teams {
		teamsResponse = append(teamsResponse, dto.TeamResponseDto{ID: team.ID, Name: team.Name, Permission: string(team.Permission)})
	}

	return &teamsResponse, nil
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
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	teams, _, err := client.ListOrgTeams(gitOrgRef, gitea.ListTeamsOptions{})
	if err != nil {
		return nil, err
	}

	teamsResponse := make([]dto.TeamResponseDto, 0)
	for _, team := range teams {
		teamsResponse = append(teamsResponse, dto.TeamResponseDto{ID: team.ID, Name: team.Name, Permission: string(team.Permission)})
	}

	return &teamsResponse, nil
}

func (giteaApi *GiteaApi) GetTeamMembers(gitSource *model.GitSource, user *model.User, teamId int64) (*[]dto.UserTeamResponseDto, error) {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	members, _, err := client.ListTeamMembers(teamId, gitea.ListTeamMembersOptions{})
	if err != nil {
		return nil, err
	}

	retVal := make([]dto.UserTeamResponseDto, 0)
	for _, member := range members {
		memberDto := dto.UserTeamResponseDto{
			ID:       member.ID,
			Username: member.UserName,
			Email:    member.Email,
		}
		retVal = append(retVal, memberDto)
	}

	return &retVal, nil
}

func (giteaApi *GiteaApi) GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (map[string]bool, error) {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	branchList, _, err := client.ListRepoBranches(gitOrgRef, repositoryRef, gitea.ListRepoBranchesOptions{})
	if err != nil {
		return nil, err
	}

	retVal := make(map[string]bool)

	for _, branche := range branchList {
		retVal[branche.Name] = true
	}

	return retVal, nil
}

func (giteaApi *GiteaApi) getRepositoryAgolaMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, branchName string) (*[]MetadataResponseDto, error) {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	contents, _, err := client.ListContents(gitOrgRef, repositoryRef, branchName, ".agola")
	if err != nil {
		return nil, err
	}

	retVal := make([]MetadataResponseDto, 0)
	for _, content := range contents {
		metadataDto := MetadataResponseDto{
			Name: content.Name,
			Type: content.Type,
			Size: content.Size,
		}
		retVal = append(retVal, metadataDto)
	}

	return &retVal, nil
}

func (giteaApi *GiteaApi) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error) {
	branchList, err := giteaApi.GetBranches(gitSource, user, gitOrgRef, repositoryRef)
	if err != nil {
		return false, err
	}

	for branch := range branchList {
		metadata, err := giteaApi.getRepositoryAgolaMetadata(gitSource, user, gitOrgRef, repositoryRef, branch)
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
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	commitMetadata, _, err := client.GetSingleCommit(gitOrgRef, repositoryRef, commitSha)
	if err != nil {
		return nil, err
	}

	if commitMetadata != nil && commitMetadata.CommitMeta != nil {
		author := make(map[string]string)
		if commitMetadata.Author != nil {
			author["email"] = commitMetadata.Author.Email
		}

		retVal := dto.CommitMetadataDto{
			Sha:    commitMetadata.CommitMeta.SHA,
			Author: author,
		}
		if commitMetadata.Parents != nil {
			retVal.Parents = make([]dto.CommitParentDto, 0)
			for _, parent := range commitMetadata.Parents {
				retVal.Parents = append(retVal.Parents, dto.CommitParentDto{Sha: parent.SHA})
			}
		}

		return &retVal, nil
	}

	return nil, nil
}

func (giteaApi *GiteaApi) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error) {
	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	organizations, _, err := client.ListMyOrgs(gitea.ListOrgsOptions{})
	if err != nil {
		return nil, err
	}

	retVal := make([]dto.OrganizationDto, 0)
	for _, org := range organizations {
		orgDto := dto.OrganizationDto{
			Path:      org.UserName,
			Name:      org.FullName,
			AvatarURL: org.AvatarURL,
			ID:        org.ID,
		}
		retVal = append(retVal, orgDto)
	}

	return &retVal, nil
}

func (giteaApi *GiteaApi) IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error) {
	teams, err := giteaApi.GetOrganizationTeams(gitSource, user, gitOrgRef)

	if err != nil || teams == nil {
		log.Println("IsUserOwner error in GetOrganizationTeams:", err)
		return false, err
	}

	for _, team := range *teams {
		if team.HasOwnerPermission() {
			members, err := giteaApi.GetTeamMembers(gitSource, user, team.ID)
			if err != nil || members == nil {
				log.Println("IsUserOwner error in GetTeamMembers:", err)
				return false, err
			}

			for _, member := range *members {
				if uint64(member.ID) == user.ID {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (giteaApi *GiteaApi) GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	log.Println("GetUserInfo start")

	client, err := giteaApi.getClient(gitSource, user)
	if err != nil {
		return nil, err
	}

	userInfo, _, err := client.GetMyUserInfo()
	if err != nil || userInfo == nil {
		return nil, err
	}

	retVal := dto.UserInfoDto{
		ID:          userInfo.ID,
		Login:       userInfo.UserName,
		Email:       userInfo.Email,
		FullName:    userInfo.FullName,
		AvatarURL:   userInfo.AvatarURL,
		IsAdmin:     userInfo.IsAdmin,
		UserPageURL: gitSource.GitAPIURL + "/" + userInfo.UserName,
	}

	return &retVal, nil
}

func (giteaApi *GiteaApi) GetUserByLogin(gitSource *model.GitSource, login string) (*dto.UserInfoDto, error) {
	log.Println("GetUserByLogin start")

	client, err := giteaApi.getClient(gitSource, nil)
	if err != nil {
		return nil, err
	}

	userInfo, _, err := client.GetUserInfo(login)
	if err != nil || userInfo == nil {
		return nil, err
	}

	retVal := dto.UserInfoDto{
		ID:          userInfo.ID,
		Login:       userInfo.UserName,
		Email:       userInfo.Email,
		FullName:    userInfo.FullName,
		AvatarURL:   userInfo.AvatarURL,
		IsAdmin:     userInfo.IsAdmin,
		UserPageURL: gitSource.GitAPIURL + "/" + userInfo.UserName,
	}

	return &retVal, nil
}

func (giteaApi *GiteaApi) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error) {
	log.Println("GetOauth2AccessToken start")

	client := &http.Client{}

	URLApi := fmt.Sprintf(oauth2AccessTokenPath, gitSource.GitAPIURL)
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

	URLApi := fmt.Sprintf(oauth2AccessTokenPath, gitSource.GitAPIURL)
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

func GetOauth2AuthorizeUrl(gitApiUrl string, gitClientId string, redirectUrl string, state string) string {
	return fmt.Sprintf(oauth2AuthorizePath, gitApiUrl, gitClientId, url.QueryEscape(redirectUrl), state)
}

type Extra struct {
	Expiry int `json:"expires_in,omitempty"`
}

///////////////

func (giteabApi *GiteaApi) getClient(gitSource *model.GitSource, user *model.User) (*gitea.Client, error) {
	if user == nil {
		return gitea.NewClient(gitSource.GitAPIURL)
	}

	if common.IsAccessTokenExpired(user.Oauth2AccessTokenExpiresAt) {
		log.Println("Token expired is to refresh")
		token, err := giteabApi.RefreshToken(gitSource, user.Oauth2RefreshToken)

		if err != nil {
			log.Println("error during refresh token")
			return nil, err
		}

		user.Oauth2AccessToken = token.AccessToken
		user.Oauth2RefreshToken = token.RefreshToken
		user.Oauth2AccessTokenExpiresAt = token.ExpiryAt

		err = giteabApi.Db.SaveUser(user)

		if err != nil {
			log.Println("error in SaveUser:", err)
			return nil, err
		}
	}

	return gitea.NewClient(gitSource.GitAPIURL, gitea.SetToken(user.Oauth2AccessToken))
}
