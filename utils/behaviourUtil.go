package utils

import (
	"strings"

	"wecode.sorint.it/opensource/papagaio-be/model"
)

//TODO
func evaluateBehaviour(organization *model.Organization, repositoryName string) bool {
	if strings.Compare(organization.BehaviourType, "regex") == 0 {

	} else {

	}

	return false
}

//TODO
func validateBehaviour(organization *model.Organization) bool {

	return true
}
