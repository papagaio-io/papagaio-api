package test

import "wecode.sorint.it/opensource/papagaio-api/model"

func MakeGitSourceMap() *map[string]model.GitSource {
	retVal := make(map[string]model.GitSource)

	retVal["gitea"] = model.GitSource{Name: "gitea", GitType: model.Gitea}
	retVal["github"] = model.GitSource{Name: "github", GitType: model.Github}

	return &retVal
}

func MakeOrganizationMap() *map[string]model.Organization {
	retVal := make(map[string]model.Organization)

	//TODO
	//retVal["Organization1"] = model.

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
