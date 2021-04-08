package dto

import "time"

type BranchDto struct {
	Name   string    `json:"name"`
	State  RunState  `json:"state"` //state of last run
	Report ReportDto `json:"report"`

	LastSuccessRunDate time.Time `json:"lastSuccessRunDate,omitempty"`
	LastFailedRunDate  time.Time `json:"lastFailedRunDate,omitempty"`
	LastRunDuration    time.Time `json:"lastRunDuration,omitempty"`
}
