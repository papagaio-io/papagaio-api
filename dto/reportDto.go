package dto

type ReportDto struct {
	BranchName       string `json:"branchName"`
	ProjectName      string `json:"projectName"`
	OrganizationName string `json:"organizationName"`

	FailedRuns            uint `json:"failedRuns"`
	TotalRuns             uint `json:"totalRuns"`
	SuccessRunsPercentage uint `json:"successRunsPercentage"`
}
