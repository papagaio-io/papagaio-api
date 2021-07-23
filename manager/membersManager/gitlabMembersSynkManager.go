package membersManager

import (
	"log"
	"strings"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitlab"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Sincronizzo i membri della organization tra github e agola
func SyncMembersForGitlab(organization *model.Organization, gitSource *model.GitSource, agolaApi agolaApi.AgolaApiInterface, gitGateway *git.GitGateway, user *model.User) {
	gitlabUsers, _ := gitGateway.GitlabApi.GetOrganizationMembers(gitSource, user, organization.GitPath)
	agolaMembers, _ := agolaApi.GetOrganizationMembers(organization)

	agolaUsersMap := utils.GetUsersMapByRemotesource(agolaApi, gitSource.AgolaRemoteSource)

	for _, gitMember := range *gitlabUsers {
		agolaUserRef, usersExists := (*agolaUsersMap)[gitMember.Username]
		if !usersExists {
			continue
		}

		var role string
		if gitMember.HasOwnerPermission() {
			role = "owner"
		} else {
			role = "member"
		}
		err := agolaApi.AddOrUpdateOrganizationMember(organization, agolaUserRef, role)
		if err != nil {
			log.Println("AddOrUpdateOrganizationMember error:", err)
		}
	}

	//Verifico i membri eliminati su git
	for _, agolaMember := range agolaMembers.Members {
		if findGitlabMemberByAgolaUserRef(gitlabUsers, agolaUsersMap, agolaMember.User.Username) == nil {
			err := agolaApi.RemoveOrganizationMember(organization, agolaMember.User.Username)
			if err != nil {
				log.Println("RemoveOrganizationMember error:", err)
			}
		}
	}
}

func findGitlabMemberByAgolaUserRef(gitMembers *[]gitlab.GitlabUser, agolaUsersMap *map[string]string, agolaUserRef string) *gitlab.GitlabUser {
	for _, gitMember := range *gitMembers {
		if strings.Compare(agolaUserRef, (*agolaUsersMap)[gitMember.Username]) == 0 {
			return &gitMember
		}
	}

	return nil
}
