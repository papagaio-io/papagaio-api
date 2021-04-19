package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func CreateWebHook(gitSource *model.GitSource, gitOrgRef string) (int, error) {
	client := getClient(gitSource)

	webHookName := "web"
	active := true
	conf := make(map[string]interface{})
	conf["url"] = config.Config.Server.LocalHostAddress + controller.GetWebHookPath() + "/" + gitOrgRef
	fmt.Println("url:", conf["url"])
	conf["content_type"] = "json"
	hook := &github.Hook{Name: &webHookName, Events: []string{"repository", "push", "create", "delete"}, Active: &active, Config: conf}
	hook, resp, err := client.Organizations.CreateHook(context.Background(), gitOrgRef, hook)
	hookID := -1
	if err == nil {
		hookID = int(*hook.ID)
	}

	fmt.Println("resp:", resp)
	fmt.Println("err:", err)

	return hookID, err
}

func DeleteWebHook(gitSource *model.GitSource, gitOrgRef string, webHookID int) error {
	client := getClient(gitSource)
	_, err := client.Organizations.DeleteHook(context.Background(), gitOrgRef, int64(webHookID))
	return err
}

func GetRepositories(gitSource *model.GitSource, gitOrgRef string) (*[]string, error) {
	client := getClient(gitSource)

	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, err := client.Repositories.ListByOrg(context.Background(), gitOrgRef, opt)

	retVal := make([]string, 0)

	for _, repo := range repos {
		retVal = append(retVal, *repo.Name)
	}

	return &retVal, err
}

func CheckOrganizationExists(gitSource *model.GitSource, gitOrgRef string) bool {
	client := getClient(gitSource)
	_, _, err := client.Organizations.Get(context.Background(), gitOrgRef)
	return err == nil
}

func GetRepositoryTeams(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (*[]dto.TeamResponseDto, error) {
	client := getClient(gitSource)
	teams, _, err := client.Repositories.ListTeams(context.Background(), gitOrgRef, repositoryRef, nil)

	retVal := make([]dto.TeamResponseDto, 0)
	for _, team := range teams {
		retVal = append(retVal, dto.TeamResponseDto{ID: int(*team.ID), Name: *team.Name, Permission: *team.Permission})
	}

	return &retVal, err
}

func GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]dto.TeamResponseDto, error) {
	client := getClient(gitSource)
	teams, _, err := client.Teams.ListTeams(context.Background(), gitOrgRef, nil)

	retVal := make([]dto.TeamResponseDto, 0)
	for _, team := range teams {
		retVal = append(retVal, dto.TeamResponseDto{ID: int(*team.ID), Name: *team.Name, Permission: *team.Permission})
	}

	return &retVal, err
}

func GetTeamMembers(gitSource *model.GitSource, organizationName string, teamId int) (*[]dto.UserTeamResponseDto, error) {
	client := getClient(gitSource)
	users, _, err := client.Teams.ListTeamMembers(context.Background(), int64(teamId), nil)

	retVal := make([]dto.UserTeamResponseDto, 0)
	for _, user := range users {
		if err == nil {
			retVal = append(retVal, dto.UserTeamResponseDto{ID: int(*user.ID), Username: *user.Name})
		}
	}

	return &retVal, err
}

func GetOrganizationMembers(gitSource *model.GitSource, organizationName string) (*[]GitHubUser, error) {
	client := getClient(gitSource)
	users, _, err := client.Organizations.ListMembers(context.Background(), organizationName, nil)

	retVal := make([]GitHubUser, 0)

	for _, user := range users {
		userMembership, _, err := client.Organizations.GetOrgMembership(context.Background(), *user.Login, organizationName)
		if err == nil {
			var role string
			if strings.Compare(*userMembership.Role, "admin") == 0 {
				role = "owner"
			} else {
				role = "member"
			}

			retVal = append(retVal, GitHubUser{ID: int(*user.ID), Username: *user.Login, Role: role, Email: *user.Email})
		}
	}

	return &retVal, err
}

func GetRepositoryMembers(gitSource *model.GitSource, organizationName string, repositoryRef string) (*[]GitHubUser, error) {
	client := getClient(gitSource)
	users, _, err := client.Repositories.ListCollaborators(context.Background(), organizationName, repositoryRef, nil)

	retVal := make([]GitHubUser, 0)

	for _, user := range users {
		userMembership, _, err := client.Organizations.GetOrgMembership(context.Background(), *user.Login, organizationName)
		if err == nil {
			var role string
			if strings.Compare(*userMembership.Role, "admin") == 0 {
				role = "owner"
			} else {
				role = "member"
			}

			retVal = append(retVal, GitHubUser{ID: int(*user.ID), Username: *user.Login, Role: role, Email: *user.Email})
		}
	}

	return &retVal, err
}

func GetBranches(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) map[string]bool {
	client := getClient(gitSource)
	branchList, _, err := client.Repositories.ListBranches(context.Background(), gitOrgRef, repositoryRef, nil)

	retVal := make(map[string]bool)

	if err != nil {
		for _, branche := range branchList {
			retVal[*branche.Name] = true
		}
	}

	return retVal
}

func CheckRepositoryAgolaConfExists(gitSource *model.GitSource, gitOrgRef string, repositoryRef string) (bool, error) {
	client := getClient(gitSource)
	branchList, _, err := client.Repositories.ListBranches(context.Background(), gitOrgRef, repositoryRef, nil)

	if err != nil {
		return false, err
	}

	for _, branch := range branchList {
		if err != nil {
			return false, err
		}

		tree, _, err := client.Git.GetTree(context.Background(), gitOrgRef, repositoryRef, *branch.Commit.SHA, true)
		if err != nil {
			return false, err
		}

		for _, file := range tree.Entries {
			if strings.Compare(*file.Type, "blob") == 0 && (strings.Compare(*file.Path, ".agola/config.jsonnet") == 0 || strings.Compare(*file.Path, ".agola/config.yml") == 0 || strings.Compare(*file.Path, ".agola/config.json") == 0) {
				return true, nil
			}
		}
	}

	return false, nil
}

func GetCommitMetadata(gitSource *model.GitSource, gitOrgRef string, repositoryRef string, commitSha string) (*dto.CommitMetadataDto, error) {
	client := getClient(gitSource)
	commit, _, err := client.Repositories.GetCommit(context.Background(), gitOrgRef, repositoryRef, commitSha)
	if err != nil {
		return nil, err
	}

	author := make(map[string]string)
	author["email"] = *commit.Commit.Author.Email
	retVal := &dto.CommitMetadataDto{Sha: *commit.SHA, Author: author}

	if commit.Parents != nil {
		retVal.Parents = make([]dto.CommitParentDto, 0)
		for _, parent := range commit.Parents {
			retVal.Parents = append(retVal.Parents, dto.CommitParentDto{Sha: *parent.SHA})
		}
	}

	return retVal, nil
}

func GetOrganization(gitSource *model.GitSource, gitOrgRef string) *dto.OrganizationDto {
	client := getClient(gitSource)
	org, _, err := client.Organizations.Get(context.Background(), gitOrgRef)
	if err != nil {
		return nil
	}

	return &dto.OrganizationDto{Name: *org.Name, ID: int(*org.ID), AvatarURL: *org.AvatarURL}
}

func getClient(gitSource *model.GitSource) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitSource.GitToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
