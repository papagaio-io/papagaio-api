package model

type Organization struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	UserEmailCreator string         `json:"userEmailCreator"`
	Visibility       VisibilityType `json:"visibility"`

	GitSourceID string `json:"gitSourceId"`
	WebHookID   int    `json:"webHookId"`

	BehaviourInclude string        `json:"behaviourInclude"`
	BehaviourExclude string        `json:"behaviourExclude"`
	BehaviourType    BehaviourType `json:"behaviourType"`

	Projects map[string]Project `json:"projects"`
}

type BehaviourType string

const (
	Wildcard BehaviourType = "wildcard"
	Regex    BehaviourType = "regex"
	None     BehaviourType = "none"
)

type VisibilityType string

const (
	Public  VisibilityType = "public"
	Private VisibilityType = "private"
)
