package git

import (
	"wecode.sorint.it/opensource/papagaio-be/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-be/api/git/github"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func CreateWebHook(gitSource *model.GitSource, gitOrgRef string) (int, error) {
	if gitSource.GitType == model.Gitea {
		return gitea.CreateWebHook(gitSource, gitOrgRef)
	} else {
		return github.CreateWebHook(gitSource, gitOrgRef)
	}
}

func DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error {
	if gitSource.GitType == model.Gitea {
		return gitea.DeleteWebHook(gitSource, gitOrgRef, webHookID)
	} else {
		return github.DeleteWebHook(gitSource, gitOrgRef, webHookID)
	}
}

func GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error) { //*[]gitea.RepositoryDto
	if gitSource.GitType == model.Gitea {
		return gitea.GetRepositories(gitSource, gitOrgRef)
	} else {
		return github.GetRepositories(gitSource, gitOrgRef)
	}
}

func CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool {
	if gitSource.GitType == model.Gitea {
		return gitea.CheckOrganizationExists(gitSource, gitOrgRef)
	} else {
		return github.CheckOrganizationExists(gitSource, gitOrgRef)
	}
}

func GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]gitea.TeamResponseDto, error) {
	if gitSource.GitType == model.Gitea {
		return gitea.GetOrganizationTeams(gitSource, gitOrgRef)
	} else {
		return github.GetOrganizationTeams(gitSource, gitOrgRef)
	}

	return nil, nil
}

func GetTeamMembers(gitSource *model.GitSource, organizationName string, teamId int) (*[]gitea.UserTeamResponseDto, error) {
	if gitSource.GitType == model.Gitea {
		return gitea.GetTeamMembers(gitSource, teamId)
	} else {
		return github.GetTeamMembers(gitSource, organizationName, teamId)
	}

	return nil, nil
}
