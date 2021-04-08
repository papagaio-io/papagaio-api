package dto

type ReportDto struct {
	BranchName  string `json:"branchName"`
	ProjectName string `json:"projectName"`

	FailedRuns            uint `json:"failedRuns"`
	TotalRuns             uint `json:"totalRuns"`
	SuccessRunsPercentage uint `json:"successRunsPercentage"`
}
