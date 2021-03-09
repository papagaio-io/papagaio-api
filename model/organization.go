package model

type Organization struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	UserEmailCreator string `json:"userEmailCreator"`
	//AgolaUserRefOwner string `json:"agolaUserRefOwner"`
	Visibility string `json:"visibility"`
	//RemoteSourceName string `json:"remoteSourceName"`

	GitSourceId string `json:"gitSourceId"`
	//GitOrgRef   string `json:"gitOrgRef"`
	WebHookID int `json:"webHookId"`
}
