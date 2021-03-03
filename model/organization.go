package model

type Organization struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	UserRefOwner string `json:"userRefOwner"`
	Visibility   string `json:"visibility"`
	//Token string `json:"token"`
	RemoteSourceName string `json:"remoteSourceName"`

	GitSourceID int    `json:"gitSourceId"`
	GirOrgRef   string `json:"gitOrgRef"`
	WebHookID   int    `json:"webHookId"`
}
