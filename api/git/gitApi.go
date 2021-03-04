package git

import "wecode.sorint.it/opensource/papagaio-be/model"

//TODO
func CreateWebHook(gitSource *model.GitSource, gitOrgRef string) (int, error) {
	var webHookID int
	var err error

	return webHookID, err
}

//TODO
func DeleteWebHook(gitSource *model.GitSource, webHookID int) error {
	var err error
	return err
}

//TODO
func GetRepositories(gitSource *model.GitSource, gitOrgRef string) ([]string, error) {
	var gitRepositoryRef []string
	var err error

	return gitRepositoryRef, err
}

//TODO
func GetGitOrganizations(gitSource *model.GitSource) ([]string, error) {
	var organizations []string
	var err error
	return organizations, err
}
