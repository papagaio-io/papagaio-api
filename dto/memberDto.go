package dto

type memberDto struct {
	OrganizationID string `json:"organizationId"`
	UserRef        string `json:"userRef"`
	UserType       string `json:"userType"`
}
