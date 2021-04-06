package trigger

import (
	"fmt"
	"log"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

func StartOrganizationSync(db repository.Database) {
	go syncOrganizationRun(db)
}

func syncOrganizationRun(db repository.Database) {
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

			membersManager.SynkMembers(&org, gitSource)
			repositoryManager.SynkGitRepositorys(db, &org, gitSource)
		}

		time.Sleep(time.Duration(config.Config.TriggersConfig.OrganizationsDefaultTriggerTime) * time.Minute)
	}
}
