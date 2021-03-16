package model

import "wecode.sorint.it/opensource/papagaio-api/dto"

type Organization struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	UserEmailCreator string             `json:"userEmailCreator"`
	Visibility       dto.VisibilityType `json:"visibility"`

	GitSourceID string `json:"gitSourceId"`
	WebHookID   int    `json:"webHookId"`

	BehaviourInclude string            `json:"behaviourInclude"`
	BehaviourExclude string            `json:"behaviourExclude"`
	BehaviourType    dto.BehaviourType `json:"behaviourType"`

	Projects map[string]Project `json:"projects"`
}
