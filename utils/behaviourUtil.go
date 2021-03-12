package utils

import (
	"path/filepath"
	"regexp"

	"wecode.sorint.it/opensource/papagaio-be/dto"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func EvaluateBehaviour(organization *model.Organization, repositoryName string) bool {
	if organization.BehaviourType == dto.None {
		return true
	}
	if organization.BehaviourType == dto.Regex {

		if len(organization.BehaviourExclude) > 0 {
			isMatch := regexp.MustCompile(organization.BehaviourExclude).MatchString(repositoryName)
			if isMatch {
				return false
			}
		}

		return regexp.MustCompile(organization.BehaviourInclude).MatchString(repositoryName)
	} else {
		if len(organization.BehaviourExclude) > 0 {
			isMatch, _ := filepath.Match(organization.BehaviourExclude, repositoryName)
			if isMatch {
				return false
			}
		}
		matched, _ := filepath.Match(organization.BehaviourInclude, repositoryName)
		return matched
	}
}

func ValidateBehaviour(organization *model.Organization) bool {
	if organization.BehaviourType == dto.None {
		return true
	} else if organization.BehaviourType == dto.Regex {
		_, err := regexp.Compile(organization.BehaviourInclude)
		if err != nil {
			if len(organization.BehaviourExclude) > 0 {
				_, err := regexp.Compile(organization.BehaviourExclude)
				return err == nil
			}
		}

		return true
	} else {
		_, err := filepath.Match(organization.BehaviourInclude, "validate")
		if err != nil {
			if len(organization.BehaviourExclude) > 0 {
				_, err := filepath.Match(organization.BehaviourExclude, "validate")
				return err == nil
			}
		}

		return true
	}
}
