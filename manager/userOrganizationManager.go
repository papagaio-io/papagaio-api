package manager

import (
	"strings"

	"wecode.sorint.it/opensource/papagaio-be/api/git/gitea"
)

/*func syncUserAccounts(db repository.Database, organizationID string) {
	organization, _ := db.GetOrganizationById(organizationID)
	gitSource, _ := db.GetGitSourceById(organization.ID)

	gitTeams, _ := git.GetOrganizationTeams(gitSource, organization.Name)
	gitTeamOwners := make(map[int]gitea.UserTeamResponseDto)
	gitTeamMembers := make(map[int]gitea.UserTeamResponseDto)

	for _, team := range *gitTeams {
		teamMembers, _ := git.GetTeamMembers(gitSource, team.ID)

		var teamToCheck *map[int]gitea.UserTeamResponseDto
		if strings.Compare(team.Permission, "owner") == 0 {
			teamToCheck = &gitTeamOwners
		} else {
			teamToCheck = &gitTeamMembers
		}

		for _, member := range *teamMembers {

			teamToCheck[member.ID] = member
		}
	}
}*/

func containsMember(teamMembers *[]gitea.UserTeamResponseDto, member gitea.UserTeamResponseDto) bool {
	for _, m := range *teamMembers {
		if strings.Compare(m.Email, member.Email) == 0 {
			return true
		}
	}
	return false
}
