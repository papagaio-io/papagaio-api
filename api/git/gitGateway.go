package git

import (
	"errors"
	"strings"

	"wecode.sorint.it/opensource/papagaio-be/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func CreateWebHook(gitSource *model.GitSource, gitOrgRef string, branchFilter string) (int, error) {
	if strings.Compare(gitSource.GitType, "gitea") == 0 {
		return gitea.CreateWebHook(gitSource, gitOrgRef, branchFilter)
	}

	return -1, errors.New("Git type not found")
}

func DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error {
	if strings.Compare(gitSource.GitType, "gitea") == 0 {
		return gitea.DeleteWebHook(gitSource, gitOrgRef, webHookID)
	}

	return errors.New("Git type not found")
}

func GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]gitea.RepositoryDto, error) {
	if strings.Compare(gitSource.GitType, "gitea") == 0 {
		return gitea.GetRepositories(gitSource, gitOrgRef)
	}

	return nil, errors.New("Git type not found")
}

func CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool {
	if strings.Compare(gitSource.GitType, "gitea") == 0 {
		return gitea.CheckOrganizationExists(gitSource, gitOrgRef)
	}

	return false
}
