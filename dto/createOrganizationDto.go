package dto

type CreateOrganizationDto struct {
	Name              string `json:"name"`
	AgolaUserRefOwner string `json:"agolaUserRefOwner"`
	Visibility        string `json:"visibility"`
	RemoteSourceName  string `json:"remoteSourceName"`
	GitSourceName     string `json:"gitSourceName"`
	GitOrgRef         string `json:"gitOrgRef"`
}
