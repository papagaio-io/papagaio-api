package trigger

import (
	"fmt"
	"log"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func StartOrganizationSync(db repository.Database, tr utils.ConfigUtils) {
	go syncOrganizationRun(db, tr)
}

func syncOrganizationRun(db repository.Database, tr utils.ConfigUtils) {
	db.GetOrganizationsTriggerTime()
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

		time.Sleep(time.Duration(tr.GetOrganizationsTriggerTime()) * time.Minute)
	}
}
