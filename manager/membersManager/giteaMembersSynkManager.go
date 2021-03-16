package membersManager

import (
	"strings"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	gitApi "wecode.sorint.it/opensource/papagaio-api/api/git"
	giteaApi "wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Sincronizzo i membri della organization tra gitea e agola
func SyncMembersForGitea(organization *model.Organization, gitSource *model.GitSource) {
	gitTeams, _ := gitApi.GetOrganizationTeams(gitSource, organization.Name)
	gitTeamOwners := make(map[int]giteaApi.UserTeamResponseDto)
	gitTeamMembers := make(map[int]giteaApi.UserTeamResponseDto)

	for _, team := range *gitTeams {
		teamMembers, _ := gitApi.GetTeamMembers(gitSource, organization.Name, team.ID)

		var teamToCheck *map[int]giteaApi.UserTeamResponseDto
		if strings.Compare(team.Permission, "owner") == 0 {
			teamToCheck = &gitTeamOwners
		} else {
			teamToCheck = &gitTeamMembers
		}

		for _, member := range *teamMembers {
			(*teamToCheck)[member.ID] = member
		}
	}

	agolaMembers, _ := agolaApi.GetOrganizationMembers(organization.Name)
	agolaMembersMap := toMapMembers(&agolaMembers.Members)

	for _, gitMember := range gitTeamOwners {
		agolaUserRef := utils.ConvertGiteaToAgolaUsername(gitMember.Username)
		agolaUserRole := (*agolaMembersMap)[agolaUserRef].Role

		if _, ok := (*agolaMembersMap)[agolaUserRef]; !ok {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "owner")
		} else if agolaUserRole == agolaApi.Owner {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "member")
		}
	}

	for _, gitMember := range gitTeamMembers {
		agolaUserRef := utils.ConvertGiteaToAgolaUsername(gitMember.Username)
		agolaUserRole := (*agolaMembersMap)[agolaUserRef].Role

		if _, ok := (*agolaMembersMap)[agolaUserRef]; !ok {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "member")
		} else if agolaUserRole == agolaApi.Member {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "owner")
		}
	}

	//Verifico i membri eliminati su git

	for _, agolaMember := range agolaMembers.Members {
		if findGiteaMemberByAgolaUserRef(gitTeamOwners, agolaMember.Username) == nil && findGiteaMemberByAgolaUserRef(gitTeamMembers, agolaMember.Username) == nil {
			agolaApi.RemoveOrganizationMember(organization.Name, agolaMember.Username)
		}
	}

}

func findGiteaMemberByAgolaUserRef(gitMembers map[int]giteaApi.UserTeamResponseDto, agolaUserRef string) *giteaApi.UserTeamResponseDto {
	for _, gitMember := range gitMembers {
		if strings.Compare(agolaUserRef, utils.ConvertGiteaToAgolaUsername(gitMember.Username)) == 0 {
			return &gitMember
		}
	}

	return nil
}

func toMapMembers(members *[]agolaApi.MemberDto) *map[string]agolaApi.MemberDto {
	membersMap := make(map[string]agolaApi.MemberDto)
	for _, member := range *members {
		membersMap[member.Username] = member
	}
	return &membersMap
}