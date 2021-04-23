package utils

import "wecode.sorint.it/opensource/papagaio-api/config"

func GetOrganizationUrl(organizationName string) string {
	return config.Config.Agola.AgolaAddr + "/org/" + organizationName
}
