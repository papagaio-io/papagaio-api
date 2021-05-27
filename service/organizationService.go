package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

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

// @Summary Return the list of organizations
// @Description Return the list of organizations
// @Tags Organization
// @Produce  json
// @Success 200 {object} model.Organization "ok"
// @Failure 400 "bad request"-
// @Router /organizations [get]
// @Security OAuth2Password
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
		UnprocessableEntityResponse(w, "parameters have no correct values")
		return
	}

	log.Println("req CreateOrganizationDto: ", req)

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
		log.Println("failed to find gitSource", org.GitSourceName, "from db")
		UnprocessableEntityResponse(w, "Gitsource non found")
		return
	}

	gitOrgExists := service.GitGateway.CheckOrganizationExists(gitSource, org.Name)
	log.Println("gitOrgExists:", gitOrgExists)
	if !gitOrgExists {
		log.Println("failed to find organization", org.Name, "from git")
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.GitOrganizationNotFoundError}
		JSONokResponse(w, response)
		return
	}

	mutex := utils.ReserveOrganizationMutex(org.AgolaOrganizationRef, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(org.AgolaOrganizationRef, service.CommonMutex, mutex, &locked)

	agolaOrg, _ := service.Db.GetOrganizationByAgolaRef(org.AgolaOrganizationRef)
	if agolaOrg != nil {
		log.Println("organization", org.AgolaOrganizationRef, "just exists")
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.PapagaioOrganizationExistsError}
		JSONokResponse(w, response)
		return
	} else {
		organizations, _ := service.Db.GetOrganizationsByGitSource(org.GitSourceName)
		for _, organization := range *organizations {
			if strings.Compare(organization.Name, org.Name) == 0 {
				log.Println("organization name", org.Name, "just present in papagaio with gitSource", org.GitSourceName)

				response := dto.CreateOrganizationResponseDto{ErrorCode: dto.PapagaioOrganizationExistsError}
				JSONokResponse(w, response)
				return
			}
		}
	}

	org.UserEmailCreator = r.Header.Get(controller.XAuthEmail)

	org.WebHookID, err = service.GitGateway.CreateWebHook(gitSource, org.Name, org.AgolaOrganizationRef)
	if err != nil {
		log.Println("failed to creare webhook")
		InternalServerError(w)
		return
	}

	agolaOrganizationExists, agolaOrganizationID := service.AgolaApi.CheckOrganizationExists(org)
	log.Println("agolaOrganizationExists:", agolaOrganizationExists)
	if agolaOrganizationExists {
		log.Println("organization", org.AgolaOrganizationRef, "just exists in Agola")
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
			log.Println("failed to create organization", org.AgolaOrganizationRef, "in agola:", err)
			service.GitGateway.DeleteWebHook(gitSource, org.Name, org.WebHookID)
			InternalServerError(w)
			return
		}
	}

	log.Println("Organization created: ", org.AgolaOrganizationRef, " by:", org.UserEmailCreator)
	log.Println("WebHook created: ", org.WebHookID)

	err = service.Db.SaveOrganization(org)
	if err != nil {
		log.Println("failed to save organization in db")
		InternalServerError(w)
		return
	}

	manager.StartOrganizationCheckout(service.Db, org, gitSource, service.AgolaApi, service.GitGateway)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(org.AgolaOrganizationRef, service.CommonMutex)
	locked = false

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
		log.Println("organiztion", organizationRef, "not found in db")
		NotFoundResponse(w)
		return
	}

	gitSource, err := service.Db.GetGitSourceByName(organization.GitSourceName)
	if err != nil || gitSource == nil {
		log.Println("gitSource", organization.GitSourceName, "not found in db")
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
	log.Println("Organization deleted:", organization.AgolaOrganizationRef, " by:", organization.UserEmailCreator)
}

func (service *OrganizationService) AddExternalUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

	mutex := utils.ReserveOrganizationMutex(organizationRef, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(organizationRef, service.CommonMutex, mutex, &locked)

	organization, err := service.Db.GetOrganizationByAgolaRef(organizationRef)
	if err != nil || organization == nil {
		NotFoundResponse(w)
		return
	}

	var req *dto.ExternalUserDto
	json.NewDecoder(r.Body).Decode(&req)

	if organization.ExternalUsers == nil {
		organization.ExternalUsers = make(map[string]bool)
	}

	organization.ExternalUsers[req.Email] = true
	service.Db.SaveOrganization(organization)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
	locked = false
}

func (service *OrganizationService) RemoveExternalUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

	mutex := utils.ReserveOrganizationMutex(organizationRef, service.CommonMutex)
	mutex.Lock()

	locked := true
	defer utils.ReleaseOrganizationMutexDefer(organizationRef, service.CommonMutex, mutex, &locked)

	organization, err := service.Db.GetOrganizationByAgolaRef(organizationRef)
	if err != nil || organization == nil {
		NotFoundResponse(w)
		return
	}

	var req *dto.ExternalUserDto
	json.NewDecoder(r.Body).Decode(&req)

	delete(organization.ExternalUsers, req.Email)
	service.Db.SaveOrganization(organization)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
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
	organizationRef := vars["organizationRef"]

	organization, _ := service.Db.GetOrganizationByAgolaRef(organizationRef)
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
	organizationRef := vars["organizationRef"]
	projectName := vars["projectName"]

	organization, _ := service.Db.GetOrganizationByAgolaRef(organizationRef)
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

// @Summary Return the organization ref list
// @Description Return the organization ref list existing in Agola but not in Papagaio
// @Tags Organization
// @Produce  json
// @Success 200 {array} string "ok"
// @Failure 400 "bad request"-
// @Router /agolarefs [get]
// @Security OAuth2Password
func (service *OrganizationService) GetAgolaOrganizations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	organizations, err := service.AgolaApi.GetOrganizations()
	if err != nil {
		InternalServerError(w)
		return
	}

	agolaRefList := make([]string, 0)
	if organizations != nil {
		for _, agolaOrganization := range *organizations {
			dbOrganization, _ := service.Db.GetOrganizationByAgolaRef(agolaOrganization.Name)
			if dbOrganization == nil {
				agolaRefList = append(agolaRefList, agolaOrganization.Name)
			}
		}
	}

	JSONokResponse(w, agolaRefList)
}
