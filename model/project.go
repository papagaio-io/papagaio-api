package model

type Project struct {
	OrganizationID  string `json:"organizationID"`
	GitRepoPath     string `json:"gitRepoPath"`
	AgolaProjectRef string `json:"agolaProjectRef"`
}
