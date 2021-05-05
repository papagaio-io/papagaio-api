package model

import (
	"wecode.sorint.it/opensource/papagaio-api/types"
)

type Organization struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	AgolaOrganizationRef string               `json:"agolaOrganizationRef"`
	UserEmailCreator     string               `json:"userEmailCreator"`
	Visibility           types.VisibilityType `json:"visibility"`

	GitSourceName string `json:"gitSourceName"`
	WebHookID     int    `json:"webHookId"`

	BehaviourInclude string              `json:"behaviourInclude"`
	BehaviourExclude string              `json:"behaviourExclude"`
	BehaviourType    types.BehaviourType `json:"behaviourType"`

	Projects      map[string]Project `json:"projects"`
	ExternalUsers map[string]bool    `json:"externalUsers"`
}
