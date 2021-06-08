package model

type Project struct {
	GitRepoPath     string `json:"gitRepoPath"`
	AgolaProjectRef string `json:"agolaProjectRef"`
	AgolaProjectID  string `json:"agolaProjectID"`
	Archivied       bool   `json:"archivied"`

	Branchs map[string]Branch `json:"branchs"` //use branch name as key
}

func (project *Project) GetLastRun() RunInfo {
	var lastRun RunInfo

	if project.Branchs != nil {
		for _, branch := range project.Branchs {
			if branch.LastRuns != nil && len(branch.LastRuns) > 0 {
				branchLastRun := branch.LastRuns[len(branch.LastRuns)-1]
				if branchLastRun.RunStartDate.After(lastRun.RunStartDate) {
					lastRun = branchLastRun
				}
			}
		}
	}

	return lastRun
}

/*func (project *Project) GetLastSuccessRun() *RunInfo {
	var lastSuccessRun *RunInfo = nil

	if project.Branchs != nil {
		for _, branch := range project.Branchs {
			if !branch.LastSuccessRun.RunStartDate.IsZero() && (lastSuccessRun == nil || branch.LastSuccessRun.RunStartDate.After(lastSuccessRun.RunStartDate)) {
				lastSuccessRun = &branch.LastSuccessRun
			}
		}
	}

	return lastSuccessRun
}*/

/*func (project *Project) GetLastFailedRun() *RunInfo {
	var lastFailedRun *RunInfo = nil

	if project.Branchs != nil {
		for _, branch := range project.Branchs {
			if !branch.LastFailedRun.RunStartDate.IsZero() && (lastFailedRun == nil || branch.LastFailedRun.RunStartDate.After(lastFailedRun.RunStartDate)) {
				lastFailedRun = &branch.LastFailedRun
			}
		}
	}

	return lastFailedRun
}*/

func (project *Project) PushNewRun(runInfo RunInfo) {
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
