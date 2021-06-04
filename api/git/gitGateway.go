package git

import (
	"golang.org/x/oauth2"
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

func (gitGateway *GitGateway) CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.CreateWebHook(gitSource, user, gitOrgRef, organizationRef)
	} else {
		return gitGateway.GithubApi.CreateWebHook(gitSource, user, gitOrgRef, organizationRef)
	}
}

func (gitGateway *GitGateway) DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int) error {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.DeleteWebHook(gitSource, user, gitOrgRef, webHookID)
	} else {
		return gitGateway.GithubApi.DeleteWebHook(gitSource, user, gitOrgRef, webHookID)
	}
}

func (gitGateway *GitGateway) GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetRepositories(gitSource, user, gitOrgRef)
	} else {
		return gitGateway.GithubApi.GetRepositories(gitSource, user, gitOrgRef)
	}
}

func (gitGateway *GitGateway) CheckOrganizationExists(gitSource *model.GitSource, user *model.User, gitOrgRef string) bool {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.CheckOrganizationExists(gitSource, user, gitOrgRef)
	} else {
		return gitGateway.GithubApi.CheckOrganizationExists(gitSource, user, gitOrgRef)
	}
}

func (gitGateway *GitGateway) GetOrganizationTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOrganizationTeams(gitSource, user, gitOrgRef)
	} else {
		return gitGateway.GithubApi.GetOrganizationTeams(gitSource, user, gitOrgRef)
	}
}

func (gitGateway *GitGateway) GetTeamMembers(gitSource *model.GitSource, user *model.User, organizationName string, teamId int) (*[]dto.UserTeamResponseDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetTeamMembers(gitSource, user, teamId)
	} else {
		return gitGateway.GithubApi.GetTeamMembers(gitSource, user, organizationName, teamId)
	}
}

func (gitGateway *GitGateway) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.CheckRepositoryAgolaConfExists(gitSource, user, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GithubApi.CheckRepositoryAgolaConfExists(gitSource, user, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetCommitMetadata(gitSource, user, gitOrgRef, repositoryRef, commitSha)
	} else {
		return gitGateway.GithubApi.GetCommitMetadata(gitSource, user, gitOrgRef, repositoryRef, commitSha)
	}
}

func (gitGateway *GitGateway) GetRepositoryTeams(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetRepositoryTeams(gitSource, user, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GithubApi.GetRepositoryTeams(gitSource, user, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) map[string]bool {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetBranches(gitSource, user, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GithubApi.GetBranches(gitSource, user, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) *dto.OrganizationDto {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOrganization(gitSource, user, gitOrgRef)
	} else {
		return gitGateway.GithubApi.GetOrganization(gitSource, user, gitOrgRef)
	}
}

//TODO ritornare solo le organizations dove l'utente sia owner
func (gitGateway *GitGateway) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]string, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOrganizations(gitSource, user)
	} else {
		return gitGateway.GithubApi.GetOrganizations(gitSource, user)
	}
}

func (gitGateway *GitGateway) GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetUserInfo(gitSource, user)
	} else {
		return gitGateway.GithubApi.GetUserInfo(gitSource, user)
	}
}

func (gitGateway *GitGateway) GetOauth2AuthorizePathUrl(gitSource *model.GitSource, redirectUrl string, state string) string {
	if gitSource.GitType == types.Gitea {
		return gitea.GetOauth2AuthorizeUrl(gitSource.GitAPIURL, gitSource.GitClientID, redirectUrl, state)
	} else {
		return github.GetOauth2AuthorizeUrl(gitSource.GitClientID, redirectUrl, state)
	}
}

func (gitGateway *GitGateway) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*oauth2.Token, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOauth2AccessToken(gitSource, code)
	} else {
		return gitGateway.GithubApi.GetOauth2AccessToken(gitSource, code)
	}
}

func (gitGateway *GitGateway) RefreshToken(gitSource *model.GitSource, refreshToken string) (*oauth2.Token, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.RefreshToken(gitSource, refreshToken)
	} else {
		return gitGateway.GithubApi.RefreshToken(gitSource, refreshToken)
	}
}
