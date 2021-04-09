package dto

type RunState string

const (
	RunStateSuccess RunState = "success"
	RunStateFailed  RunState = "error"
	//RunStateRunning RunState = "running"
	RunStateNone RunState = "none"
)
