package agola

import (
	"fmt"
	"net/url"

	"wecode.sorint.it/opensource/papagaio-api/config"
)

const organizationPath string = "%s/api/v1alpha/orgs/%s"
const orgPath string = "%s/api/v1alpha/orgs"
const createMemberPath string = "%s/api/v1alpha/orgs/%s/members/%s"
const getRemoteSourcesPath string = "%s/api/v1alpha/remotesources"
const createProjectPath string = "%s/api/v1alpha/projects"
const deleteProjectPath string = "%s/api/v1alpha/projects/%s"
const organizationMembersPath string = "%s/api/v1alpha/orgs/%s/members"
const runsListPath string = "%s/api/v1alpha/runs?group=%s"

func getOrganizationUrl(agolaOrganizationRef string) string {
	return fmt.Sprintf(organizationPath, config.Config.Agola.AgolaAddr, agolaOrganizationRef)
}

func getOrgUrl() string {
	return fmt.Sprintf(orgPath, config.Config.Agola.AgolaAddr)
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

func getRunsListPath(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) string {
	query := url.QueryEscape("/project/" + projectRef)
	if lastRun {
		query += "&lastrun"
	}
	if len(phase) > 0 {
		query += "&phase=" + phase
	}
	if startRunID != nil {
		query += "&start=" + *startRunID
	}
	if limit > 0 {
		query += "&limit=" + fmt.Sprint(limit)
	}
	if asc {
		query += "&asc"
	}

	return fmt.Sprintf(runsListPath, config.Config.Agola.AgolaAddr, query)
}
