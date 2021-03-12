package manager

import (
	"strings"
	"time"

	agolaApi "wecode.sorint.it/opensource/papagaio-be/api/agola"
	gitApi "wecode.sorint.it/opensource/papagaio-be/api/git"
	giteaApi "wecode.sorint.it/opensource/papagaio-be/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-be/model"
	"wecode.sorint.it/opensource/papagaio-be/repository"
	"wecode.sorint.it/opensource/papagaio-be/utils"
)

func StartSyncMembers(db repository.Database) {
	go syncMembersRun(db)
}

func syncMembersRun(db repository.Database) {
	for {

		organizations, _ := db.GetOrganizations()
		for _, org := range *organizations {
			gitSource, _ := db.GetGitSourceById(org.ID)
			SyncMembers(&org, gitSource)
		}

		time.Sleep(time.Hour)
	}
}

//Sincronizzo i membri della organization tra git e agola
func SyncMembers(organization *model.Organization, gitSource *model.GitSource) {
	gitTeams, _ := gitApi.GetOrganizationTeams(gitSource, organization.Name)
	gitTeamOwners := make(map[int]giteaApi.UserTeamResponseDto)
	gitTeamMembers := make(map[int]giteaApi.UserTeamResponseDto)

	for _, team := range *gitTeams {
		teamMembers, _ := gitApi.GetTeamMembers(gitSource, team.ID)

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
		agolaUserRef := utils.ConvertGitToAgolaUsername(gitMember.Username)
		agolaUserRole := (*agolaMembersMap)[agolaUserRef].Role

		if _, ok := (*agolaMembersMap)[agolaUserRef]; !ok {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "owner")
		} else if agolaUserRole == agolaApi.Owner {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "member")
		}
	}

	for _, gitMember := range gitTeamMembers {
		agolaUserRef := utils.ConvertGitToAgolaUsername(gitMember.Username)
		agolaUserRole := (*agolaMembersMap)[agolaUserRef].Role

		if _, ok := (*agolaMembersMap)[agolaUserRef]; !ok {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "member")
		} else if agolaUserRole == agolaApi.Member {
			agolaApi.AddOrUpdateOrganizationMember(organization.Name, agolaUserRef, "owner")
		}
	}

	//Verifico i membri eliminati su git

	for _, agolaMember := range agolaMembers.Members {
		if findGitMemberByAgolaUserRef(gitTeamOwners, agolaMember.Username) == nil || findGitMemberByAgolaUserRef(gitTeamMembers, agolaMember.Username) == nil {
			agolaApi.RemoveOrganizationMember(organization.Name, agolaMember.Username)
		}
	}

}

func findGitMemberByAgolaUserRef(gitMembers map[int]giteaApi.UserTeamResponseDto, agolaUserRef string) *giteaApi.UserTeamResponseDto {
	for _, gitMember := range gitMembers {
		if strings.Compare(agolaUserRef, utils.ConvertGitToAgolaUsername(gitMember.Username)) == 0 {
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
