package model

type GitSource struct {
	Name      string `json:"name"`
	GitType   string `json:"gitType"`
	GitAPIURL string `json:"gitApiUrl"`
	GitToken  string `json:"gitToken"`
}
