package dto

import "time"

type ConfigTriggersDto struct {
	OrganizationsTriggerTime uint `json:"organizationsTriggerTime"`
	RunFailedTriggerTime     uint `json:"runFailedTriggerTime"`
	UsersTriggerTime         uint `json:"usersTriggerTime"`

	OrganizationsTriggerLastRun      time.Time `json:"organizationsTriggerLastRun"`
	DiscoveryRunFailedTriggerLastRun time.Time `json:"discoveryRunFailedTriggerLastRun"`
	UsersTriggerLastRun              time.Time `json:"usersTriggerLastRun"`
}
