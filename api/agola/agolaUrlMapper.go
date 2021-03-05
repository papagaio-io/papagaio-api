package agola

import (
	"fmt"

	"wecode.sorint.it/opensource/papagaio-be/config"
)

const createOrgPath string = "%s/api/v1alpha/orgs"
const createMemberPath string = "%s/api/v1alpha/orgs/%s/members/%s"
const getRemoteSourcesPath string = "%s/api/v1alpha/remotesources"

func getCreateORGUrl() string {
	return fmt.Sprintf(createOrgPath, config.Config.Agola.AgolaAddr)
}

func getAddOrgMemberUrl(agolaOrganizationRef string, agolaUserRef string) string {
	return fmt.Sprintf(createMemberPath, config.Config.Agola.AgolaAddr, agolaOrganizationRef, agolaUserRef)
}

func getRemoteSourcesUrl() string {
	return fmt.Sprintf(getRemoteSourcesPath, config.Config.Agola.AgolaAddr)
}
