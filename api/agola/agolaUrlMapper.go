package agola

import (
	"fmt"
	"net/url"

	"wecode.sorint.it/opensource/papagaio-api/config"
)

const createOrgPath string = "%s/api/v1alpha/orgs"
const createMemberPath string = "%s/api/v1alpha/orgs/%s/members/%s"
const getRemoteSourcesPath string = "%s/api/v1alpha/remotesources"
const createProjectPath string = "%s/api/v1alpha/projects"
const deleteProjectPath string = "%s/api/v1alpha/projects/%s"
const organizationMembersPath string = "%s/api/v1alpha/orgs/%s/members"

func getCreateORGUrl() string {
	return fmt.Sprintf(createOrgPath, config.Config.Agola.AgolaAddr)
}

func getAddOrgMemberUrl(agolaOrganizationRef string, agolaUserRef string) string {
	return fmt.Sprintf(createMemberPath, config.Config.Agola.AgolaAddr, agolaOrganizationRef, agolaUserRef)
}

func getRemoteSourcesUrl() string {
	return fmt.Sprintf(getRemoteSourcesPath, config.Config.Agola.AgolaAddr)
}

func getCreateProjectUrl() string {
	return fmt.Sprintf(createProjectPath, config.Config.Agola.AgolaAddr)
}

func getDeleteProjectUrl(organizationName string, projectName string) string {
	projectref := url.QueryEscape("org/" + organizationName + "/" + projectName)
	return fmt.Sprintf(deleteProjectPath, config.Config.Agola.AgolaAddr, projectref)
}

func getOrganizationMembersUrl(organizationName string) string {
	return fmt.Sprintf(organizationMembersPath, config.Config.Agola.AgolaAddr, organizationName)
}
