package manager

import (
	"fmt"
	"log"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

func StartSynkOrganization(db repository.Database, organization *model.Organization, gitSource *model.GitSource) {
	go synkOrganization(db, organization, gitSource)
}

func synkOrganization(db repository.Database, organization *model.Organization, gitSource *model.GitSource) error {
	log.Println("Start organization synk")

	if gitSource.GitType == model.Gitea {
		membersManager.SyncMembersForGitea(organization, gitSource)
	} else {
		membersManager.SyncMembersForGithub(organization, gitSource)
	}

	err := repositoryManager.AddAllGitRepository(db, organization, gitSource)
	if err != nil {
		return err
	}

	return nil
}

func StartSyncMembers(db repository.Database) {
	go syncMembersRun(db)
}

func syncMembersRun(db repository.Database) {
	for {
		log.Println("start members synk")

		organizations, _ := db.GetOrganizations()
		for _, org := range *organizations {
			if !agola.CheckOrganizationExists(org.Name) {
				continue
			}

			log.Println("start synk organization", org.Name)

			gitSource, _ := db.GetGitSourceById(org.GitSourceID)
			fmt.Println("gitSource:", gitSource)

			if gitSource != nil {
				if gitSource.GitType == model.Gitea {
					membersManager.SyncMembersForGitea(&org, gitSource)
				} else {
					membersManager.SyncMembersForGithub(&org, gitSource)
				}
			} else {
				log.Println("Warning!!! Found gitSource null: ", org.ID)
			}
		}

		time.Sleep(10 * time.Minute)
	}
}

func StartRunFailsDiscovery(db repository.Database) {
	go discoveryRunFails(db)
}
