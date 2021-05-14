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

//Synchronize projects and members of organizations
func syncOrganizationRun(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) {
	for {
		log.Println("start members synk")

		organizationsRef, _ := db.GetOrganizationsRef()
		for _, organizationRef := range organizationsRef {
			mutex := utils.ReserveOrganizationMutex(organizationRef, commonMutex)
			mutex.Lock()

			locked := true
			defer utils.ReleaseOrganizationMutexDefer(organizationRef, commonMutex, mutex, &locked)

			org, _ := db.GetOrganizationByAgolaRef(organizationRef)
			if org == nil {
				log.Println("syncOrganizationRun organization ", organizationRef, "not found")
				continue
			}

			//If organization deleted in Agola, delete in Papagaio
			if agolaOrganizationExists, _ := agolaApi.CheckOrganizationExists(org); !agolaOrganizationExists {
				gitSource, _ := db.GetGitSourceByName(org.GitSourceName)
				if gitSource != nil {
					gitGateway.DeleteWebHook(gitSource, org.Name, org.WebHookID)
				}

				db.DeleteOrganization(organizationRef)
				continue
			}

			log.Println("start synk organization", org.Name)

			gitSource, _ := db.GetGitSourceByName(org.GitSourceName)

			membersManager.SynkMembers(org, gitSource, agolaApi, gitGateway)
			repositoryManager.SynkGitRepositorys(db, org, gitSource, agolaApi, gitGateway)

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
			locked = false
		}

		time.Sleep(time.Duration(tr.GetOrganizationsTriggerTime()) * time.Minute)
	}
}
