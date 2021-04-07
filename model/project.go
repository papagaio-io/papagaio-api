package model

type Project struct {
	OrganizationID string `json:"organizationID"`
	GitRepoPath    string `json:"gitRepoPath"`
	AgolaProjectID string `json:"agolaProjectID"`
	Archivied      bool   `json:"archivied"`

	//Agola run info. Use branch in key map
	LastBranchRunFailsMap map[string]RunInfo `json:"lastBranchRunMap"`
	LastRun               RunInfo            `json:"lastRun"`
	OlderRunFaild         RunInfo            `json:"olderRunFaild"`

	//Dati per la dashboard
	LastRuns []RunInfo         `json:"lastRuns"`
	Branchs  map[string]Branch `json:"branchs"` //use branch as key
}
