package trigger

import (
	"log"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func StartOrganizationSync(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) {
	go syncOrganizationRun(db, tr, commonMutex, agolaApi, gitGateway)
}

func syncOrganizationRun(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) {
	db.GetOrganizationsTriggerTime()
	for {
		log.Println("start members synk")

		organizationsName, _ := db.GetOrganizationsName()
		for _, organizationName := range organizationsName {
			mutex := utils.ReserveOrganizationMutex(organizationName, commonMutex)
			mutex.Lock()

			locked := true
			defer utils.ReleaseOrganizationMutexDefer(organizationName, commonMutex, mutex, &locked)

			org, _ := db.GetOrganizationByName(organizationName)
			if org == nil {
				log.Println("syncOrganizationRun organization ", organizationName, "not found")
				continue
			}

			if !agolaApi.CheckOrganizationExists(org.Name) {
				continue
			}

			log.Println("start synk organization", org.Name)

			gitSource, _ := db.GetGitSourceByName(org.GitSourceName)

			membersManager.SynkMembers(org, gitSource, agolaApi, gitGateway)
			repositoryManager.SynkGitRepositorys(db, org, gitSource, agolaApi, gitGateway)

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationName, commonMutex)
			locked = false
		}

		time.Sleep(time.Duration(tr.GetOrganizationsTriggerTime()) * time.Minute)
	}
}
