package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

type WebHookService struct {
	Db          repository.Database
	CommonMutex *utils.CommonMutex
	AgolaApi    agolaApi.AgolaApiInterface
	GitGateway  *git.GitGateway
}

func (service *WebHookService) WebHookOrganization(w http.ResponseWriter, r *http.Request) {
	log.Println("WebHookOrganization start...")

	data, _ := ioutil.ReadAll(r.Body)
	var webHookMessage dto.WebHookDto
	json.Unmarshal(data, &webHookMessage)

	log.Println("webHook message: ", webHookMessage)

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

	mutex := utils.ReserveOrganizationMutex(organizationRef, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(organizationRef, service.CommonMutex, mutex, &locked)

	organization, _ := service.Db.GetOrganizationByAgolaRef(organizationRef)
	if organization == nil {
		log.Println("warning!!! Organization", organizationRef, "not found in db")
		return
	}

	if !utils.EvaluateBehaviour(organization, webHookMessage.Repository.Name) {
		log.Println("webhook", webHookMessage.Repository.Name, "excluded by behaviour settings")
		return
	}

	gitSource, _ := service.Db.GetGitSourceByName(organization.GitSourceName)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	if webHookMessage.IsRepositoryCreated() {
		log.Println("repository created: ", webHookMessage.Repository.Name)
		project := model.Project{GitRepoPath: webHookMessage.Repository.Name, Archivied: true, AgolaProjectRef: utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)}

		agolaConfExists, _ := service.GitGateway.CheckRepositoryAgolaConfExists(gitSource, organization.Name, webHookMessage.Repository.Name)
		if agolaConfExists {
			projectID, err := service.AgolaApi.CreateProject(webHookMessage.Repository.Name, project.AgolaProjectRef, organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)
			project.AgolaProjectID = projectID
			if err != nil {
				log.Println("warning!!! Agola CreateProject API error!")
			} else {
				project.Archivied = false
			}
		}

		organization.Projects[webHookMessage.Repository.Name] = project
		service.Db.SaveOrganization(organization)
	} else if webHookMessage.IsRepositoryDeleted() {
		log.Println("repository deleted: ", webHookMessage.Repository.Name)

		orgProject, ok := organization.Projects[webHookMessage.Repository.Name]

		if !ok {
			log.Println("warning!!! project", orgProject, "not found in db")
			return
		}

		service.AgolaApi.DeleteProject(organization, orgProject.AgolaProjectRef, gitSource.AgolaToken)
		delete(organization.Projects, webHookMessage.Repository.Name)

		project := organization.Projects[webHookMessage.Repository.Name]
		project.Archivied = true

		service.Db.SaveOrganization(organization)
	} else if webHookMessage.IsPush() {
		log.Println("repository push: ", webHookMessage.Repository.Name)

		project, projectExist := organization.Projects[webHookMessage.Repository.Name]
		agolaConfExists, _ := service.GitGateway.CheckRepositoryAgolaConfExists(gitSource, organization.Name, webHookMessage.Repository.Name)

		if agolaConfExists {
			if !projectExist {
				agolaProjectRef := utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)
				projectID, err := service.AgolaApi.CreateProject(webHookMessage.Repository.Name, agolaProjectRef, organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)
				if err != nil {
					log.Println("warning!!! Agola CreateProject API error!")
					return
				}

				project := model.Project{GitRepoPath: webHookMessage.Repository.Name, AgolaProjectID: projectID, AgolaProjectRef: agolaProjectRef}
				organization.Projects[webHookMessage.Repository.Name] = project
				service.Db.SaveOrganization(organization)
			} else if project.Archivied {
				err := service.AgolaApi.UnarchiveProject(organization, project.AgolaProjectRef)
				if err != nil {
					InternalServerError(w)
				}
				project.Archivied = false
				organization.Projects[webHookMessage.Repository.Name] = project
				service.Db.SaveOrganization(organization)
			}
		} else {
			if projectExist && !project.Archivied {
				err := service.AgolaApi.ArchiveProject(organization, project.AgolaProjectRef)
				if err != nil {
					InternalServerError(w)
				}
				project.Archivied = true
				organization.Projects[webHookMessage.Repository.Name] = project
				service.Db.SaveOrganization(organization)
			}
		}

		repositoryManager.BranchSynck(service.Db, gitSource, organization, webHookMessage.Repository.Name, service.GitGateway)
	}

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
	locked = false

	log.Println("WebHookOrganization end...")
}
