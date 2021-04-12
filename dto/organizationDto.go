package dto

import "time"

type OrganizationDto struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Visibility VisibilityType `json:"visibility"`

	Projects    []ProjectDto `json:"projects"`
	WorstReport *ReportDto   `json:"worstReport,omitempty"`

	LastSuccessRunDate *time.Time    `json:"lastSuccessRunDate,omitempty"`
	LastFailedRunDate  *time.Time    `json:"lastFailedRunDate,omitempty"`
	LastRunDuration    time.Duration `json:"lastRunDuration,omitempty"`

	LastSuccessRunURL *string `json:"lastSuccessRunURL,omitempty"`
	LastFailedRunURL  *string `json:"lastFailedRunURL,omitempty"`
}
