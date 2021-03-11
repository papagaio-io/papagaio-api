package dto

type CreateOrganizationDto struct {
	Name       string `json:"name"`
	Visibility string `json:"visibility"`

	GitSourceId string `json:"gitSourceId"`

	BehaviourInclude string `json:"behaviourInclude"`
	BehaviourExclude string `json:"behaviourExclude"`
	BehaviourType    string `json:"behaviourType"` // wildcard, regex
}
