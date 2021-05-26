package dto

type ProjectDto struct {
	Name    string      `json:"projectName"`
	Branchs []BranchDto `json:"branchs"`

	WorstReport *ReportDto `json:"worstReport"`
	ProjectUrl  *string    `json:"projectURL"`
}
