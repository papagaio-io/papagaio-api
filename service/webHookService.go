package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xanzy/go-gitlab"
	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/types"
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

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

	mutex := utils.ReserveOrganizationMutex(organizationRef, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(organizationRef, service.CommonMutex, mutex, &locked)

	organization, _ := service.Db.GetOrganizationByAgolaRef(organizationRef)
	if organization == nil {
		log.Println("warning!!! Organization", organizationRef, "not found in db")
		InternalServerError(w)
		return
	}

	gitSource, _ := service.Db.GetGitSourceByName(organization.GitSourceName)
	if gitSource == nil {
		log.Println("gitSource", organization.GitSourceName, "not found")
		InternalServerError(w)
		return
	}

	var webHookMessage dto.WebHookDto
	data, _ := ioutil.ReadAll(r.Body)

	if gitSource.GitType == types.Gitlab {
		var gitLabHookMessage gitlab.PushEvent
		err := json.Unmarshal(data, &gitLabHookMessage)
		if err != nil {
			log.Println("gitlab unmarshal error:", err)
			InternalServerError(w)
			return
		}

		webHookMessage.Action = ""
		webHookMessage.Sha = gitLabHookMessage.CheckoutSHA
		webHookMessage.Repository.Name = gitLabHookMessage.Repository.Name
		webHookMessage.Repository.ID = gitLabHookMessage.ProjectID
	} else {
		err := json.Unmarshal(data, &webHookMessage)
		if err != nil {
			log.Println("unmarshal error:", err)
			InternalServerError(w)
			return
		}
	}

	log.Println("webHook message: ", webHookMessage)

	if !utils.EvaluateBehaviour(organization, webHookMessage.Repository.Name) {
		log.Println("webhook", webHookMessage.Repository.Name, "excluded by behaviour settings")
		UnprocessableEntityResponse(w, "behaviour exclude")
		return
	}

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	user, _ := service.Db.GetUserByUserId(organization.UserIDConnected)
	if user == nil {
		log.Println("user not found")
		InternalServerError(w)
		return
	}

	if webHookMessage.IsRepositoryCreated() {
		log.Println("repository created: ", webHookMessage.Repository.Name)
		project := model.Project{GitRepoPath: webHookMessage.Repository.Name, Archivied: true, AgolaProjectRef: utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)}

		agolaConfExists, _ := service.GitGateway.CheckRepositoryAgolaConfExists(gitSource, user, organization.GitPath, webHookMessage.Repository.Name)
		if agolaConfExists {
			projectID, err := service.AgolaApi.CreateProject(webHookMessage.Repository.Name, project.AgolaProjectRef, organization, gitSource.AgolaRemoteSource, user)
			project.AgolaProjectID = projectID
			if err != nil {
				log.Println("warning!!! Agola CreateProject API error!")
				InternalServerError(w)
				return
			} else {
				project.Archivied = false
			}
		}

		organization.Projects[webHookMessage.Repository.Name] = project
		err := service.Db.SaveOrganization(organization)

		if err != nil {
			log.Println("SaveOrganization error:", err)
			InternalServerError(w)
			return
		}
	} else if webHookMessage.IsRepositoryDeleted() {
		log.Println("repository deleted: ", webHookMessage.Repository.Name)

		orgProject, ok := organization.Projects[webHookMessage.Repository.Name]

		if !ok {
			log.Println("warning!!! project", orgProject, "not found in db")
			UnprocessableEntityResponse(w, "repository not found")
			return
		}

		err := service.AgolaApi.DeleteProject(organization, orgProject.AgolaProjectRef, user)
		if err != nil {
			log.Println("agola DeleteProject error:", err)
			InternalServerError(w)
			return
		}

		delete(organization.Projects, webHookMessage.Repository.Name)

		project := organization.Projects[webHookMessage.Repository.Name]
		project.Archivied = true

		err = service.Db.SaveOrganization(organization)

		if err != nil {
			log.Println("SaveOrganization error:", err)
			InternalServerError(w)
			return
		}
	} else if webHookMessage.IsPush() {
		log.Println("repository push: ", webHookMessage.Repository.Name)

		project, projectExist := organization.Projects[webHookMessage.Repository.Name]
		agolaConfExists, _ := service.GitGateway.CheckRepositoryAgolaConfExists(gitSource, user, organization.GitPath, webHookMessage.Repository.Name)

		if agolaConfExists {
			if !projectExist || !project.ExistsInAgola() {
				agolaProjectRef := utils.ConvertToAgolaProjectRef(webHookMessage.Repository.Name)
				projectID, err := service.AgolaApi.CreateProject(webHookMessage.Repository.Name, agolaProjectRef, organization, gitSource.AgolaRemoteSource, user)
				if err != nil {
					log.Println("warning!!! Agola CreateProject API error!")
					InternalServerError(w)
					return
				}

				project := model.Project{GitRepoPath: webHookMessage.Repository.Name, AgolaProjectID: projectID, AgolaProjectRef: agolaProjectRef}
				organization.Projects[webHookMessage.Repository.Name] = project
				err = service.Db.SaveOrganization(organization)

				if err != nil {
					log.Println("SaveOrganization error:", err)
					InternalServerError(w)
					return
				}
			} else if project.Archivied {
				err := service.AgolaApi.UnarchiveProject(organization, project.AgolaProjectRef)
				if err != nil {
					log.Println("UnarchiveProject error:", err)
					InternalServerError(w)
					return
				}
				project.Archivied = false
				organization.Projects[webHookMessage.Repository.Name] = project
				err = service.Db.SaveOrganization(organization)
				if err != nil {
					log.Println("SaveOrganization error:", err)
					InternalServerError(w)
					return
				}
			}
		} else {
			if projectExist && !project.Archivied {
				err := service.AgolaApi.ArchiveProject(organization, project.AgolaProjectRef)
				if err != nil {
					log.Println("ArchiveProject error:", err)
					InternalServerError(w)
					return
				}
				project.Archivied = true
				organization.Projects[webHookMessage.Repository.Name] = project
				err = service.Db.SaveOrganization(organization)

				if err != nil {
					log.Println("SaveOrganization error:", err)
					InternalServerError(w)
					return
				}
			}
		}

		repositoryManager.BranchSynck(service.Db, user, gitSource, organization, webHookMessage.Repository.Name, service.GitGateway)
	}

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
	locked = false

	log.Println("WebHookOrganization end...")
}
