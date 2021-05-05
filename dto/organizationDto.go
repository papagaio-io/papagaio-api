package dto

import (
	"time"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

type OrganizationDto struct {
	ID         string               `json:"id"`
	Name       string               `json:"organizationName"`
	AgolaRef   string               `json:"agolaRef"`
	Visibility types.VisibilityType `json:"visibility"`
	AvatarURL  string               `json:"avatarUrl"`

	Projects    []ProjectDto `json:"projects"`
	WorstReport *ReportDto   `json:"worstReport"`

	LastSuccessRunDate *time.Time    `json:"lastSuccessRunDate"`
	LastFailedRunDate  *time.Time    `json:"lastFailedRunDate"`
	LastRunDuration    time.Duration `json:"lastRunDuration"`

	LastSuccessRunURL string `json:"lastSuccessRunURL"`
	LastFailedRunURL  string `json:"lastFailedRunURL"`
	OrganizationURL   string `json:"organizationURL"`
}
