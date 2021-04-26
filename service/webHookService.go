package service

import (
	"encoding/json"
	"fmt"
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
}

func (service *WebHookService) WebHookOrganization(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebHookOrganization start...")

	data, _ := ioutil.ReadAll(r.Body)
	var webHookMessage dto.WebHookDto
	json.Unmarshal(data, &webHookMessage)

	log.Println("webHook message: ", webHookMessage)

	vars := mux.Vars(r)
	gitOrgRef := vars["gitOrgRef"]

	mutex := utils.ReserveOrganizationMutex(gitOrgRef, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(gitOrgRef, service.CommonMutex, mutex, &locked)

	organization, _ := service.Db.GetOrganizationByName(gitOrgRef)
	if organization == nil {
		log.Println("Warning!!! Organization", gitOrgRef, "not found in db")
		return
	}

	if !utils.EvaluateBehaviour(organization, webHookMessage.Repository.Name) {
		log.Println("Web ", webHookMessage.Repository.Name, "excluded by behaviour settings")
		return
	}

	gitSource, _ := service.Db.GetGitSourceByName(organization.GitSourceName)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	if webHookMessage.IsRepositoryCreated() {
		log.Println("Repository created: ", webHookMessage.Repository.Name)
		project := model.Project{GitRepoPath: webHookMessage.Repository.Name, Archivied: true, AgolaProjectRef: utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)}

		agolaConfExists, _ := git.CheckRepositoryAgolaConf(gitSource, organization.Name, webHookMessage.Repository.Name)
		if agolaConfExists {
			projectID, err := agolaApi.CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)
			project.AgolaProjectID = projectID
			if err != nil {
				log.Println("Warning!!! Agola CreateProject API error!")
			} else {
				project.Archivied = false
			}
		}

		organization.Projects[webHookMessage.Repository.Name] = project
		service.Db.SaveOrganization(organization)
	} else if webHookMessage.IsRepositoryDeleted() {
		log.Println("Repository deleted: ", webHookMessage.Repository.Name)

		if orgProject, ok := organization.Projects[webHookMessage.Repository.Name]; !ok {
			log.Println("Warning!!! project", orgProject, "not found in db")
			return
		}

		agolaApi.DeleteProject(organization, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), gitSource.AgolaToken)
		delete(organization.Projects, webHookMessage.Repository.Name)

		project := organization.Projects[webHookMessage.Repository.Name]
		project.Archivied = true

		service.Db.SaveOrganization(organization)
	} else if webHookMessage.IsPush() {
		log.Println("Repository push: ", webHookMessage.Repository.Name)

		project, projectExist := organization.Projects[webHookMessage.Repository.Name]
		agolaConfExists, _ := git.CheckRepositoryAgolaConf(gitSource, organization.Name, webHookMessage.Repository.Name)

		if agolaConfExists {
			if !projectExist {
				projectID, err := agolaApi.CreateProject(webHookMessage.Repository.Name, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name), organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)
				if err != nil {
					fmt.Println("Warning!!! Agola CreateProject API error!")
					return
				}

				project := model.Project{GitRepoPath: webHookMessage.Repository.Name, AgolaProjectID: projectID, AgolaProjectRef: utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)}
				organization.Projects[webHookMessage.Repository.Name] = project
				service.Db.SaveOrganization(organization)
			} else if project.Archivied {
				err := agolaApi.UnarchiveProject(organization, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name))
				if err != nil {
					InternalServerError(w)
				}
				project.Archivied = false
				service.Db.SaveOrganization(organization)
			}
		} else {
			if projectExist && !project.Archivied {
				err := agolaApi.ArchiveProject(organization, utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name))
				if err != nil {
					InternalServerError(w)
				}
				project.Archivied = true
				service.Db.SaveOrganization(organization)
			}
		}

		repositoryManager.BranchSynck(service.Db, gitSource, organization, webHookMessage.Repository.Name)
	}

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(gitOrgRef, service.CommonMutex)
	locked = false

	fmt.Println("WebHookOrganization end...")
}
