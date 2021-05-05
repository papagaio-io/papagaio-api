package dto

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/types"
)

type CreateOrganizationRequestDto struct {
	Name       string               `json:"name"`
	AgolaRef   string               `json:"agolaRef"`
	Visibility types.VisibilityType `json:"visibility"`

	GitSourceName string `json:"gitSourceName"`

	BehaviourInclude string              `json:"behaviourInclude"`
	BehaviourExclude string              `json:"behaviourExclude"`
	BehaviourType    types.BehaviourType `json:"behaviourType"`
}

var organizationRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*([-]?[a-zA-Z0-9]+)+$`)

func (org *CreateOrganizationRequestDto) IsAgolaRefValid() bool {
	return organizationRegexp.MatchString(org.AgolaRef)
}

func (org *CreateOrganizationRequestDto) IsValid() error {
	if org.Visibility.IsValid() == nil && org.BehaviourType.IsValid() == nil && org.IsBehaviourValid() && len(org.Name) > 0 && len(org.GitSourceName) > 0 && len(org.AgolaRef) > 0 && !strings.Contains(org.AgolaRef, ".") {
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
	OrganizationURL string                               `json:"organizationURL"`
	ErrorCode       CreateOrganizationResponseStatusCode `json:"errorCode"`
}

type CreateOrganizationResponseStatusCode string

const (
	NoError                         CreateOrganizationResponseStatusCode = "NO_ERROR"
	AgolaOrganizationExistsError    CreateOrganizationResponseStatusCode = "ORG_AGOLA_EXISTS"
	PapagaioOrganizationExistsError CreateOrganizationResponseStatusCode = "ORG_PAPAGAIO_EXISTS"
	GitOrganizationNotFoundError    CreateOrganizationResponseStatusCode = "ORG_GIT_NOT_FOUND"
	AgolaRefNotValid                CreateOrganizationResponseStatusCode = "AGOLA_REF_NOT_VALID"
)
