package utils

import (
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func GetOrganizationUrl(organization *model.Organization) string {
	return config.Config.Agola.AgolaAddr + "/org/" + organization.AgolaOrganizationRef
}

func GetProjectUrl(organization *model.Organization, project *model.Project) *string {
	url := config.Config.Agola.AgolaAddr + "/org/" + organization.AgolaOrganizationRef + "/projects/" + project.AgolaProjectRef + ".proj"
	return &url
}

func ConvertToAgolaProjectRef(projectName string) string {
	return strings.ReplaceAll(projectName, ".", "")
}

//Return the users map by the agola remoteSource. Key is the git username and value agola userref
func GetUsersMapByRemotesource(agolaApi agola.AgolaApiInterface, agolaRemoteSource string) *map[string]string {
	usersMap := make(map[string]string)

	remotesource, _ := agolaApi.GetRemoteSource(agolaRemoteSource)
	if remotesource == nil {
		return &usersMap
	}
	users, _ := agolaApi.GetUsers()

	if users != nil {
		for _, user := range *users {
			for _, linkedAccount := range user.LinkedAccounts {
				if strings.Compare(remotesource.ID, linkedAccount.RemoteSourceID) == 0 {
					usersMap[linkedAccount.RemoteUserName] = user.Username
					break
				}
			}
		}
	}

	return &usersMap
}

func GetAgolaUserRefByGitUsername(agolaApi agola.AgolaApiInterface, agolaRemoteSource string, gitUsername string) *string {
	users := GetUsersMapByRemotesource(agolaApi, agolaRemoteSource)

	userRef, exists := (*users)[gitUsername]
	if !exists {
		return nil
	}

	return &userRef
}
