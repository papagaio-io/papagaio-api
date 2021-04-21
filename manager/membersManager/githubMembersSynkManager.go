package membersManager

import (
	"strings"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Sincronizzo i membri della organization tra github e agola
func SyncMembersForGithub(organization *model.Organization, gitSource *model.GitSource, agolaApi agolaApi.AgolaApiInterface, gitGateway *git.GitGateway) {
	githubUsers, _ := gitGateway.GithubApi.GetOrganizationMembers(gitSource, organization.Name)
	agolaMembers, _ := agolaApi.GetOrganizationMembers(organization.Name)

	for _, gitMember := range *githubUsers {
		agolaUserRef := utils.ConvertGithubToAgolaUsername(gitMember.Username)
		agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, gitMember.Role)
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
