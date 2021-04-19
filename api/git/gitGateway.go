package git

import (
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/model"
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

func GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error) {
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

func GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	if gitSource.GitType == model.Gitea {
		return gitea.GetOrganizationTeams(gitSource, gitOrgRef)
	} else {
		return github.GetOrganizationTeams(gitSource, gitOrgRef)
	}
}

func GetTeamMembers(gitSource *model.GitSource, organizationName string, teamId int) (*[]dto.UserTeamResponseDto, error) {
	if gitSource.GitType == model.Gitea {
		return gitea.GetTeamMembers(gitSource, teamId)
	} else {
		return github.GetTeamMembers(gitSource, organizationName, teamId)
	}
}

func CheckRepositoryAgolaConf(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (bool, error) {
	if gitSource.GitType == model.Gitea {
		return gitea.CheckRepositoryAgolaConfExists(gitSource, gitOrgRef, repositoryRef)
	} else {
		return github.CheckRepositoryAgolaConfExists(gitSource, gitOrgRef, repositoryRef)
	}
}

func GetCommitMetadata(gitSource *model.GitSource, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	if gitSource.GitType == model.Gitea {
		return gitea.GetCommitMetadata(gitSource, gitOrgRef, repositoryRef, commitSha)
	} else {
		return github.GetCommitMetadata(gitSource, gitOrgRef, repositoryRef, commitSha)
	}
}

func GetRepositoryTeams(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	if gitSource.GitType == model.Gitea {
		return gitea.GetRepositoryTeams(gitSource, gitOrgRef, repositoryRef)
	} else {
		return github.GetRepositoryTeams(gitSource, gitOrgRef, repositoryRef)
	}
}

func GetBranches(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) map[string]bool {
	if gitSource.GitType == model.Gitea {
		return gitea.GetBranches(gitSource, gitOrgRef, repositoryRef)
	} else {
		return github.GetBranches(gitSource, gitOrgRef, repositoryRef)
	}
}

func GetOrganization(gitSource *model.GitSource, gitOrgRef string) *dto.OrganizationDto {
	if gitSource.GitType == model.Gitea {
		return gitea.GetOrganization(gitSource, gitOrgRef)
	} else {
		return github.GetOrganization(gitSource, gitOrgRef)
	}
}
