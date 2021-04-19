package dto

import "time"

type OrganizationDto struct {
	ID         string         `json:"id"`
	Name       string         `json:"organizationName"`
	Visibility VisibilityType `json:"visibility"`
	AvatarURL  string         `json:"avatarUrl"`

	Projects    []ProjectDto `json:"projects"`
	WorstReport *ReportDto   `json:"worstReport"`

	LastSuccessRunDate *time.Time    `json:"lastSuccessRunDate"`
	LastFailedRunDate  *time.Time    `json:"lastFailedRunDate"`
	LastRunDuration    time.Duration `json:"lastRunDuration"`

	LastSuccessRunURL *string `json:"lastSuccessRunURL"`
	LastFailedRunURL  *string `json:"lastFailedRunURL"`
}
