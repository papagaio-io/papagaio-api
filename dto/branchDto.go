package dto

import (
	"time"

	"wecode.sorint.it/opensource/papagaio-api/types"
)

type BranchDto struct {
	Name   string         `json:"name"`
	State  types.RunState `json:"state"` //state of last run
	Report *ReportDto     `json:"report"`

	LastSuccessRunDate *time.Time    `json:"lastSuccessRunDate"`
	LastFailedRunDate  *time.Time    `json:"lastFailedRunDate"`
	LastRunDuration    time.Duration `json:"lastRunDuration" swaggertype:"integer"`

	LastSuccessRunURL string `json:"lastSuccessRunURL"`
	LastFailedRunURL  string `json:"lastFailedRunURL"`
}
