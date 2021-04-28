package dto

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

type CreateOrganizationRequestDto struct {
	Name       string         `json:"name"`
	AgolaRef   string         `json:"agolaRef"`
	Visibility VisibilityType `json:"visibility"`

	GitSourceName string `json:"gitSourceName"`

	BehaviourInclude string        `json:"behaviourInclude"`
	BehaviourExclude string        `json:"behaviourExclude"`
	BehaviourType    BehaviourType `json:"behaviourType"`
}

func (org CreateOrganizationRequestDto) IsValid() error {
	if org.Visibility.IsValid() == nil && org.BehaviourType.IsValid() == nil && org.IsBehaviourValid() && len(org.Name) > 0 && len(org.GitSourceName) > 0 && len(org.AgolaRef) > 0 && !strings.Contains(org.AgolaRef, ".") {
		return nil
	}
	return errors.New("Invalid visibility type")
}

type BehaviourType string

const (
	Wildcard BehaviourType = "wildcard"
	Regex    BehaviourType = "regex"
	None     BehaviourType = "none"
)

func (bt BehaviourType) IsValid() error {
	switch bt {
	case Wildcard, Regex, None:
		return nil
	}
	return errors.New("Invalid visibility type")
}

type VisibilityType string

const (
	Public  VisibilityType = "public"
	Private VisibilityType = "private"
)

func (vt VisibilityType) IsValid() error {
	switch vt {
	case Public, Private:
		return nil
	}
	return errors.New("Invalid visibility type")
}

func (org CreateOrganizationRequestDto) IsBehaviourValid() bool {
	if org.BehaviourType.IsValid() != nil {
		return false
	}

	if org.BehaviourType == None {
		return true
	} else if org.BehaviourType == Regex {
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
	OrganizationURL string `json:"organizationURL"`
	AgolaExists     bool   `json:"agolaExists"`
}
