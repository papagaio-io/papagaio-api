package dto

import "time"

type TriggerRunTimeDto struct {
	Chan           chan TriggerStarter
	Starter        TriggerStarter
	TriggerLastRun time.Time
	TimerLastRun   time.Time
	IsRunning      bool
}

type TriggerStarter string

const (
	Trigger TriggerStarter = "TRIGGER"
	Service TriggerStarter = "SERVICE"
)
