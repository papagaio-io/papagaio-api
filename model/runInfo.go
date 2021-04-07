package model

import "time"

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
