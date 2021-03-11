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
	"wecode.sorint.it/opensource/papagaio-be/model"
	"wecode.sorint.it/opensource/papagaio-be/repository"
)

type WebHookService struct {
	Db repository.Database
}

func (service *WebHookService) WebHookOrganization(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebHookOrganization start...")

	data, _ := ioutil.ReadAll(r.Body)
	var webHookMessage *dto.WebHookDto
	json.Unmarshal(data, webHookMessage)

	log.Println("webHook message: ", webHookMessage)

	vars := mux.Vars(r)
	gitOrgRef := vars["gitOrgRef"]
	organization, _ := service.Db.GetOrganizationByName(gitOrgRef)
	if organization == nil {
		return
	}
	gitSource, _ := service.Db.GetGitSourceById(organization.GitSourceID)

	if strings.Compare(webHookMessage.Action, "created") == 0 {
		projectID, _ := agolaApi.CreateProject(webHookMessage.Repository.Name, organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)
		project := model.Project{OrganizationID: organization.ID, GitRepoPath: webHookMessage.Repository.Name, AgolaProjectRef: projectID}
		organization.Projects = append(organization.Projects, project)
		service.Db.SaveOrganization(organization)
	}
}
