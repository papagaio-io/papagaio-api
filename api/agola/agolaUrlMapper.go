package agola

import (
	"fmt"

	"wecode.sorint.it/opensource/papagaio-be/config"
)

const createOrgPath string = "%s/api/v1alpha/orgs"

func getCreateORGUrl() string {
	return fmt.Sprintf(createOrgPath, config.Config.Agola.AgolaAddr)
}
