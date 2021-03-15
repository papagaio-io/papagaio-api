package manager

import (
	"wecode.sorint.it/opensource/papagaio-be/model"
	"wecode.sorint.it/opensource/papagaio-be/repository"
)

func StartSynkOrganization(db repository.Database, organization *model.Organization, gitSource *model.GitSource) {
	go synkOrganization(db, organization, gitSource)
}

func synkOrganization(db repository.Database, organization *model.Organization, gitSource *model.GitSource) {
	AddAllGitRepository(db, organization, gitSource)
	if gitSource.GitType == model.Gitea {
		SyncMembersForGitea(organization, gitSource)
	} else {
		SyncMembersForGithub(organization, gitSource)
	}
}
