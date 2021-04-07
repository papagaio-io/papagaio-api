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

const lastProjectRunsSize int = 4

func (project *Project) PushNewRun(runInfo RunInfo) {
	if len(project.LastRuns) > 0 {
		lastRun := project.LastRuns[len(project.LastRuns)-1]
		if !runInfo.RunStartDate.After(lastRun.RunStartDate) {
			return
		}
	}

	project.LastRuns = append(project.LastRuns, runInfo)
	if len(project.LastRuns) > lastProjectRunsSize {
		project.LastRuns = project.LastRuns[1:len(project.LastRuns)]
	}
}
