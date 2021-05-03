package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
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
	AgolaApi    agolaApi.AgolaApiInterface
	GitGateway  *git.GitGateway
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

	forceCreateQuery, ok := r.URL.Query()["force"]
	forceCreate := false
	if ok {
		if len(forceCreateQuery[0]) == 0 {
			forceCreate = true
		} else {
			var parsError error
			forceCreate, parsError = strconv.ParseBool(forceCreateQuery[0])
			if parsError != nil {
				UnprocessableEntityResponse(w, "forceCreate param value is not valid")
				return
			}
		}
	}

	var req *dto.CreateOrganizationRequestDto
	json.NewDecoder(r.Body).Decode(&req)

	if !req.IsAgolaRefValid() {
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.AgolaRefNotValid}
		JSONokResponse(w, response)
		return
	}

	if req.IsValid() != nil {
		UnprocessableEntityResponse(w, "Parameters have no correct values")
		return
	}

	log.Println("Req CreateOrganizationDto: ", req)

	org := &model.Organization{}
	org.Name = req.Name
	org.AgolaOrganizationRef = req.AgolaRef
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

	gitOrgExists := service.GitGateway.CheckOrganizationExists(gitSource, org.Name)
	if gitOrgExists == false {
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.GitOrganizationNotFoundError}
		JSONokResponse(w, response)
		return
	}

	agolaOrg, err := service.Db.GetOrganizationByAgolaRef(org.AgolaOrganizationRef)
	if agolaOrg != nil {
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.PapagaioOrganizationExistsError}
		JSONokResponse(w, response)
		return
	}

	org.UserEmailCreator = r.Header.Get(controller.XAuthEmail)

	org.WebHookID, err = service.GitGateway.CreateWebHook(gitSource, org.Name, org.AgolaOrganizationRef)
	if err != nil {
		InternalServerError(w)
		return
	}

	agolaOrganizationExists, agolaOrganizationID := service.AgolaApi.CheckOrganizationExists(org)
	if agolaOrganizationExists {
		log.Println(org.Name, "just exists in Agola")
		if !forceCreate {
			service.GitGateway.DeleteWebHook(gitSource, org.Name, org.WebHookID)
			response := dto.CreateOrganizationResponseDto{ErrorCode: dto.AgolaOrganizationExistsError}
			JSONokResponse(w, response)
			return
		}
		org.ID = agolaOrganizationID
	} else {
		org.ID, err = service.AgolaApi.CreateOrganization(org, org.Visibility)
		if err != nil {
			log.Println("Agola CreateOrganization error:", err)
			service.GitGateway.DeleteWebHook(gitSource, org.Name, org.WebHookID)
			InternalServerError(w)
			return
		}
	}

	log.Println("Organization created: ", org.ID)
	log.Println("WebHook created: ", org.WebHookID)

	err = service.Db.SaveOrganization(org)
	if err != nil {
		InternalServerError(w)
		return
	}

	manager.StartOrganizationCheckout(service.Db, org, gitSource, service.AgolaApi, service.GitGateway)

	response := dto.CreateOrganizationResponseDto{OrganizationURL: utils.GetOrganizationUrl(org), ErrorCode: dto.NoError}
	JSONokResponse(w, response)
}

func (service *OrganizationService) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

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

	mutex := utils.ReserveOrganizationMutex(organizationRef, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(organizationRef, service.CommonMutex, mutex, &locked)

	organization, err := service.Db.GetOrganizationByAgolaRef(organizationRef)
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
		err = service.AgolaApi.DeleteOrganization(organization, gitSource.AgolaToken)
		if err != nil {
			InternalServerError(w)
			return
		}
	}

	service.GitGateway.DeleteWebHook(gitSource, organization.Name, organization.WebHookID)
	err = service.Db.DeleteOrganization(organization.AgolaOrganizationRef)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
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

	organization, err := service.Db.GetOrganizationByAgolaRef(organizationName)
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

	organization, err := service.Db.GetOrganizationByAgolaRef(organizationName)
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
		gitsource, _ := service.Db.GetGitSourceByName(organization.GitSourceName)
		retVal = append(retVal, manager.GetOrganizationDto(&organization, gitsource, service.GitGateway))
	}

	JSONokResponse(w, retVal)
}

func (service *OrganizationService) GetOrganizationReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationName := vars["organizationName"]

	organization, _ := service.Db.GetOrganizationByAgolaRef(organizationName)
	if organization == nil {
		NotFoundResponse(w)
		return
	}

	gitsource, _ := service.Db.GetGitSourceByName(organization.GitSourceName)

	JSONokResponse(w, manager.GetOrganizationDto(organization, gitsource, service.GitGateway))
}

func (service *OrganizationService) GetProjectReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationName := vars["organizationName"]
	projectName := vars["projectName"]

	organization, _ := service.Db.GetOrganizationByAgolaRef(organizationName)
	if organization == nil {
		NotFoundResponse(w)
		return
	}

	if organization.Projects == nil {
		NotFoundResponse(w)
		return
	} else if _, ok := organization.Projects[projectName]; !ok {
		NotFoundResponse(w)
		return
	}

	project := organization.Projects[projectName]

	JSONokResponse(w, manager.GetProjectDto(&project, organization))
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
