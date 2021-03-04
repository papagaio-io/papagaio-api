package agola

import (
	"fmt"
	"wecode.sorint.it/opensource/papagaio-be/config"
)

const createTokenPath string = "%s/api/v1alpha/users/%s/tokens"

func getCreateTokenUrl(agolaUserRef string) string {
	return fmt.Sprintf(createTokenPath, config.Config.Agola.AgolaAddr, agolaUserRef)
}
