package model

type Branch struct {
	Name           string    `json:"name"`
	LastSuccessRun RunInfo   `json:"lastSuccessRun"`
	LastFailedRun  RunInfo   `json:"lastFailedRun"`
	LastRuns       []RunInfo `json:"lastRuns"`
}

const lastBranchRunsSize int = 5

func (branch *Branch) PushNewRun(runInfo RunInfo) {
	if runInfo.Result != RunResultFailed && runInfo.Result != RunResultSuccess {
		return
	}

	if runInfo.Result == RunResultFailed {
		if runInfo.RunStartDate.After(branch.LastFailedRun.RunStartDate) {
			branch.LastFailedRun = runInfo
		}
	} else if runInfo.Result == RunResultSuccess {
		if runInfo.RunStartDate.After(branch.LastSuccessRun.RunStartDate) {
			branch.LastSuccessRun = runInfo
		}
	}

	if len(branch.LastRuns) > 0 {
		lastRun := branch.LastRuns[len(branch.LastRuns)-1]
		if !runInfo.RunStartDate.After(lastRun.RunStartDate) {
			return
		}
	}

	branch.LastRuns = append(branch.LastRuns, runInfo)
	if len(branch.LastRuns) > lastBranchRunsSize {
		branch.LastRuns = branch.LastRuns[1:len(branch.LastRuns)]
	}
}
