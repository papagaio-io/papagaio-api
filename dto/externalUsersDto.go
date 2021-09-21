package dto

import "regexp"

type ExternalUsersDto struct {
	ErrorCode OrganizationResponseStatusCode `json:"errorCode"`
	EmailList *[]string                      `json:"emailList,omitempty"`
}

var emailRegex = regexp.MustCompile(`[a-z0-9!#$%&'*+/=?^_"{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_"{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?`)

func (org ExternalUserDto) IsEmailValid() bool {
	return len(org.Email) < 101 && emailRegex.MatchString(org.Email)
}
