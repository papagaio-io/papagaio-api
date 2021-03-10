package dto

type CreateOrganizationDto struct {
	Name        string `json:"name"`
	Visibility  string `json:"visibility"`
	GitSourceId string `json:"gitSourceId"`
}
