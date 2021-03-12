package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	agolaApi "wecode.sorint.it/opensource/papagaio-be/api/agola"
	"wecode.sorint.it/opensource/papagaio-be/dto"
	"wecode.sorint.it/opensource/papagaio-be/manager"
	"wecode.sorint.it/opensource/papagaio-be/model"
	"wecode.sorint.it/opensource/papagaio-be/repository"
	"wecode.sorint.it/opensource/papagaio-be/utils"
)

type WebHookService struct {
	Db repository.Database
}

func (service *WebHookService) WebHookOrganization(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebHookOrganization start...")

	data, _ := ioutil.ReadAll(r.Body)
	var webHookMessage dto.WebHookDto
	json.Unmarshal(data, &webHookMessage)

	log.Println("webHook message: ", webHookMessage)

	vars := mux.Vars(r)
	gitOrgRef := vars["gitOrgRef"]
	organization, _ := service.Db.GetOrganizationByName(gitOrgRef)
	if organization == nil {
		log.Println("Warning!!! Organization", gitOrgRef, "not found in db")
		return
	}

	if !utils.EvaluateBehaviour(organization, webHookMessage.Repository.Name) {
		log.Println("Web ", webHookMessage.Repository.Name, "excluded by behaviour settings")
		return
	}

	gitSource, _ := service.Db.GetGitSourceById(organization.GitSourceID)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	if strings.Compare(webHookMessage.Action, "created") == 0 {
		projectID, err := agolaApi.CreateProject(webHookMessage.Repository.Name, organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)
		if err != nil {
			fmt.Println("Warning!!! Agola CreateProject API error!")
			return
		}

		project := model.Project{OrganizationID: organization.ID, GitRepoPath: webHookMessage.Repository.Name, AgolaProjectID: projectID}
		organization.Projects[webHookMessage.Repository.Name] = project
		service.Db.SaveOrganization(organization)

		manager.StartSynkOrganization(service.Db, organization, gitSource)
	} else if strings.Compare(webHookMessage.Action, "deleted") == 0 {
		if orgProject, ok := organization.Projects[webHookMessage.Repository.Name]; !ok {
			log.Println("Warning!!! project", orgProject, "not found in db")
			return
		}

		agolaApi.DeleteProject(organization.Name, webHookMessage.Repository.Name, gitSource.AgolaToken)
		delete(organization.Projects, webHookMessage.Repository.Name)
		service.Db.SaveOrganization(organization)
	}

	fmt.Println("WebHookOrganization end...")
}
