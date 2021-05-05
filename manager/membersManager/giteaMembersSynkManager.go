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
func SyncMembersForGitea(organization *model.Organization, gitSource *model.GitSource, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) {
	log.Println("SyncMembersForGitea start")

	gitTeams, _ := gitGateway.GetOrganizationTeams(gitSource, organization.Name)
	gitTeamOwners := make(map[int]dto.UserTeamResponseDto)
	gitTeamMembers := make(map[int]dto.UserTeamResponseDto)

	for _, team := range *gitTeams {
		teamMembers, _ := gitGateway.GetTeamMembers(gitSource, organization.Name, team.ID)

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

	agolaMembers, _ := agolaApi.GetOrganizationMembers(organization)
	agolaMembersMap := toMapMembers(&agolaMembers.Members)

	agolaUsersMap := utils.GetUsersMapByRemotesource(agolaApi, gitSource.AgolaRemoteSource)

	for _, gitMember := range gitTeamMembers {
		agolaUserRef := (*agolaUsersMap)[gitMember.Username]
		agolaUserRole := (*agolaMembersMap)[agolaUserRef].Role

		if _, ok := (*agolaMembersMap)[agolaUserRef]; !ok {
			agolaApi.AddOrUpdateOrganizationMember(organization, agolaUserRef, "member")
		} else if agolaUserRole == agola.Member {
			agolaApi.AddOrUpdateOrganizationMember(organization, agolaUserRef, "owner")
		}
	}

	for _, gitMember := range gitTeamOwners {
		agolaUserRef := (*agolaUsersMap)[gitMember.Username]
		agolaUserRole := (*agolaMembersMap)[agolaUserRef].Role

		if _, ok := (*agolaMembersMap)[agolaUserRef]; !ok {
			agolaApi.AddOrUpdateOrganizationMember(organization, agolaUserRef, "owner")
		} else if agolaUserRole == agola.Owner {
			agolaApi.AddOrUpdateOrganizationMember(organization, agolaUserRef, "member")
		}
	}

	//Verifico i membri eliminati su git

	for _, agolaMember := range agolaMembers.Members {
		if findGiteaMemberByAgolaUserRef(gitTeamOwners, agolaUsersMap, agolaMember.Username) == nil && findGiteaMemberByAgolaUserRef(gitTeamMembers, agolaUsersMap, agolaMember.Username) == nil {
			agolaApi.RemoveOrganizationMember(organization, agolaMember.Username)
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
		membersMap[member.Username] = member
	}
	return &membersMap
}
