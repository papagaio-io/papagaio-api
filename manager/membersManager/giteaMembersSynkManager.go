package membersManager

import (
	"log"
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Sincronizzo i membri della organization tra gitea e agola
func SyncMembersForGitea(organization *model.Organization, gitSource *model.GitSource, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway, user *model.User) {
	log.Println("SyncMembersForGitea start")

	gitTeams, err := gitGateway.GetOrganizationTeams(gitSource, user, organization.Name)
	if err != nil {
		log.Println("error in GetOrganizationTeams:", err)
		return
	}
	gitTeamOwners := make(map[int]dto.UserTeamResponseDto)
	gitTeamMembers := make(map[int]dto.UserTeamResponseDto)

	for _, team := range *gitTeams {
		teamMembers, _ := gitGateway.GetTeamMembers(gitSource, user, organization.GitOrganizationID, team.ID)

		var teamToCheck *map[int]dto.UserTeamResponseDto
		if strings.Compare(team.Permission, "owner") == 0 {
			teamToCheck = &gitTeamOwners
		} else {
			teamToCheck = &gitTeamMembers
		}

		for _, member := range *teamMembers {
			(*teamToCheck)[member.ID] = member
		}
	}

	agolaOrganizationMembers, _ := agolaApi.GetOrganizationMembers(organization)
	agolaOrganizationMembersMap := toMapMembers(&agolaOrganizationMembers.Members)

	agolaUsersMap := utils.GetUsersMapByRemotesource(agolaApi, gitSource.AgolaRemoteSource)

	for _, gitMember := range gitTeamMembers {
		agolaUserRef, usersExists := (*agolaUsersMap)[gitMember.Username]
		if !usersExists {
			continue
		}

		if agolaMember, ok := (*agolaOrganizationMembersMap)[agolaUserRef]; !ok || agolaMember.Role == agola.Owner {
			agolaApi.AddOrUpdateOrganizationMember(organization, agolaUserRef, "member")
		}
	}

	for _, gitMember := range gitTeamOwners {
		agolaUserRef, usersExists := (*agolaUsersMap)[gitMember.Username]
		if !usersExists {
			continue
		}

		if agolaMember, ok := (*agolaOrganizationMembersMap)[agolaUserRef]; !ok || agolaMember.Role == agola.Member {
			agolaApi.AddOrUpdateOrganizationMember(organization, agolaUserRef, "owner")
		}
	}

	//Verifico i membri eliminati su git

	for _, agolaMember := range agolaOrganizationMembers.Members {
		log.Println("Check user in git:", agolaMember)

		if findGiteaMemberByAgolaUserRef(gitTeamOwners, agolaUsersMap, agolaMember.User.Username) == nil && findGiteaMemberByAgolaUserRef(gitTeamMembers, agolaUsersMap, agolaMember.User.Username) == nil {
			agolaApi.RemoveOrganizationMember(organization, agolaMember.User.Username)
		}
	}

	log.Println("SyncMembersForGitea end")
}

func findGiteaMemberByAgolaUserRef(gitMembers map[int]dto.UserTeamResponseDto, agolaUsersMap *map[string]string, agolaUserRef string) *dto.UserTeamResponseDto {
	for _, gitMember := range gitMembers {
		if strings.Compare(agolaUserRef, (*agolaUsersMap)[gitMember.Username]) == 0 {
			return &gitMember
		}

	}

	return nil
}

func toMapMembers(members *[]agola.MemberDto) *map[string]agola.MemberDto {
	membersMap := make(map[string]agola.MemberDto)
	for _, member := range *members {
		membersMap[member.User.Username] = member
	}
	return &membersMap
}
