package model

import (
	"wecode.sorint.it/opensource/papagaio-api/types"
)

type Organization struct {
	ID                   string               `json:"id"`
	GitName              string               `json:"gitName" example:"TestDemo"`
	GitPath              string               `json:"gitPath" example:"TestDemo"`
	AgolaOrganizationRef string               `json:"agolaOrganizationRef" example:"TestDemo"`
	UserIDCreator        uint64               `json:"userIdCreator" example:"1"`
	UserIDConnected      uint64               `json:"userIdConnected" example:"1"`
	Visibility           types.VisibilityType `json:"visibility" example:"public"`

	GitSourceName     string `json:"gitSourceName" example:"wecodedev"`
	WebHookID         int    `json:"webHookId"`
	GitOrganizationID int64  `json:"gitOrganizationId"`

	BehaviourInclude string              `json:"behaviourInclude"`
	BehaviourExclude string              `json:"behaviourExclude"`
	BehaviourType    types.BehaviourType `json:"behaviourType" example:"none"`

	Projects      map[string]Project `json:"projects"`
	ExternalUsers map[string]bool    `json:"externalUsers"`
}
