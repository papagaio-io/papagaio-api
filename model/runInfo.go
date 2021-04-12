package model

import (
	"fmt"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/config"
)

type RunInfo struct {
	ID           string    `json:"id"`
	Branch       string    `json:"branch"`
	RunStartDate time.Time `json:"runStartDate"`
	RunEndDate   time.Time `json:"runEndDate,omitempty"`
	Phase        RunPhase  `json:"phase"`
	Result       RunResult `json:"result"`
}

type RunPhase string

const (
	RunPhaseSetupError RunPhase = "setuperror"
	RunPhaseQueued     RunPhase = "queued"
	RunPhaseCancelled  RunPhase = "cancelled"
	RunPhaseRunning    RunPhase = "running"
	RunPhaseFinished   RunPhase = "finished"
)

type RunResult string

const (
	RunResultUnknown RunResult = "unknown"
	RunResultStopped RunResult = "stopped"
	RunResultSuccess RunResult = "success"
	RunResultFailed  RunResult = "failed"
)

const runURL string = "%s/org/%s/projects/%s.proj/runs/%s"

func (run *RunInfo) GetURL(organizationName string, projectName string) string {
	return fmt.Sprintf(runURL, config.Config.Agola.AgolaAddr, organizationName, projectName, run.ID)
}
