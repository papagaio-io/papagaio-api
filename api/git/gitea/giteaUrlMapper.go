package gitea

import "fmt"

const createWebHookPath string = "%s/api/v1/orgs/%s/hooks?token=%s"
const webhookPath string = "%s/api/v1/orgs/%s/hooks/%s?token=%s"
const listRepositoryesPath string = "%s/api/v1/orgs/%s/repos?token=%s"
const organizationPath string = "%s/api/v1/orgs/%s?token=%s"
const organizationTeamsListPath string = "%s/api/v1/orgs/%s/teams?token=%s"
const repositoryTeamsListPath string = "%s/api/v1/orgs/%s/%s/teams?token=%s"
const teamUsersListPath string = "%s/api/v1/teams/%s/members?token=%s"
const listBranchPath string = "%s/api/v1/repos/%s/%s/branches?token=%s"
const listMetadataPath string = "%s/api/v1/repos/%s/%s/contents/%s?ref=%s&token=%s"
const commitMetadataPath string = "%s/api/v1/repos/%s/%s/commits?sha=%s&page=1&limit=1&token=%s"
const repositoryListPath string = "%s/api/v1/orgs/%s/repos?token=%s"

func getCreateWebHookUrl(gitApiUrl string, gitOrgRef string, gitToken string) string {
	return fmt.Sprintf(createWebHookPath, gitApiUrl, gitOrgRef, gitToken)
}

func getWehHookUrl(gitApiUrl string, gitOrgRef string, webHookID string, gitToken string) string {
	return fmt.Sprintf(webhookPath, gitApiUrl, gitOrgRef, webHookID, gitToken)
}

func getGetListRepositoryUrl(gitApiUrl string, gitOrgRef string, gitToken string) string {
	return fmt.Sprintf(listRepositoryesPath, gitApiUrl, gitOrgRef, gitToken)
}

func getOrganizationUrl(gitApiUrl string, gitOrgRef string, gitToken string) string {
	return fmt.Sprintf(organizationPath, gitApiUrl, gitOrgRef, gitToken)
}

func getOrganizationTeamsListUrl(gitApiUrl string, gitOrgRef string, gitToken string) string {
	return fmt.Sprintf(organizationTeamsListPath, gitApiUrl, gitOrgRef, gitToken)
}

func getRepositoryTeamsListUrl(gitApiUrl string, gitOrgRef string, repositoryRef string, gitToken string) string {
	return fmt.Sprintf(repositoryTeamsListPath, gitApiUrl, gitOrgRef, repositoryRef, gitToken)
}

func getTeamUsersListUrl(gitApiUrl string, teamId string, gitToken string) string {
	return fmt.Sprintf(teamUsersListPath, gitApiUrl, teamId, gitToken)
}

func getListBranchPath(gitApiUrl string, gitOrgRef string, repositoryRef string, gitToken string) string {
	return fmt.Sprintf(listBranchPath, gitApiUrl, gitOrgRef, repositoryRef, gitToken)
}

func getListMetadataPath(gitApiUrl string, gitOrgRef string, repositoryRef string, dirSourcePath string, branch string, token string) string {
	return fmt.Sprintf(listMetadataPath, gitApiUrl, gitOrgRef, repositoryRef, dirSourcePath, branch, token)
}

func getCommitMetadataPath(gitApiUrl string, gitOrgRef string, repositoryRef string, commitSha string, token string) string {
	return fmt.Sprintf(commitMetadataPath, gitApiUrl, gitOrgRef, repositoryRef, commitSha, token)
}
