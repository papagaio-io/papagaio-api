package dto

import "time"

type TriggerRunTimeDto struct {
	Chan           chan TriggerMessage
	TriggerLastRun time.Time
	TimerLastRun   time.Time
	IsRunning      bool
}

type TriggerMessage string

const (
	Restart TriggerMessage = "RESTART"
	Stop    TriggerMessage = "STOP"
)
