package membersManager

import (
	"strings"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Sincronizzo i membri della organization tra github e agola
func SyncMembersForGithub(organization *model.Organization, gitSource *model.GitSource) {
	githubUsers, _ := github.GetOrganizationMembers(gitSource, organization.Name)

	agolaMembers, _ := agolaApi.GetOrganizationMembers(organization.Name)
	agolaMembersMap := toMapMembers(&agolaMembers.Members)

	for _, gitMember := range *githubUsers {
		agolaUserRef := utils.ConvertGithubToAgolaUsername(gitMember.Username)
		agolaUserRole := (*agolaMembersMap)[agolaUserRef].Role

		if _, ok := (*agolaMembersMap)[agolaUserRef]; !ok {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "owner")
		} else if agolaUserRole == agolaApi.Owner {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "member")
		}
	}

	//Verifico i membri eliminati su git
	for _, agolaMember := range agolaMembers.Members {
		if findGithubMemberByAgolaUserRef(githubUsers, agolaMember.Username) == nil {
			agolaApi.RemoveOrganizationMember(organization.Name, agolaMember.Username)
		}
	}
}

func findGithubMemberByAgolaUserRef(gitMembers *[]github.GitHubUser, agolaUserRef string) *github.GitHubUser {
	for _, gitMember := range *gitMembers {
		if strings.Compare(agolaUserRef, utils.ConvertGithubToAgolaUsername(gitMember.Username)) == 0 {
			return &gitMember
		}
	}

	return nil
}
