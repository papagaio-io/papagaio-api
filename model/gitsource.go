package model

type GitSource struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GitType   string `json:"gitType"`
	GitAPIURL string `json:"gitApiUrl"`
	GitToken  string `json:"gitToken"`
}
