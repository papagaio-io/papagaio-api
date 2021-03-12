package model

type GitSource struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	GitType           GitType `json:"gitType"`
	GitAPIURL         string  `json:"gitApiUrl"`
	GitToken          string  `json:"gitToken"`
	AgolaRemoteSource string  `json:"agolaRemoteSource"`
	AgolaToken        string  `json:"agolaToken"`
}

type GitType string

const (
	Gitea  GitType = "gitea"
	Github GitType = "github"
)
