package gitea

import "fmt"

const createWebHookPath string = "%s/api/v1/orgs/%s/hooks?token=%s"
const deleteWebHookPath string = "%s/api/v1/orgs/%s/hooks/%s?token=%s"
const listRepositoryesPath string = "%s/api/v1/orgs/%s/repos?token=%s"
const organizationPath string = "%s/api/v1/orgs/%s?token=%s"
const organizationTeamsListPath string = "%s/api/v1/orgs/%s/teams?token=%s"
const teamUsersListPath string = "%s/api/v1/teams/%s/members?token=%s"

func getCreateWebHookUrl(gitApiUrl string, gitOrgRef string, gitToken string) string {
	return fmt.Sprintf(createWebHookPath, gitApiUrl, gitOrgRef, gitToken)
}

func getDeleteWehHookUrl(gitApiUrl string, gitOrgRef string, webHookID string, gitToken string) string {
	return fmt.Sprintf(deleteWebHookPath, gitApiUrl, gitOrgRef, webHookID, gitToken)
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

func getTeamUsersListUrl(gitApiUrl string, teamId string, gitToken string) string {
	return fmt.Sprintf(teamUsersListPath, gitApiUrl, teamId, gitToken)
}
