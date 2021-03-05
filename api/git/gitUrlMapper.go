package git

import "fmt"

const createWebHookPath string = "%s/api/v1/orgs/%s/hooks?token=%s"
const listRepositoryesPath string = "%s/api/v1/orgs/%s/repos?token=%s"

func getCreateWebHookUrl(gitApiUrl string, gitOrgRef string, gitToken string) string {
	return fmt.Sprintf(createWebHookPath, gitApiUrl, gitOrgRef, gitToken)
}

func getGetListRepositoryPath(gitApiUrl string, gitOrgRef string, gitToken string) string {
	return fmt.Sprintf(listRepositoryesPath, gitApiUrl, gitOrgRef, gitToken)
}
