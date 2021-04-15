package dto

import "time"

type BranchDto struct {
	Name   string     `json:"name"`
	State  RunState   `json:"state"` //state of last run
	Report *ReportDto `json:"report"`

	LastSuccessRunDate *time.Time    `json:"lastSuccessRunDate"`
	LastFailedRunDate  *time.Time    `json:"lastFailedRunDate"`
	LastRunDuration    time.Duration `json:"lastRunDuration"`

	LastSuccessRunURL *string `json:"lastSuccessRunURL"`
	LastFailedRunURL  *string `json:"lastFailedRunURL"`
}
