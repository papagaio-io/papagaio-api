package utils

import (
	"regexp"

	"wecode.sorint.it/opensource/papagaio-be/model"
)

//TODO
func EvaluateBehaviour(organization *model.Organization, repositoryName string) bool {
	if organization.BehaviourType == model.Regex {

	} else {

	}

	return true
}

//TODO
func ValidateBehaviour(organization *model.Organization) bool {
	if organization.BehaviourType == model.Regex {
		_, err := regexp.Compile(organization.BehaviourInclude)
		if err != nil {
			_, err := regexp.Compile(organization.BehaviourExclude)
			return err == nil
		}

		return true
	} else {

	}

	return true
}
