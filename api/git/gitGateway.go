package git

import (
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

type GitGateway struct {
	GiteaApi  gitea.GiteaInterface
	GithubApi github.GithubInterface
}

func (gitGateway *GitGateway) CreateWebHook(gitSource *model.GitSource, gitOrgRef string, organizationRef string) (int, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.CreateWebHook(gitSource, gitOrgRef, organizationRef)
	} else {
		return gitGateway.GithubApi.CreateWebHook(gitSource, gitOrgRef, organizationRef)
	}
}

func (gitGateway *GitGateway) DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.DeleteWebHook(gitSource, gitOrgRef, webHookID)
	} else {
		return gitGateway.GithubApi.DeleteWebHook(gitSource, gitOrgRef, webHookID)
	}
}

func (gitGateway *GitGateway) GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetRepositories(gitSource, gitOrgRef)
	} else {
		return gitGateway.GithubApi.GetRepositories(gitSource, gitOrgRef)
	}
}

func (gitGateway *GitGateway) CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.CheckOrganizationExists(gitSource, gitOrgRef)
	} else {
		return gitGateway.GithubApi.CheckOrganizationExists(gitSource, gitOrgRef)
	}
}

func (gitGateway *GitGateway) GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOrganizationTeams(gitSource, gitOrgRef)
	} else {
		return gitGateway.GithubApi.GetOrganizationTeams(gitSource, gitOrgRef)
	}
}

func (gitGateway *GitGateway) GetTeamMembers(gitSource *model.GitSource, organizationName string, teamId int) (*[]dto.UserTeamResponseDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetTeamMembers(gitSource, teamId)
	} else {
		return gitGateway.GithubApi.GetTeamMembers(gitSource, organizationName, teamId)
	}
}

func (gitGateway *GitGateway) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (bool, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.CheckRepositoryAgolaConfExists(gitSource, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GithubApi.CheckRepositoryAgolaConfExists(gitSource, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) GetCommitMetadata(gitSource *model.GitSource, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetCommitMetadata(gitSource, gitOrgRef, repositoryRef, commitSha)
	} else {
		return gitGateway.GithubApi.GetCommitMetadata(gitSource, gitOrgRef, repositoryRef, commitSha)
	}
}

func (gitGateway *GitGateway) GetRepositoryTeams(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetRepositoryTeams(gitSource, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GithubApi.GetRepositoryTeams(gitSource, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) GetBranches(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) map[string]bool {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetBranches(gitSource, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GithubApi.GetBranches(gitSource, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) GetOrganization(gitSource *model.GitSource, gitOrgRef string) *dto.OrganizationDto {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOrganization(gitSource, gitOrgRef)
	} else {
		return gitGateway.GithubApi.GetOrganization(gitSource, gitOrgRef)
	}
}

func (gitGateway *GitGateway) GetOrganizations(gitSource *model.GitSource) (*[]string, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOrganizations(gitSource)
	} else {
		return gitGateway.GithubApi.GetOrganizations(gitSource)
	}
}
