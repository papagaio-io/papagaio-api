package dto

import "time"

type TriggerDto struct {
	IsRunning bool      `json:"isRunning"`
	LastRun   time.Time `json:"lastRun"`
	TimeLeft  uint      `json:"timeLeft"`
}

type TriggersStatusDto struct {
	OrganizationStatus      TriggerDto `json:"organizationStatus"`
	DiscoveryRunFailsStatus TriggerDto `json:"discoveryRunFailsStatus"`
	UserSynkStatus          TriggerDto `json:"userSynkStatus"`
}
