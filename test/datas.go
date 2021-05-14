package test

import (
	"fmt"

	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

func MakeGitSourceMap() *map[string]model.GitSource {
	retVal := make(map[string]model.GitSource)

	retVal["gitea"] = model.GitSource{Name: "gitea", GitType: types.Gitea, AgolaRemoteSource: "gitea"}
	retVal["github"] = model.GitSource{Name: "github", GitType: types.Github, AgolaRemoteSource: "github"}

	return &retVal
}

func MakeOrganizationMap() *map[string]model.Organization {
	retVal := make(map[string]model.Organization)

	for i := 1; i <= 10; i++ {
		organizationName := "Organization" + fmt.Sprint(i)
		retVal[organizationName] = model.Organization{
			Name:                 organizationName,
			AgolaOrganizationRef: organizationName,
			Visibility:           types.Public,
			GitSourceName:        "gitea",
			UserEmailCreator:     "testuser",
			BehaviourInclude:     "*",
			BehaviourType:        types.Wildcard,
			WebHookID:            i,
		}
	}

	return &retVal
}

func MakeOrganizationList() *[]model.Organization {
	organizationMap := MakeOrganizationMap()
	retVal := make([]model.Organization, 0)

	for _, organization := range *organizationMap {
		retVal = append(retVal, organization)
	}

	return &retVal
}
