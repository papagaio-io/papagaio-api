package dto

type ConfigTriggersDto struct {
	OrganizationsDefaultTriggerTime uint `json:"organizationsDefaultTriggerTime"`
	RunFailedDefaultTriggerTime     uint `json:"runFailedDefaultTriggerTime"`
}
