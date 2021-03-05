package git

import "fmt"

const createWebHookPath string = "%s/api/v1/orgs/%s/hooks/?token=%s"

func getCreateWebHookUrl(gitApiUrl string, gitOrgRef string, gitToken string) string {
	return fmt.Sprintf(createWebHookPath, gitApiUrl, gitOrgRef, gitToken)
}
