package model

type GitSource struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	GitType                string `json:"gitType"`
	WebHookOrganizationURL string `json:"iwebHookOrganizationUrl"`
	GitToken               string `json:"gitToken"`
}
