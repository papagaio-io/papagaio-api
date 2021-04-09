package manager

import (
	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func GetOrganizationDto(organization *model.Organization) dto.OrganizationDto {
	retVal := dto.OrganizationDto{
		ID:               organization.ID,
		Name:             organization.Name,
		Visibility:       organization.Visibility,
		BehaviourInclude: organization.BehaviourInclude,
		BehaviourExclude: organization.BehaviourExclude,
		BehaviourType:    organization.BehaviourType,
		ExternalUsers:    organization.ExternalUsers,
	}

	projectList := make([]dto.ProjectDto, 0)
	if organization.Projects != nil {
		for _, project := range organization.Projects {
			projectList = append(projectList, GetProjectDto(&project, organization.Name))
		}
	}
	retVal.Projects = projectList

	var worstReport *dto.ReportDto = nil
	if len(retVal.Projects) > 0 {
		for _, project := range retVal.Projects {
			if project.WorstReport != nil && (worstReport == nil || project.WorstReport.SuccessRunsPercentage < worstReport.SuccessRunsPercentage) {
				worstReport = project.WorstReport
			}
		}
	}
	retVal.WorstReport = worstReport

	return retVal
}

func GetProjectDto(project *model.Project, organizationName string) dto.ProjectDto {
	retVal := dto.ProjectDto{Name: project.GitRepoPath}

	branchList := make([]dto.BranchDto, 0)
	if project.Branchs != nil {
		for _, branch := range project.Branchs {
			branchList = append(branchList, GetBranchDto(branch, project, organizationName))
		}
	}
	retVal.Branchs = branchList

	var worstReport *dto.ReportDto = nil
	if len(retVal.Branchs) > 0 {
		worstReport = retVal.Branchs[0].Report
		for _, branch := range retVal.Branchs {
			if branch.Report.SuccessRunsPercentage < worstReport.SuccessRunsPercentage {
				worstReport = branch.Report
			}
		}
	}
	retVal.WorstReport = worstReport

	return retVal
}

func GetBranchDto(branch model.Branch, project *model.Project, organizationName string) dto.BranchDto {
	retVal := dto.BranchDto{Name: branch.Name}

	if branch.LastRuns == nil || len(branch.LastRuns) == 0 {
		retVal.State = dto.RunStateNone
	} else {
		lastRun := branch.LastRuns[len(branch.LastRuns)-1]
		if lastRun.Result == model.RunResultSuccess {
			retVal.State = dto.RunStateSuccess
		} else {
			retVal.State = dto.RunStateFailed
		}
	}

	retVal.Report = GetBranchReport(branch, project.GitRepoPath, organizationName)

	lastRun := project.GetLastRun()
	if !lastRun.RunStartDate.IsZero() {
		if lastRun.Result == model.RunResult(agola.RunResultSuccess) {
			retVal.State = dto.RunStateSuccess
		} else {
			retVal.State = dto.RunStateFailed
		}

		retVal.LastRunDuration = lastRun.RunEndDate.Sub(lastRun.RunStartDate)
	} else {
		retVal.State = dto.RunStateNone
	}

	lastSuccessRun := project.GetLastSuccessRun()
	if lastSuccessRun != nil {
		retVal.LastSuccessRunDate = lastSuccessRun.RunStartDate
	}

	lastFailedRun := project.GetLastFailedRun()
	if lastFailedRun != nil {
		retVal.LastFailedRunDate = lastFailedRun.RunStartDate
	}

	return retVal
}

func GetBranchReport(branch model.Branch, projectName string, organizationName string) *dto.ReportDto {
	report := dto.ReportDto{BranchName: branch.Name, ProjectName: projectName, OrganizationName: organizationName}

	failedRuns := uint(0)
	for _, run := range branch.LastRuns {
		if run.Result == model.RunResultFailed {
			failedRuns++
		}
	}

	report.FailedRuns = failedRuns
	report.TotalRuns = uint(len(branch.LastRuns))
	report.SuccessRunsPercentage = 100 - ((report.TotalRuns * 100) / report.TotalRuns)

	return &report
}
