package dto

type ExternalUsersDto struct {
	ErrorCode OrganizationResponseStatusCode `json:"errorCode"`
	EmailList *[]string                      `json:"emailList,omitempty"`
}
