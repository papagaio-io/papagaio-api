package dto

import "wecode.sorint.it/opensource/papagaio-be/model"

type CreateOrganizationDto struct {
	Name       string               `json:"name"`
	Visibility model.VisibilityType `json:"visibility"`

	GitSourceId string `json:"gitSourceId"`

	BehaviourInclude string `json:"behaviourInclude"`
	BehaviourExclude string `json:"behaviourExclude"`
	BehaviourType    string `json:"behaviourType"` // wildcard, regex
}
