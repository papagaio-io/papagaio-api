package dto

type OrganizationDto struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	UserEmailCreator string         `json:"userEmailCreator"`
	Visibility       VisibilityType `json:"visibility"`

	BehaviourInclude string        `json:"behaviourInclude"`
	BehaviourExclude string        `json:"behaviourExclude"`
	BehaviourType    BehaviourType `json:"behaviourType"`

	ExternalUsers map[string]bool `json:"externalUsers"`

	Projects    []BranchDto `json:"projects"`
	WorstReport ReportDto   `json:"worstReport,omitempty"`
}
