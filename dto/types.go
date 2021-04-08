package dto

type RunState string

const (
	RunStateSuccess RunState = "success"
	RunStateFailed  RunState = "error"
	RunStateNon     RunState = "none"
)
