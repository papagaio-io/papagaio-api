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

			org, _ := db.GetOrganizationByAgolaRef(organizationRef)
			if org == nil {
				log.Println("syncOrganizationRun organization ", organizationRef, "not found")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			gitSource, _ := db.GetGitSourceByName(org.GitSourceName)
			if gitSource == nil {
				log.Println("gitSource", org.GitSourceName, "not found")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			user, _ := db.GetUserByUserId(org.UserIDConnected)
			if user == nil {
				log.Println("user not found")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			//if organization deleted in git, delete in Agola, else update data in db
			gitOrganization, err := gitGateway.GetOrganization(gitSource, user, org.GitPath)
			if err != nil {
				log.Println("GetOrganization error:", err)

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			if gitOrganization == nil {
				log.Println("organization", organizationRef, "not found")

				err = agolaApi.DeleteOrganization(org, user)
				if err == nil {
					err := db.DeleteOrganization(organizationRef)
					if err != nil {
						log.Println("error in DeleteOrganization:", err)
					}
				} else {
					log.Println("error in agola DeleteOrganization:", err)
				}

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			} else {
				if len(gitOrganization.Name) > 0 {
					org.GitName = gitOrganization.Name
				} else {
					org.GitName = gitOrganization.Path
				}
				err := db.SaveOrganization(org)

				if err != nil {
					log.Println("error in SaveOrganization:", err)
				}
			}

			//If organization deleted in Agola, recreate
			agolaOrganizationExists, _, err := agolaApi.CheckOrganizationExists(org)
			if err != nil {
				log.Println("Agola CheckOrganizationExists error:", err)

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			if !agolaOrganizationExists {
				orgID, err := agolaApi.CreateOrganization(org, org.Visibility)
				if err != nil {
					log.Println("failed to recreate organization", org.AgolaOrganizationRef, "in agola:", err)

					mutex.Unlock()
					utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

					continue
				}

				org.ID = orgID
				err = db.SaveOrganization(org)

				if err != nil {
					log.Println("error in SaveOrganization:", err)
				}
			}

			log.Println("start synk organization", org.GitPath)

			err = membersManager.SynkMembers(org, gitSource, agolaApi, gitGateway, user)
			if err != nil {
				log.Println("SynkMembers error:", err)
			}

			err = repositoryManager.SynkGitRepositorys(db, user, org, gitSource, agolaApi, gitGateway)
			if err != nil {
				log.Println("SynkGitRepositorys error:", err)
			}

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
		}

		fmt.Println("syncOrganizationRun:", <-c)
	}
}
