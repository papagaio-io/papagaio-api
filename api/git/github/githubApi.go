package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func CreateWebHook(gitSource *model.GitSource, gitOrgRef string) (int, error) {
	client := getClient(gitSource)

	webHookName := "web"
	active := true
	conf := make(map[string]interface{})
	conf["url"] = config.Config.Server.LocalHostAddress + controller.WebHookPath + "/" + gitOrgRef
	fmt.Println("url:", conf["url"])
	conf["content_type"] = "json"
	hook := &github.Hook{Name: &webHookName, Events: []string{"repository"}, Active: &active, Config: conf}
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

func GetOrganizationTeams(gitSource *model.GitSource, gitOrgRef string) (*[]gitea.TeamResponseDto, error) {
	client := getClient(gitSource)
	teams, _, err := client.Teams.ListTeams(context.Background(), gitOrgRef, nil)

	retVal := make([]gitea.TeamResponseDto, 0)
	for _, team := range teams {
		retVal = append(retVal, gitea.TeamResponseDto{ID: int(*team.ID), Name: *team.Name, Permission: *team.Permission})
	}

	return &retVal, err
}

func GetTeamMembers(gitSource *model.GitSource, organizationName string, teamId int) (*[]gitea.UserTeamResponseDto, error) {
	client := getClient(gitSource)
	users, _, err := client.Teams.ListTeamMembers(context.Background(), int64(teamId), nil)

	retVal := make([]gitea.UserTeamResponseDto, 0)
	for _, user := range users {
		if err == nil {
			retVal = append(retVal, gitea.UserTeamResponseDto{ID: int(*user.ID), Username: *user.Name})
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

			retVal = append(retVal, GitHubUser{ID: int(*user.ID), Username: *user.Login, Role: role})
		}
	}

	return &retVal, err
}

func getClient(gitSource *model.GitSource) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitSource.GitToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
