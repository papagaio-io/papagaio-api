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
	OlderRunFaild         RunInfo            `json:"olderRunFaild"`
}

type RunInfo struct {
	ID           string    `json:"id"`
	Branch       string    `json:"branch"`
	RunStartDate time.Time `json:"runStartDate"`
}
