package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type GitlabInterface interface {
	CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int, error)
	DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int) error
	GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error)
	GetEmailsRepositoryUsersOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]string, error)
	GetOrganizationMembers(gitSource *model.GitSource, user *model.User, organizationName string) (*[]GitlabUser, error)
	GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) map[string]bool
	CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error)
	GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error)
	GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) *dto.OrganizationDto
	GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error)
	IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error)

	GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error)
	GetUserByLogin(gitSource *model.GitSource, id int) (*dto.UserInfoDto, error)

	GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error)
	RefreshToken(gitSource *model.GitSource, refreshToken string) (*common.Token, error)
}

type GitlabApi struct {
	Db repository.Database
}

func (gitlabApi *GitlabApi) CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int, error) {
	client, _ := gitlabApi.getClient(gitSource, user)

	url := config.Config.Server.LocalHostAddress + controller.GetWebHookPath() + "/" + organizationRef
	pushEvents := true
	groupHook, _, err := client.Groups.AddGroupHook(gitOrgRef, &gitlab.AddGroupHookOptions{
		URL:        &url,
		PushEvents: &pushEvents,
	})
	hookID := -1
	if err == nil {
		hookID = groupHook.ID
	}

	return hookID, err
}

func (gitlabApi *GitlabApi) DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int) error {
	client, _ := gitlabApi.getClient(gitSource, user)
	_, err := client.Groups.DeleteGroupHook(gitOrgRef, webHookID)
	return err
}

func (gitlabApi *GitlabApi) GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error) {
	client, _ := gitlabApi.getClient(gitSource, user)

	projectList, _, err := client.Groups.ListGroupProjects(gitOrgRef, nil)

	retVal := make([]string, 0)

	for _, project := range projectList {
		retVal = append(retVal, project.Name)
	}

	return &retVal, err
}

func (gitlabApi *GitlabApi) GetEmailsRepositoryUsersOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]string, error) {
	client, _ := gitlabApi.getClient(gitSource, user)
	users, _, err := client.ProjectMembers.ListAllProjectMembers(gitOrgRef+"/"+repositoryRef, nil)
	if err != nil {
		return nil, err
	}

	retVal := make([]string, 0)

	for _, user := range users {
		if user.AccessLevel == gitlab.OwnerPermissions || user.AccessLevel == gitlab.MaintainerPermissions {
			retVal = append(retVal, user.Email)
		}
	}

	return &retVal, nil
}

func (gitlabApi *GitlabApi) GetOrganizationMembers(gitSource *model.GitSource, user *model.User, organizationName string) (*[]GitlabUser, error) {
	client, _ := gitlabApi.getClient(gitSource, user)
	members, _, err := client.Groups.ListAllGroupMembers(organizationName, nil)
	if err != nil {
		return nil, err
	}

	retVal := make([]GitlabUser, 0)

	for _, member := range members {
		user := GitlabUser{
			ID:          member.ID,
			Username:    member.Username,
			AccessLevel: member.AccessLevel,
		}

		retVal = append(retVal, user)
	}

	return &retVal, nil
}

func (gitlabApi *GitlabApi) GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) map[string]bool {
	client, _ := gitlabApi.getClient(gitSource, user)
	branches, _, _ := client.Branches.ListBranches(gitOrgRef+"/"+repositoryRef, nil)

	retVal := make(map[string]bool)

	for _, branche := range branches {
		retVal[branche.Name] = true
	}

	return retVal
}

func (gitlabApi *GitlabApi) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error) {
	client, _ := gitlabApi.getClient(gitSource, user)
	branchList, _, err := client.Branches.ListBranches(gitOrgRef+"/"+repositoryRef, nil)

	if err != nil {
		return false, err
	}

	for _, branch := range branchList {
		if err != nil {
			return false, err
		}

		optionPath := ".agola"
		options := gitlab.ListTreeOptions{Ref: &branch.Name, Path: &optionPath}
		tree, _, err := client.Repositories.ListTree(gitOrgRef+"/"+repositoryRef, &options)
		if err != nil {
			return false, err
		}

		for _, file := range tree {
			if strings.Compare(file.Type, "blob") == 0 && (strings.Compare(file.Path, ".agola/config.jsonnet") == 0 || strings.Compare(file.Path, ".agola/config.yml") == 0 || strings.Compare(file.Path, ".agola/config.json") == 0) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (gitlabApi *GitlabApi) GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	client, _ := gitlabApi.getClient(gitSource, user)
	commit, _, err := client.Commits.GetCommit(gitOrgRef+"/"+repositoryRef, commitSha)
	if err != nil {
		return nil, err
	}

	author := make(map[string]string)
	author["email"] = commit.CommitterEmail
	retVal := &dto.CommitMetadataDto{Sha: commitSha, Author: author}

	retVal.Parents = make([]dto.CommitParentDto, 0)
	for _, parent := range commit.ParentIDs {
		retVal.Parents = append(retVal.Parents, dto.CommitParentDto{Sha: parent})
	}

	return retVal, nil
}

func (gitlabApi *GitlabApi) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) *dto.OrganizationDto {
	client, _ := gitlabApi.getClient(gitSource, user)
	org, _, _ := client.Groups.GetGroup(gitOrgRef)
	if org == nil {
		return nil
	}

	response := &dto.OrganizationDto{Name: org.Name, Path: org.Path, ID: int64(org.ID), AvatarURL: org.AvatarURL}

	return response
}

func (gitlabApi *GitlabApi) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error) {
	client, _ := gitlabApi.getClient(gitSource, user)
	minAccessLevel := gitlab.OwnerPermissions
	organizations, _, err := client.Groups.ListGroups(&gitlab.ListGroupsOptions{MinAccessLevel: &minAccessLevel})

	if err != nil {
		return nil, err
	}

	retVal := make([]dto.OrganizationDto, 0)
	for _, org := range organizations {
		orgDto := dto.OrganizationDto{Name: org.Name, Path: org.Path, ID: int64(org.ID), AvatarURL: org.AvatarURL}
		retVal = append(retVal, orgDto)
	}

	return &retVal, nil
}

func (gitlabApi *GitlabApi) IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error) {
	gitlabUsers, _ := gitlabApi.GetOrganizationMembers(gitSource, user, gitOrgRef)

	if gitlabUsers != nil {
		for _, u := range *gitlabUsers {
			if u.ID == int(user.ID) {
				return u.HasOwnerPermission(), nil
			}
		}
	}

	return false, nil
}

func (gitlabApi *GitlabApi) GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	client, _ := gitlabApi.getClient(gitSource, user)
	userInfo, _, err := client.Users.CurrentUser()

	if err != nil {
		return nil, err
	}

	response := &dto.UserInfoDto{
		ID:          int64(userInfo.ID),
		Login:       userInfo.Username,
		AvatarURL:   userInfo.AvatarURL,
		IsAdmin:     userInfo.IsAdmin,
		UserPageURL: userInfo.WebURL,
		FullName:    userInfo.Name,
		Email:       userInfo.Email,
	}

	return response, nil
}

func (gitlabApi *GitlabApi) GetUserByLogin(gitSource *model.GitSource, id int) (*dto.UserInfoDto, error) {
	client, _ := gitlabApi.getClient(gitSource, nil)
	user, resp, err := client.Users.GetUser(id, gitlab.GetUsersOptions{})

	if err != nil {
		log.Println("error in GetUserByLogin:", err)
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	userInfo := dto.UserInfoDto{
		ID:          int64(user.ID),
		Login:       user.Username,
		AvatarURL:   user.AvatarURL,
		IsAdmin:     user.IsAdmin,
		UserPageURL: user.WebURL,
		FullName:    user.Name,
		Email:       user.Email,
	}

	return &userInfo, nil
}

func (gitlabApi *GitlabApi) getClient(gitSource *model.GitSource, user *model.User) (*gitlab.Client, error) {
	if user == nil {
		return gitlab.NewClient("")
	}

	if common.IsAccessTokenExpired(user.Oauth2AccessTokenExpiresAt) {
		log.Println("Token expired is to refresh")
		token, err := gitlabApi.RefreshToken(gitSource, user.Oauth2RefreshToken)

		if err != nil {
			log.Println("error during refresh token")
			return nil, err
		}

		user.Oauth2AccessToken = token.AccessToken
		user.Oauth2RefreshToken = token.RefreshToken
		user.Oauth2AccessTokenExpiresAt = token.ExpiryAt

		gitlabApi.Db.SaveUser(user)
	}

	return gitlab.NewOAuthClient(user.Oauth2AccessToken)
}

const oauth2AuthorizePath string = "https://gitlab.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=%s"

const oauth2AccessTokenPath string = "https://gitlab.com/oauth/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s"
const oauth2RefreshTokenPath string = "https://gitlab.com/oauth/token?client_id=%s&client_secret=%s&grant_type=refresh_token&refresh_token=%s"

func GetOauth2AuthorizeUrl(gitClientId string, redirectUrl string, state string) string {
	return fmt.Sprintf(oauth2AuthorizePath, gitClientId, url.QueryEscape(redirectUrl), state, "api%20read_repository%20read_api")
}

func (gitlabApi *GitlabApi) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error) {
	client := &http.Client{}

	URLApi := fmt.Sprintf(oauth2AccessTokenPath, gitSource.GitClientID, gitSource.GitSecret, code, url.QueryEscape(controller.GetRedirectUrl()))
	req, _ := http.NewRequest("POST", URLApi, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response common.Token
		json.Unmarshal(body, &response)

		if response.Expiry > 0 {
			response.ExpiryAt = time.Now().Add(time.Second * time.Duration(response.Expiry))
		}

		return &response, nil
	}

	return nil, err
}

func (gitlabApi *GitlabApi) RefreshToken(gitSource *model.GitSource, refreshToken string) (*common.Token, error) {
	client := &http.Client{}

	URLApi := fmt.Sprintf(oauth2RefreshTokenPath, gitSource.GitClientID, gitSource.GitSecret, refreshToken)
	req, _ := http.NewRequest("POST", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		var response common.Token
		json.Unmarshal(body, &response)

		response.ExpiryAt = time.Now().Add(time.Second * time.Duration(response.Expiry))

		return &response, nil
	}

	return nil, err
}
