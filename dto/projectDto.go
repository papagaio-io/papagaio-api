package dto

type ProjectDto struct {
	Name    string      `json:"name"`
	Branchs []BranchDto `json:"branchs"`

	WorstReport *ReportDto `json:"worstReport,omitempty"`
}
