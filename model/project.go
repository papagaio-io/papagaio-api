package model

type Project struct {
	OrganizationID string `json:"organizationID"`
	GitRepoPath    string `json:"gitRepoPath"`
	AgolaProjectID string `json:"agolaProjectID"`
	Archivied      bool   `json:"archivied"`

	LastRuns []RunInfo         `json:"lastRuns"`
	Branchs  map[string]Branch `json:"branchs"` //use branch as key
}

const lastProjectRunsSize int = 4

func (project *Project) GetLastRun() RunInfo {
	if len(project.LastRuns) == 0 {
		return RunInfo{}
	}

	return project.LastRuns[len(project.LastRuns)-1]
}

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

	//push into branch list

	if project.Branchs == nil {
		project.Branchs = make(map[string]Branch)
	}

	if _, ok := project.Branchs[runInfo.Branch]; !ok {
		project.Branchs[runInfo.Branch] = Branch{Name: runInfo.Branch, LastRuns: make([]RunInfo, 0)}
	}
	branch := project.Branchs[runInfo.Branch]
	branch.PushNewRun(runInfo)
	project.Branchs[runInfo.Branch] = branch
}
