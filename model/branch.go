package model

type Branch struct {
	Name           string    `json:"name"`
	LastSuccessRun RunInfo   `json:"lastSuccessRun"`
	LastFailedRun  RunInfo   `json:"lastFailedRun"`
	LastRuns       []RunInfo `json:"lastRuns"`
}
