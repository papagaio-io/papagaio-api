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

func StartOrganizationSync(db repository.Database, tr utils.ConfigUtils, CommonMutex *utils.CommonMutex) {
	go syncOrganizationRun(db, tr, CommonMutex)
}

func syncOrganizationRun(db repository.Database, tr utils.ConfigUtils, CommonMutex *utils.CommonMutex) {
	db.GetOrganizationsTriggerTime()
	for {
		log.Println("start members synk")

		organizationsName, _ := db.GetOrganizationsName()
		for _, organizationName := range organizationsName {
			mutex := utils.ReserveOrganizationMutex(organizationName, CommonMutex)
			mutex.Lock()

			locked := true
			defer utils.ReleaseOrganizationMutexDefer(organizationName, CommonMutex, mutex, &locked)

			org, _ := db.GetOrganizationByName(organizationName)
			if org == nil {
				log.Println("syncOrganizationRun organization ", organizationName, "not found")
				continue
			}

			if agolaOrganizationExists, _ := agola.CheckOrganizationExists(org.Name); !agolaOrganizationExists {
				db.DeleteOrganization(organizationName)
				continue
			}

			log.Println("start synk organization", org.Name)

			gitSource, _ := db.GetGitSourceByName(org.GitSourceName)
			fmt.Println("gitSource:", gitSource)

			membersManager.SynkMembers(org, gitSource)
			repositoryManager.SynkGitRepositorys(db, org, gitSource)

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationName, CommonMutex)
			locked = false
		}

		time.Sleep(time.Duration(tr.GetOrganizationsTriggerTime()) * time.Minute)
	}
}
