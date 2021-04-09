package dto

type OrganizationDto struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Visibility VisibilityType `json:"visibility"`

	Projects    []ProjectDto `json:"projects"`
	WorstReport *ReportDto   `json:"worstReport,omitempty"`
}
