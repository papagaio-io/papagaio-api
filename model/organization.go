package model

type Organization struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Type     string `json:"type"`
	URL      string `json:"url"` //json:"url,omitempty"
}
