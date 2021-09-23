package dto

type ConfigTriggersDto struct {
	OrganizationsTriggerTime uint `json:"organizationsTriggerTime"`
	RunFailedTriggerTime     uint `json:"runFailedTriggerTime"`
	UsersTriggerTime         uint `json:"usersTriggerTime"`
}
