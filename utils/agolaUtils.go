package utils

import (
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func GetOrganizationUrl(organization *model.Organization) string {
	return config.Config.Agola.AgolaAddr + "/org/" + organization.AgolaOrganizationRef
}

func GetProjectUrl(organization *model.Organization, projectName string) string {
	return config.Config.Agola.AgolaAddr + "/org/" + organization.AgolaOrganizationRef + "/projects/" + projectName + ".proj"
}

func ConvertToAgolaOrganizationRef(organizationName string) string {
	return strings.ReplaceAll(organizationName, ".", "")
}
