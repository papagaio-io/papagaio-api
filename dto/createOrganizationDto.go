package dto

type CreateOrganizationDto struct {
	Name string `json:"name"`
	//AgolaUserRefOwner string `json:"agolaUserRefOwner"`
	Visibility string `json:"visibility"`
	//RemoteSourceName string `json:"remoteSourceName"`
	GitSourceId string `json:"gitSourceId"`
	//GitOrgRef        string `json:"gitOrgRef"`
}
