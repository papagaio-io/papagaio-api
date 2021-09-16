package dto

import (
	"errors"
	"path/filepath"
	"regexp"

	"wecode.sorint.it/opensource/papagaio-api/types"
)

type CreateOrganizationRequestDto struct {
	GitPath    string               `json:"gitPath"`
	AgolaRef   string               `json:"agolaRef"`
	Visibility types.VisibilityType `json:"visibility"`

	BehaviourInclude string              `json:"behaviourInclude"`
	BehaviourExclude string              `json:"behaviourExclude"`
	BehaviourType    types.BehaviourType `json:"behaviourType"`
}

var organizationRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*([-]?[a-zA-Z0-9]+)+$`)

func (org *CreateOrganizationRequestDto) IsAgolaRefValid() bool {
	return organizationRegexp.MatchString(org.AgolaRef)
}

func (org *CreateOrganizationRequestDto) IsValid() error {
	if org.Visibility.IsValid() == nil && org.BehaviourType.IsValid() == nil && org.IsBehaviourValid() && len(org.GitPath) > 0 && len(org.AgolaRef) > 0 && org.IsAgolaRefValid() {
		return nil
	}
	return errors.New("fields not valid")
}

func (org CreateOrganizationRequestDto) IsBehaviourValid() bool {
	if org.BehaviourType.IsValid() != nil {
		return false
	}

	if org.BehaviourType == types.None {
		return true
	} else if org.BehaviourType == types.Regex {
		_, err := regexp.Compile(org.BehaviourInclude)
		if err != nil {
			if len(org.BehaviourExclude) > 0 {
				_, err := regexp.Compile(org.BehaviourExclude)
				return err == nil
			}
		}

		return true
	} else {
		_, err := filepath.Match(org.BehaviourInclude, "validate")
		if err != nil {
			if len(org.BehaviourExclude) > 0 {
				_, err := filepath.Match(org.BehaviourExclude, "validate")
				return err == nil
			}
		}

		return true
	}
}

type CreateOrganizationResponseDto struct {
	OrganizationURL string                         `json:"organizationURL"`
	ErrorCode       OrganizationResponseStatusCode `json:"errorCode"`
}

type OrganizationResponseStatusCode string

const (
	NoError                         OrganizationResponseStatusCode = "NO_ERROR"
	AgolaOrganizationExistsError    OrganizationResponseStatusCode = "ORG_AGOLA_EXISTS"
	PapagaioOrganizationExistsError OrganizationResponseStatusCode = "ORG_PAPAGAIO_EXISTS"
	GitOrganizationNotFoundError    OrganizationResponseStatusCode = "ORG_GIT_NOT_FOUND"
	AgolaRefNotValid                OrganizationResponseStatusCode = "AGOLA_REF_NOT_VALID"
	UserNotOwnerError               OrganizationResponseStatusCode = "USER_NOT_OWNER"
	UserAgolaRefNotFoundError       OrganizationResponseStatusCode = "USER_AGOLAREF_NOT_FOUND"
)

type DeleteOrganizationResponseDto struct {
	ErrorCode OrganizationResponseStatusCode `json:"errorCode"`
}

type OrganizationResponseDto struct {
	ErrorCode OrganizationResponseStatusCode `json:"errorCode"`
}
