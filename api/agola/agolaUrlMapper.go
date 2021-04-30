package agola

import (
	"fmt"
	"net/url"

	"wecode.sorint.it/opensource/papagaio-api/config"
)

const organizationPath string = "%s/api/v1alpha/orgs/%s"
const orgPath string = "%s/api/v1alpha/orgs"
const createMemberPath string = "%s/api/v1alpha/orgs/%s/members/%s"
const createProjectPath string = "%s/api/v1alpha/projects"
const projectPath string = "%s/api/v1alpha/projects/%s"
const organizationMembersPath string = "%s/api/v1alpha/orgs/%s/members"
const runsListPath string = "%s/api/v1alpha/runs?group=%s"
const runPath string = "%s/api/v1alpha/runs/%s"
const taskPath string = "%s/api/v1alpha/runs/%s/tasks/%s"
const logsPath string = "%s/api/v1alpha/logs?runID=%s&taskID=%s&%s"

func getOrganizationUrl(agolaOrganizationRef string) string {
	return fmt.Sprintf(organizationPath, config.Config.Agola.AgolaAddr, agolaOrganizationRef)
}

func getOrgUrl() string {
	return fmt.Sprintf(orgPath, config.Config.Agola.AgolaAddr)
}

func getAddOrgMemberUrl(agolaOrganizationRef string, agolaUserRef string) string {
	return fmt.Sprintf(createMemberPath, config.Config.Agola.AgolaAddr, agolaOrganizationRef, agolaUserRef)
}

func getCreateProjectUrl() string {
	return fmt.Sprintf(createProjectPath, config.Config.Agola.AgolaAddr)
}

func getProjectUrl(organizationName string, projectName string) string {
	projectref := url.QueryEscape("org/" + organizationName + "/" + projectName)
	return fmt.Sprintf(projectPath, config.Config.Agola.AgolaAddr, projectref)
}

func getOrganizationMembersUrl(organizationName string) string {
	return fmt.Sprintf(organizationMembersPath, config.Config.Agola.AgolaAddr, organizationName)
}

func getRunsListUrl(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) string {
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

func getRunUrl(runID string) string {
	return fmt.Sprintf(runPath, config.Config.Agola.AgolaAddr, runID)
}

func getTaskUrl(runID string, taskID string) string {
	return fmt.Sprintf(taskPath, config.Config.Agola.AgolaAddr, runID, taskID)
}

func getLogsUrl(runID string, taskID string, step int) string {
	stepParam := "setup"
	if step != -1 {
		stepParam = "step=" + fmt.Sprint(step)
	}

	return fmt.Sprintf(logsPath, config.Config.Agola.AgolaAddr, runID, taskID, stepParam)
}
