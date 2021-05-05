package model

import "wecode.sorint.it/opensource/papagaio-api/types"

type Branch struct {
	Name           string    `json:"name"`
	LastSuccessRun RunInfo   `json:"lastSuccessRun"`
	LastFailedRun  RunInfo   `json:"lastFailedRun"`
	LastRuns       []RunInfo `json:"lastRuns"`
}

const lastBranchRunsSize int = 10

func (branch *Branch) PushNewRun(runInfo RunInfo) {
	if runInfo.Result != types.RunResultFailed && runInfo.Result != types.RunResultSuccess {
		return
	}

	if runInfo.Result == types.RunResultFailed {
		if runInfo.RunStartDate.After(branch.LastFailedRun.RunStartDate) {
			branch.LastFailedRun = runInfo
		}
	} else if runInfo.Result == types.RunResultSuccess {
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
