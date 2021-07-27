package git

import (
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitlab"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

type GitGateway struct {
	GiteaApi  gitea.GiteaInterface
	GithubApi github.GithubInterface
	GitlabApi gitlab.GitlabInterface
}

func (gitGateway *GitGateway) CreateWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, organizationRef string) (int64, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.CreateWebHook(gitSource, user, gitOrgRef, organizationRef)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.CreateWebHook(gitSource, user, gitOrgRef, organizationRef)
	} else {
		return gitGateway.GitlabApi.CreateWebHook(gitSource, user, gitOrgRef, organizationRef)
	}
}

func (gitGateway *GitGateway) DeleteWebHook(gitSource *model.GitSource, user *model.User, gitOrgRef string, webHookID int64) error {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.DeleteWebHook(gitSource, user, gitOrgRef, webHookID)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.DeleteWebHook(gitSource, user, gitOrgRef, webHookID)
	} else {
		return gitGateway.GitlabApi.DeleteWebHook(gitSource, user, gitOrgRef, webHookID)
	}
}

func (gitGateway *GitGateway) GetRepositories(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*[]string, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetRepositories(gitSource, user, gitOrgRef)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetRepositories(gitSource, user, gitOrgRef)
	} else {
		return gitGateway.GitlabApi.GetRepositories(gitSource, user, gitOrgRef)
	}
}

func (gitGateway *GitGateway) GetEmailsRepositoryUsersOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (*[]string, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetEmailsRepositoryUsersOwner(gitSource, user, gitOrgRef, repositoryRef)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetEmailsRepositoryUsersOwner(gitSource, user, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GitlabApi.GetEmailsRepositoryUsersOwner(gitSource, user, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) CheckRepositoryAgolaConfExists(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (bool, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.CheckRepositoryAgolaConfExists(gitSource, user, gitOrgRef, repositoryRef)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.CheckRepositoryAgolaConfExists(gitSource, user, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GitlabApi.CheckRepositoryAgolaConfExists(gitSource, user, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) GetCommitMetadata(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetCommitMetadata(gitSource, user, gitOrgRef, repositoryRef, commitSha)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetCommitMetadata(gitSource, user, gitOrgRef, repositoryRef, commitSha)
	} else {
		return gitGateway.GitlabApi.GetCommitMetadata(gitSource, user, gitOrgRef, repositoryRef, commitSha)
	}
}

func (gitGateway *GitGateway) GetBranches(gitSource *model.GitSource, user *model.User, gitOrgRef string, repositoryRef string) (map[string]bool, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetBranches(gitSource, user, gitOrgRef, repositoryRef)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetBranches(gitSource, user, gitOrgRef, repositoryRef)
	} else {
		return gitGateway.GitlabApi.GetBranches(gitSource, user, gitOrgRef, repositoryRef)
	}
}

func (gitGateway *GitGateway) GetOrganization(gitSource *model.GitSource, user *model.User, gitOrgRef string) (*dto.OrganizationDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOrganization(gitSource, user, gitOrgRef)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetOrganization(gitSource, user, gitOrgRef)
	} else {
		return gitGateway.GitlabApi.GetOrganization(gitSource, user, gitOrgRef)
	}
}

func (gitGateway *GitGateway) GetOrganizations(gitSource *model.GitSource, user *model.User) (*[]dto.OrganizationDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOrganizations(gitSource, user)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetOrganizations(gitSource, user)
	} else {
		return gitGateway.GitlabApi.GetOrganizations(gitSource, user)
	}
}

func (gitGateway *GitGateway) IsUserOwner(gitSource *model.GitSource, user *model.User, gitOrgRef string) (bool, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.IsUserOwner(gitSource, user, gitOrgRef)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.IsUserOwner(gitSource, user, gitOrgRef)
	} else {
		return gitGateway.GitlabApi.IsUserOwner(gitSource, user, gitOrgRef)
	}
}

func (gitGateway *GitGateway) GetUserInfo(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetUserInfo(gitSource, user)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetUserInfo(gitSource, user)
	} else {
		return gitGateway.GitlabApi.GetUserInfo(gitSource, user)
	}
}

func (gitGateway *GitGateway) GetUserByLogin(gitSource *model.GitSource, user *model.User) (*dto.UserInfoDto, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetUserByLogin(gitSource, user.Login)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetUserByLogin(gitSource, user.Login)
	} else {
		return gitGateway.GitlabApi.GetUserByLogin(gitSource, int(user.ID))
	}
}

func (gitGateway *GitGateway) GetOauth2AuthorizePathUrl(gitSource *model.GitSource, redirectUrl string, state string) string {
	if gitSource.GitType == types.Gitea {
		return gitea.GetOauth2AuthorizeUrl(gitSource.GitAPIURL, gitSource.GitClientID, redirectUrl, state)
	} else if gitSource.GitType == types.Github {
		return github.GetOauth2AuthorizeUrl(gitSource.GitClientID, redirectUrl, state)
	} else {
		return gitlab.GetOauth2AuthorizeUrl(gitSource.GitClientID, redirectUrl, state)
	}
}

func (gitGateway *GitGateway) GetOauth2AccessToken(gitSource *model.GitSource, code string) (*common.Token, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.GetOauth2AccessToken(gitSource, code)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.GetOauth2AccessToken(gitSource, code)
	} else {
		return gitGateway.GitlabApi.GetOauth2AccessToken(gitSource, code)
	}
}

func (gitGateway *GitGateway) RefreshToken(gitSource *model.GitSource, refreshToken string) (*common.Token, error) {
	if gitSource.GitType == types.Gitea {
		return gitGateway.GiteaApi.RefreshToken(gitSource, refreshToken)
	} else if gitSource.GitType == types.Github {
		return gitGateway.GithubApi.RefreshToken(gitSource, refreshToken)
	} else {
		return gitGateway.GitlabApi.RefreshToken(gitSource, refreshToken)
	}
}
