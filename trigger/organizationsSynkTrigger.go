package trigger

import (
	"fmt"
	"log"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/trigger/dto"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func StartOrganizationSync(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway, c chan string, rtDto *dto.TriggersRunTimeDto) {
	go syncOrganizationRun(db, tr, commonMutex, agolaApi, gitGateway, c, rtDto)
	go syncOrganizationRunTimer(tr, c)
}

func syncOrganizationRunTimer(tr utils.ConfigUtils, c chan string) {
	for {
		log.Println("syncOrganizationRunTimer wait for", time.Duration(time.Minute.Nanoseconds()*int64(tr.GetOrganizationsTriggerTime())))
		time.Sleep(time.Duration(time.Minute.Nanoseconds() * int64(tr.GetOrganizationsTriggerTime())))
		c <- "resume from timer"
	}
}

//Synchronize projects and members of organizations
func syncOrganizationRun(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway, c chan string, rtDto *dto.TriggersRunTimeDto) {
	for {
		log.Println("start syncOrganizationRun")
		rtDto.OrganizationsTriggerLastRun = time.Now()

		organizationsRef, _ := db.GetOrganizationsRef()
		for _, organizationRef := range organizationsRef {
			log.Println("syncOrganizationRun organizationRef:", organizationRef)
			mutex := utils.ReserveOrganizationMutex(organizationRef, commonMutex)
			mutex.Lock()

			locked := true
			defer utils.ReleaseOrganizationMutexDefer(organizationRef, commonMutex, mutex, &locked)

			org, _ := db.GetOrganizationByAgolaRef(organizationRef)
			if org == nil {
				log.Println("syncOrganizationRun organization ", organizationRef, "not found")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
				locked = false

				continue
			}

			gitSource, _ := db.GetGitSourceByName(org.GitSourceName)
			if gitSource == nil {
				log.Println("gitSource", org.GitSourceName, "not found")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
				locked = false

				continue
			}

			user, _ := db.GetUserByUserId(org.UserIDConnected)
			if user == nil {
				log.Println("user not found")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
				locked = false

				continue
			}

			//if organization deleted in git, delete in Agola, else update data in db
			gitOrganization, err := gitGateway.GetOrganization(gitSource, user, org.GitPath)
			if err != nil {
				log.Println("GetOrganization error:", err)

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
				locked = false

				continue
			}

			if gitOrganization == nil {
				log.Println("organization", organizationRef, "not found")

				err = agolaApi.DeleteOrganization(org, user)
				if err == nil {
					db.DeleteOrganization(organizationRef)
				} else {
					log.Println("error in agola DeleteOrganization:", err)
				}

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
				locked = false

				continue
			} else {
				if len(gitOrganization.Name) > 0 {
					org.GitName = gitOrganization.Name
				} else {
					org.GitName = gitOrganization.Path
				}
				db.SaveOrganization(org)
			}

			//If organization deleted in Agola, recreate
			if agolaOrganizationExists, _ := agolaApi.CheckOrganizationExists(org); !agolaOrganizationExists {
				orgID, err := agolaApi.CreateOrganization(org, org.Visibility)
				if err != nil {
					log.Println("failed to recreate organization", org.AgolaOrganizationRef, "in agola:", err)

					mutex.Unlock()
					utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
					locked = false

					continue
				}

				org.ID = orgID
				db.SaveOrganization(org)
			}

			log.Println("start synk organization", org.GitPath)

			membersManager.SynkMembers(org, gitSource, agolaApi, gitGateway, user)
			repositoryManager.SynkGitRepositorys(db, user, org, gitSource, agolaApi, gitGateway)

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
			locked = false
		}

		fmt.Println("syncOrganizationRun:", <-c)
	}
}
