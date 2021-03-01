package model

type Organization struct {
	Name     string `json:"name,omitempty"`
	UserName string `json:"username,omitempty"`
	Type     string `json:"type,omitempty"`
	URL      string `json:"url,omitempty"`
}
