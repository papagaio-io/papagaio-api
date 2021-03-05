package model

type Organization struct {
	ID                string `json:"id"`
	Name              string `json:"name"` //TODO remove?
	UserEmail         string `json:"userEmail"`
	AgolaUserRefOwner string `json:"agolaUserRefOwner"`
	Visibility        string `json:"visibility"`
	//Token string `json:"token"`
	RemoteSourceName string `json:"remoteSourceName"`

	GitSourceID int    `json:"gitSourceId"`
	GitOrgRef   string `json:"gitOrgRef"`
	WebHookID   int    `json:"webHookId"`
}
