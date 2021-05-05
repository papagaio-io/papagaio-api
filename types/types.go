package types

import "errors"

type GitType string

const (
	Gitea  GitType = "gitea"
	Github GitType = "github"
)

type VisibilityType string

const (
	Public  VisibilityType = "public"
	Private VisibilityType = "private"
)

func (vt VisibilityType) IsValid() error {
	switch vt {
	case Public, Private:
		return nil
	}
	return errors.New("invalid visibility type")
}

type BehaviourType string

const (
	Wildcard BehaviourType = "wildcard"
	Regex    BehaviourType = "regex"
	None     BehaviourType = "none"
)

func (bt BehaviourType) IsValid() error {
	switch bt {
	case Wildcard, Regex, None:
		return nil
	}
	return errors.New("invalid visibility type")
}

type RunState string

const (
	RunStateSuccess RunState = "success"
	RunStateFailed  RunState = "error"
	//RunStateRunning RunState = "running"
	RunStateNone RunState = "none"
)

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
