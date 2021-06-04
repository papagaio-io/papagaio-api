package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type GithubInterface interface {
	CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int, error)
	DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int) error
	GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error)
	CheckOrganizationExists(gitSource *model.GitSource, user *model.User, gitOrgRef string) bool
	GetRepositoryTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error)
	GetOrganizationTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]dto.TeamResponseDto, error)
	GetTeamMembers(gitSource *model.GitSource, user *model.User, organizationName string, teamId int) (*[]dto.UserTeamResponseDto, error)
	GetOrganizationMembers(gitSource *model.GitSource, user *model.User, organizationName string) (*[]GitHubUser, error)
	GetRepositoryMembers(gitSource *model.GitSource, user *model.User, organizationName string, repositoryRef string) (*[]GitHubUser, error)
	GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) map[string]bool
	CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error)
	GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error)
	GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) *dto.OrganizationDto
	GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]string, error)

	GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error)

	GetOauth2AccessToken(gitSource *model.GitSource, code string) (*oauth2.Token, error)
	RefreshToken(gitSource *model.GitSource, refreshToken string) (*oauth2.Token, error)
}

type GithubApi struct {
	Db repository.Database
}

func (githubApi *GithubApi) CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int, error) {
	client, _ := githubApi.getClient(gitSource, user)

	webHookName := "web"
	active := true
	conf := make(map[string]interface{})
	conf["url"] = config.Config.Server.LocalHostAddress + controller.GetWebHookPath() + "/" + organizationRef
	conf["content_type"] = "json"
	hook := &github.Hook{Name: &webHookName, Events: []string{"repository", "push", "create", "delete"}, Active: &active, Config: conf}
	hook, _, err := client.Organizations.CreateHook(context.Background(), gitOrgRef, hook)
	hookID := -1
	if err == nil {
		hookID = int(*hook.ID)
	}

	return hookID, err
}

func (githubApi *GithubApi) DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int) error {
	client, _ := githubApi.getClient(gitSource, user)
	_, err := client.Organizations.DeleteHook(context.Background(), gitOrgRef, int64(webHookID))
	return err
}

func (githubApi *GithubApi) GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error) {
	client, _ := githubApi.getClient(gitSource, user)

	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, err := client.Repositories.ListByOrg(context.Background(), gitOrgRef, opt)

	retVal := make([]string, 0)

	for _, repo := range repos {
		retVal = append(retVal, *repo.Name)
	}

	return &retVal, err
}

func (githubApi *GithubApi) CheckOrganizationExists(gitSource *model.GitSource, user *model.User, gitOrgRef string) bool {
	client, _ := githubApi.getClient(gitSource, user)
	_, _, err := client.Organizations.Get(context.Background(), gitOrgRef)
	return err == nil
}

func (githubApi *GithubApi) GetRepositoryTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	client, _ := githubApi.getClient(gitSource, user)
	teams, _, err := client.Repositories.ListTeams(context.Background(), gitOrgRef, repositoryRef, nil)

	retVal := make([]dto.TeamResponseDto, 0)
	for _, team := range teams {
		retVal = append(retVal, dto.TeamResponseDto{ID: int(*team.ID), Name: *team.Name, Permission: *team.Permission})
	}

	return &retVal, err
}

func (githubApi *GithubApi) GetOrganizationTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	client, _ := githubApi.getClient(gitSource, user)
	teams, _, err := client.Teams.ListTeams(context.Background(), gitOrgRef, nil)

	retVal := make([]dto.TeamResponseDto, 0)
	for _, team := range teams {
		retVal = append(retVal, dto.TeamResponseDto{ID: int(*team.ID), Name: *team.Name, Permission: *team.Permission})
	}

	return &retVal, err
}

func (githubApi *GithubApi) GetTeamMembers(gitSource *model.GitSource, user *model.User, organizationName string, teamId int) (*[]dto.UserTeamResponseDto, error) {
	client, _ := githubApi.getClient(gitSource, user)
	users, _, err := client.Teams.ListTeamMembers(context.Background(), int64(teamId), nil)

	retVal := make([]dto.UserTeamResponseDto, 0)
	for _, user := range users {
		if err == nil {
			retVal = append(retVal, dto.UserTeamResponseDto{ID: int(*user.ID), Username: *user.Name})
		}
	}

	return &retVal, err
}

func (githubApi *GithubApi) GetOrganizationMembers(gitSource *model.GitSource, user *model.User, organizationName string) (*[]GitHubUser, error) {
	client, _ := githubApi.getClient(gitSource, user)
	users, _, err := client.Organizations.ListMembers(context.Background(), organizationName, nil)

	retVal := make([]GitHubUser, 0)

	for _, user := range users {
		userMembership, _, err := client.Organizations.GetOrgMembership(context.Background(), *user.Login, organizationName)
		if err == nil {
			var role string
			if strings.Compare(*userMembership.Role, "admin") == 0 {
				role = "owner"
			} else {
				role = "member"
			}

			retVal = append(retVal, GitHubUser{ID: int(*user.ID), Username: *user.Login, Role: role, Email: *user.Email})
		}
	}

	return &retVal, err
}

func (githubApi *GithubApi) GetRepositoryMembers(gitSource *model.GitSource, user *model.User, organizationName string, repositoryRef string) (*[]GitHubUser, error) {
	client, _ := githubApi.getClient(gitSource, user)
	users, _, err := client.Repositories.ListCollaborators(context.Background(), organizationName, repositoryRef, nil)

	retVal := make([]GitHubUser, 0)

	for _, user := range users {
		userMembership, _, err := client.Organizations.GetOrgMembership(context.Background(), *user.Login, organizationName)
		if err == nil {
			var role string
			if strings.Compare(*userMembership.Role, "admin") == 0 {
				role = "owner"
			} else {
				role = "member"
			}

			retVal = append(retVal, GitHubUser{ID: int(*user.ID), Username: *user.Login, Role: role, Email: *user.Email})
		}
	}

	return &retVal, err
}

func (githubApi *GithubApi) GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) map[string]bool {
	client, _ := githubApi.getClient(gitSource, user)
	branchList, _, err := client.Repositories.ListBranches(context.Background(), gitOrgRef, repositoryRef, nil)

	retVal := make(map[string]bool)

	if err != nil {
		for _, branche := range branchList {
			retVal[*branche.Name] = true
		}
	}

	return retVal
}

func (githubApi *GithubApi) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error) {
	client, _ := githubApi.getClient(gitSource, user)
	branchList, _, err := client.Repositories.ListBranches(context.Background(), gitOrgRef, repositoryRef, nil)

	if err != nil {
		return false, err
	}

	for _, branch := range branchList {
		if err != nil {
			return false, err
		}

		tree, _, err := client.Git.GetTree(context.Background(), gitOrgRef, repositoryRef, *branch.Commit.SHA, true)
		if err != nil {
			return false, err
		}

		for _, file := range tree.Entries {
			if strings.Compare(*file.Type, "blob") == 0 && (strings.Compare(*file.Path, ".agola/config.jsonnet") == 0 || strings.Compare(*file.Path, ".agola/config.yml") == 0 || strings.Compare(*file.Path, ".agola/config.json") == 0) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (githubApi *GithubApi) GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	client, _ := githubApi.getClient(gitSource, user)
	commit, _, err := client.Repositories.GetCommit(context.Background(), gitOrgRef, repositoryRef, commitSha)
	if err != nil {
		return nil, err
	}

	author := make(map[string]string)
	author["email"] = *commit.Commit.Author.Email
	retVal := &dto.CommitMetadataDto{Sha: *commit.SHA, Author: author}

	if commit.Parents != nil {
		retVal.Parents = make([]dto.CommitParentDto, 0)
		for _, parent := range commit.Parents {
			retVal.Parents = append(retVal.Parents, dto.CommitParentDto{Sha: *parent.SHA})
		}
	}

	return retVal, nil
}

func (githubApi *GithubApi) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) *dto.OrganizationDto {
	client, _ := githubApi.getClient(gitSource, user)
	org, _, err := client.Organizations.Get(context.Background(), gitOrgRef)
	if err != nil {
		return nil
	}

	return &dto.OrganizationDto{Name: *org.Name, ID: int(*org.ID), AvatarURL: *org.AvatarURL}
}

func (githubApi *GithubApi) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]string, error) {
	client, _ := githubApi.getClient(gitSource, user)
	organizations, _, err := client.Organizations.List(context.Background(), "", nil)

	if err != nil {
		return nil, err
	}

	retVal := make([]string, 0)
	for _, org := range organizations {
		retVal = append(retVal, *org.Login)
	}

	return &retVal, nil
}

func (githubApi *GithubApi) GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	client, _ := githubApi.getClient(gitSource, user)
	userInfo, _, err := client.Users.Get(context.Background(), "")

	if err != nil {
		return nil, err
	}

	response := &dto.UserInfoDto{
		ID:        *userInfo.ID,
		Login:     *userInfo.Login,
		Email:     *userInfo.Email,
		FullName:  *userInfo.Name,
		AvatarURL: *userInfo.AvatarURL,
		IsAdmin:   *userInfo.SiteAdmin,
	}

	return response, nil
}

func (githubApi *GithubApi) getClient(gitSource *model.GitSource, user *model.User) (*github.Client, error) {
	if common.IsAccessTokenExpired(user.Oauth2AccessTokenExpiresAt) {
		token, err := githubApi.RefreshToken(gitSource, user.Oauth2RefreshToken)

		if err != nil {
			log.Println("error during refresh token")
			return nil, err
		}

		user.Oauth2AccessToken = token.AccessToken
		user.Oauth2RefreshToken = token.RefreshToken
		user.Oauth2AccessTokenExpiresAt = token.Expiry

		githubApi.Db.SaveUser(user)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken:  user.Oauth2AccessToken,
			TokenType:    "bearer",
			RefreshToken: user.Oauth2RefreshToken,
			Expiry:       user.Oauth2AccessTokenExpiresAt},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

const oauth2AuthorizePath string = "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=&state=%s"
const oauth2AccessTokenPath string = "https://github.com/login/oauth/access_token"

func GetOauth2AuthorizeUrl(gitClientId string, redirectUrl string, state string) string {
	return fmt.Sprintf(oauth2AuthorizePath, gitClientId, redirectUrl, state)
}

func (githubApi *GithubApi) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*oauth2.Token, error) {
	client := &http.Client{}

	URLApi := oauth2AccessTokenPath
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

func (githubApi *GithubApi) RefreshToken(gitSource *model.GitSource, refreshToken string) (*oauth2.Token, error) {
	client := &http.Client{}

	URLApi := oauth2AccessTokenPath
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
