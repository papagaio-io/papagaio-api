package membersManager

import (
	"log"
	"strings"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Sincronizzo i membri della organization tra github e agola
func SyncMembersForGithub(organization *model.Organization, gitSource *model.GitSource, agolaApi agolaApi.AgolaApiInterface, gitGateway *git.GitGateway, user *model.User) {
	githubUsers, _ := gitGateway.GithubApi.GetOrganizationMembers(gitSource, user, organization.GitPath)
	agolaMembers, _ := agolaApi.GetOrganizationMembers(organization)

	agolaUsersMap := utils.GetUsersMapByRemotesource(agolaApi, gitSource.AgolaRemoteSource)

	for _, gitMember := range *githubUsers {
		agolaUserRef, usersExists := (*agolaUsersMap)[gitMember.Username]
		if !usersExists {
			continue
		}

		err := agolaApi.AddOrUpdateOrganizationMember(organization, agolaUserRef, gitMember.Role)
		if err != nil {
			log.Println("AddOrUpdateOrganizationMember error:", err)
		}
	}

	//Verifico i membri eliminati su git
	for _, agolaMember := range agolaMembers.Members {
		if findGithubMemberByAgolaUserRef(githubUsers, agolaUsersMap, agolaMember.User.Username) == nil {
			err := agolaApi.RemoveOrganizationMember(organization, agolaMember.User.Username)
			if err != nil {
				log.Println("RemoveOrganizationMember error:", err)
			}
		}
	}
}

func findGithubMemberByAgolaUserRef(gitMembers *[]github.GitHubUser, agolaUsersMap *map[string]string, agolaUserRef string) *github.GitHubUser {
	for _, gitMember := range *gitMembers {
		if strings.Compare(agolaUserRef, (*agolaUsersMap)[gitMember.Username]) == 0 {
			return &gitMember
		}
	}

	return nil
}
