package agola

import (
	"fmt"

	"wecode.sorint.it/opensource/papagaio-be/config"
)

const createOrgPath string = "%s/api/v1alpha/orgs"
const createMemberPath string = "%s/api/v1alpha/orgs/%s/members/%s"

func getCreateORGUrl() string {
	return fmt.Sprintf(createOrgPath, config.Config.Agola.AgolaAddr)
}

func getAddOrgMemberUrl(agolaOrganizationRef string, agolaUserRef string) string {
	return fmt.Sprintf(createMemberPath, config.Config.Agola.AgolaAddr, agolaOrganizationRef, agolaUserRef)
}
