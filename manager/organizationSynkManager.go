package manager

import (
	"time"

	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

func StartSynkOrganization(db repository.Database, organization *model.Organization, gitSource *model.GitSource) {
	go synkOrganization(db, organization, gitSource)
}

func synkOrganization(db repository.Database, organization *model.Organization, gitSource *model.GitSource) {
	repositoryManager.AddAllGitRepository(db, organization, gitSource)
	if gitSource.GitType == model.Gitea {
		membersManager.SyncMembersForGitea(organization, gitSource)
	} else {
		membersManager.SyncMembersForGithub(organization, gitSource)
	}
}

func StartSyncMembers(db repository.Database) {
	go syncMembersRun(db)
}

func syncMembersRun(db repository.Database) {
	for {

		organizations, _ := db.GetOrganizations()
		for _, org := range *organizations {
			gitSource, _ := db.GetGitSourceById(org.ID)
			if gitSource.GitType == model.Gitea {
				membersManager.SyncMembersForGitea(&org, gitSource)
			} else {
				membersManager.SyncMembersForGithub(&org, gitSource)
			}
		}

		time.Sleep(30 * time.Minute)
	}
}
