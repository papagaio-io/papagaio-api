package model

type Organization struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	UserEmailOwner    string `json:"userEmailOwner"`
	AgolaUserRefOwner string `json:"agolaUserRefOwner"`
	Visibility        string `json:"visibility"`
	RemoteSourceName  string `json:"remoteSourceName"`

	GitSourceName string `json:"gitSourceName"`
	GitOrgRef     string `json:"gitOrgRef"`
	WebHookID     int    `json:"webHookId"`
}
