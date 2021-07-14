package dto

import "time"

type TriggersRunTimeDto struct {
	OrganizationsTriggerLastRun      time.Time
	DiscoveryRunFailedTriggerLastRun time.Time
	UsersTriggerLastRun              time.Time
}
