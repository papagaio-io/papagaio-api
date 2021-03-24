package model

import "time"

type Project struct {
	OrganizationID string `json:"organizationID"`
	GitRepoPath    string `json:"gitRepoPath"`
	AgolaProjectID string `json:"agolaProjectID"`
	Archivied      bool   `json:"archivied"`

	//Agola run info. Use branch in key map
	LastBranchRunMap map[string]RunInfo `json:"lastBranchRunMap"`
}

type RunInfo struct {
	LastRunID        string    `json:"lastRunID"`
	LastRunStartDate time.Time `json:"lastRunStartDate"`
	ISLastRunFailed  bool      `json:"sSLastRunFailed"`
}
