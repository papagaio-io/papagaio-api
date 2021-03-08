package gitea

import "fmt"

const createWebHookPath string = "%s/api/v1/orgs/%s/hooks?token=%s"
const deleteWebHookPath string = "%s/api/v1/orgs/%s/hooks/%s?token=%s"
const listRepositoryesPath string = "%s/api/v1/orgs/%s/repos?token=%s"
const organizationPath string = "%s/api/v1/orgs/%s?token=%s"

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
