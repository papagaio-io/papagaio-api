package gitea

import (
	"fmt"
	"net/url"
)

const createWebHookPath string = "%s/api/v1/orgs/%s/hooks"
const webhookPath string = "%s/api/v1/orgs/%s/hooks/%s"
const listRepositoryesPath string = "%s/api/v1/orgs/%s/repos"
const organizationPath string = "%s/api/v1/orgs/%s?"
const organizationTeamsListPath string = "%s/api/v1/orgs/%s/teams"
const repositoryTeamsListPath string = "%s/api/v1/orgs/%s/%s/teams"
const teamUsersListPath string = "%s/api/v1/teams/%s/members?"
const listBranchPath string = "%s/api/v1/repos/%s/%s/branches"
const listMetadataPath string = "%s/api/v1/repos/%s/%s/contents/%s?ref=%s"
const commitMetadataPath string = "%s/api/v1/repos/%s/%s/commits?sha=%s&page=1&limit=1"
const organizationsPath string = "%s/api/v1/user/orgs"
const loggedUserInfoPath string = "%s/api/v1/user"
const userInfoPath string = "%s/api/v1/users/%s"
const createOauth2AppPath string = "%s/api/v1/user/applications/oauth2"

const oauth2AuthorizePath string = "%s/login/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=&state=%s"
const oauth2AccessTokenPath string = "%s/login/oauth/access_token"

func getCreateWebHookUrl(gitApiUrl string, gitOrgRef string) string {
	return fmt.Sprintf(createWebHookPath, gitApiUrl, gitOrgRef)
}

func getWehHookUrl(gitApiUrl string, gitOrgRef string, webHookID string) string {
	return fmt.Sprintf(webhookPath, gitApiUrl, gitOrgRef, webHookID)
}

func getGetListRepositoryUrl(gitApiUrl string, gitOrgRef string) string {
	return fmt.Sprintf(listRepositoryesPath, gitApiUrl, gitOrgRef)
}

func getOrganizationUrl(gitApiUrl string, gitOrgRef string) string {
	return fmt.Sprintf(organizationPath, gitApiUrl, gitOrgRef)
}

func getOrganizationTeamsListUrl(gitApiUrl string, gitOrgRef string) string {
	return fmt.Sprintf(organizationTeamsListPath, gitApiUrl, gitOrgRef)
}

func getRepositoryTeamsListUrl(gitApiUrl string, gitOrgRef string, repositoryRef string) string {
	return fmt.Sprintf(repositoryTeamsListPath, gitApiUrl, gitOrgRef, repositoryRef)
}

func getTeamUsersListUrl(gitApiUrl string, teamId string) string {
	return fmt.Sprintf(teamUsersListPath, gitApiUrl, teamId)
}

func getListBranchUrl(gitApiUrl string, gitOrgRef string, repositoryRef string) string {
	return fmt.Sprintf(listBranchPath, gitApiUrl, gitOrgRef, repositoryRef)
}

func getListMetadataUrl(gitApiUrl string, gitOrgRef string, repositoryRef string, dirSourcePath string, branch string) string {
	return fmt.Sprintf(listMetadataPath, gitApiUrl, gitOrgRef, repositoryRef, dirSourcePath, branch)
}

func getCommitMetadataPath(gitApiUrl string, gitOrgRef string, repositoryRef string, commitSha string) string {
	return fmt.Sprintf(commitMetadataPath, gitApiUrl, gitOrgRef, repositoryRef, commitSha)
}

func getOrganizationsUrl(gitApiUrl string) string {
	return fmt.Sprintf(organizationsPath, gitApiUrl)
}

func getLoggedUserInfoUrl(gitApiUrl string) string {
	return fmt.Sprintf(loggedUserInfoPath, gitApiUrl)
}

func getUserInfoUrl(gitApiUrl string, login string) string {
	return fmt.Sprintf(userInfoPath, gitApiUrl, login)
}

func GetOauth2AuthorizeUrl(gitApiUrl string, gitClientId string, redirectUrl string, state string) string {
	return fmt.Sprintf(oauth2AuthorizePath, gitApiUrl, gitClientId, url.QueryEscape(redirectUrl), state)
}

func getOauth2AccessTokenUrl(gitApiUrl string) string {
	return fmt.Sprintf(oauth2AccessTokenPath, gitApiUrl)
}

func getCreateOauth2AppUrl(gitApiUrl string) string {
	return fmt.Sprintf(createOauth2AppPath, gitApiUrl)
}
