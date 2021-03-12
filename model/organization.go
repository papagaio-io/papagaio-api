package model

type Organization struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	UserEmailCreator string `json:"userEmailCreator"`
	Visibility       string `json:"visibility"`

	GitSourceID string `json:"gitSourceId"`
	WebHookID   int    `json:"webHookId"`

	BehaviourInclude string        `json:"behaviourInclude"`
	BehaviourExclude string        `json:"behaviourExclude"`
	BehaviourType    BehaviourType `json:"behaviourType"` // wildcard, regex

	Projects map[string]Project `json:"projects"`
}

type BehaviourType string

const (
	Wildcard BehaviourType = "wildcard"
	Regex    BehaviourType = "regex"
)
