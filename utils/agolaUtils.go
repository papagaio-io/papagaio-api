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

func GetProjectUrl(organization *model.Organization, project *model.Project) string {
	return config.Config.Agola.AgolaAddr + "/org/" + organization.AgolaOrganizationRef + "/projects/" + project.AgolaProjectRef + ".proj"
}

func ConvertToAgolaProjectRef(projectName string) string {
	return strings.ReplaceAll(projectName, ".", "")
}

//Return the users map by the agola remoteSource. Key is the git username and value agola userref
func GetUsersMapByRemotesource(agolaApi agola.AgolaApiInterface, agolaRemoteSource string) *map[string]string {
	gitSource, _ := agolaApi.GetRemoteSource(agolaRemoteSource)
	users, _ := agolaApi.GetUsers()

	usersMap := make(map[string]string)

	for _, user := range *users {
		for _, linkedAccount := range user.LinkedAccounts {
			if strings.Compare(gitSource.ID, linkedAccount.RemoteSourceID) == 0 {
				usersMap[linkedAccount.RemoteUserName] = user.Username
				break
			}
		}
	}

	return &usersMap
}
