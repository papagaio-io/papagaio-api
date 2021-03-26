package model

import "time"

type Project struct {
	OrganizationID string `json:"organizationID"`
	GitRepoPath    string `json:"gitRepoPath"`
	AgolaProjectID string `json:"agolaProjectID"`
	Archivied      bool   `json:"archivied"`

	//Agola run info. Use branch in key map
	LastBranchRunFailsMap map[string]RunInfo `json:"lastBranchRunMap"`
	LastRun               RunInfo            `json:"lastRun"`
}

type RunInfo struct {
	RunID        string    `json:"runID"`
	RunStartDate time.Time `json:"runStartDate"`
}
