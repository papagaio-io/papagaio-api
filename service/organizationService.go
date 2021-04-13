package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	gitApi "wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/manager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

type OrganizationService struct {
	Db          repository.Database
	CommonMutex *utils.CommonMutex
}

func (service *OrganizationService) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	org, err := service.Db.GetOrganizations()
	if err != nil {
		InternalServerError(w)
		return
	}

	JSONokResponse(w, org)
}

func (service *OrganizationService) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req *dto.CreateOrganizationDto
	json.NewDecoder(r.Body).Decode(&req)

	if req.IsValid() != nil {
		UnprocessableEntityResponse(w, "Parameters have no correct values")
		return
	}

	log.Println("Req CreateOrganizationDto: ", req)

	org := &model.Organization{}
	org.Name = req.Name
	org.GitSourceName = req.GitSourceName
	org.Visibility = req.Visibility
	org.BehaviourType = req.BehaviourType
	org.BehaviourInclude = req.BehaviourInclude
	org.BehaviourExclude = req.BehaviourExclude

	//Some checks
	gitSource, err := service.Db.GetGitSourceByName(org.GitSourceName)
	if gitSource == nil || err != nil {
		UnprocessableEntityResponse(w, "Gitsource non found")
		return
	}

	gitOrgExists := gitApi.CheckOrganizationExists(gitSource, org.Name)
	if gitOrgExists == false {
		UnprocessableEntityResponse(w, "Organization not found")
		return
	}

	agolaOrg, err := service.Db.GetOrganizationByName(org.Name)
	if agolaOrg != nil {
		UnprocessableEntityResponse(w, "Organization just present in Agola")
		return
	}

	org.UserEmailCreator = r.Header.Get(controller.XAuthEmail)

	org.WebHookID, err = gitApi.CreateWebHook(gitSource, org.Name)
	if err != nil {
		UnprocessableEntityResponse(w, err.Error())
		return
	}

	org.ID, err = agolaApi.CreateOrganization(org.Name, org.Visibility)
	if err != nil {
		log.Println("Agola CreateOrganization error")
		InternalServerError(w)
		return
	}

	log.Println("Organization created: ", org.ID)
	log.Println("WebHook created: ", org.WebHookID)

	err = service.Db.SaveOrganization(org)
	if err != nil {
		InternalServerError(w)
		return
	}

	manager.StartOrganizationCheckout(service.Db, org, gitSource)

	JSONokResponse(w, config.Config.Agola.AgolaAddr+"/org/"+org.Name)
}

func (service *OrganizationService) GetRemoteSources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	remoteSources, _ := agolaApi.GetRemoteSources()

	JSONokResponse(w, remoteSources)
}

func (service *OrganizationService) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationName := vars["organizationName"]

	internalonlyQuery, ok := r.URL.Query()["internalonly"]
	internalonly := false
	if ok {
		if len(internalonlyQuery[0]) == 0 {
			internalonly = true
		} else {
			var parsError error
			internalonly, parsError = strconv.ParseBool(internalonlyQuery[0])
			if parsError != nil {
				UnprocessableEntityResponse(w, "internalonly param value is not valid")
				return
			}
		}
	}

	mutex := utils.ReserveOrganizationMutex(organizationName, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(organizationName, service.CommonMutex, mutex, &locked)

	organization, err := service.Db.GetOrganizationByName(organizationName)
	if err != nil || organization == nil {
		NotFoundResponse(w)
		return
	}

	gitSource, err := service.Db.GetGitSourceByName(organization.GitSourceName)
	if err != nil || gitSource == nil {
		InternalServerError(w)
		return
	}

	if !internalonly {
		err = agolaApi.DeleteOrganization(organizationName, gitSource.AgolaToken)
		if err != nil {
			InternalServerError(w)
			return
		}
	}

	gitApi.DeleteWebHook(gitSource, organizationName, organization.WebHookID)
	err = service.Db.DeleteOrganization(organization.Name)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationName, service.CommonMutex)
	locked = false

	if err != nil {
		InternalServerError(w)
	}
}

func (service *OrganizationService) AddExternalUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationName := vars["organizationName"]

	mutex := utils.ReserveOrganizationMutex(organizationName, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(organizationName, service.CommonMutex, mutex, &locked)

	organization, err := service.Db.GetOrganizationByName(organizationName)
	if err != nil || organization == nil {
		NotFoundResponse(w)
		return
	}

	var req *dto.ExternalUserDto
	json.NewDecoder(r.Body).Decode(&req)

	organization.ExternalUsers[req.Email] = true
	service.Db.SaveOrganization(organization)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationName, service.CommonMutex)
	locked = false
}

func (service *OrganizationService) RemoveExternalUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationName := vars["organizationName"]

	mutex := utils.ReserveOrganizationMutex(organizationName, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(organizationName, service.CommonMutex, mutex, &locked)

	organization, err := service.Db.GetOrganizationByName(organizationName)
	if err != nil || organization == nil {
		NotFoundResponse(w)
		return
	}

	var req *dto.ExternalUserDto
	json.NewDecoder(r.Body).Decode(&req)

	delete(organization.ExternalUsers, req.Email)
	service.Db.SaveOrganization(organization)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationName, service.CommonMutex)
	locked = false
}

func (service *OrganizationService) GetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	organizations, _ := service.Db.GetOrganizations()

	retVal := make([]dto.OrganizationDto, 0)
	for _, organization := range *organizations {
		retVal = append(retVal, manager.GetOrganizationDto(&organization))
	}

	JSONokResponse(w, retVal)
}

func (service *OrganizationService) GetOrganizationReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationName := vars["organizationName"]

	organization, _ := service.Db.GetOrganizationByName(organizationName)
	if organization == nil {
		NotFoundResponse(w)
		return
	}

	JSONokResponse(w, manager.GetOrganizationDto(organization))
}

func (service *OrganizationService) GetProjectReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationName := vars["organizationName"]
	projectName := vars["projectName"]

	organization, _ := service.Db.GetOrganizationByName(organizationName)
	if organization == nil {
		NotFoundResponse(w)
		return
	}

	if organization.Projects == nil {
		NotFoundResponse(w)
		return
	} else if project, ok := organization.Projects[projectName]; !ok || project.Archivied {
		NotFoundResponse(w)
		return
	}

	project := organization.Projects[projectName]

	JSONokResponse(w, manager.GetProjectDto(&project, organizationName))
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
