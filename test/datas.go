package test

import "wecode.sorint.it/opensource/papagaio-api/model"

func MakeGitSourceList() *map[string]model.GitSource {
	retVal := make(map[string]model.GitSource)

	retVal["gitea"] = model.GitSource{Name: "gitea", GitType: model.Gitea}
	retVal["github"] = model.GitSource{Name: "github", GitType: model.Github}

	return &retVal
}
