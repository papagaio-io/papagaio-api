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
// @Success 200 {array} model.Organization "ok"
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

// @Summary Create a new Organization in Papagaio/Agola
// @Description Create an organization in Papagaio and in Agola. If already exists on Agola and you want to use the same organization then use the query parameter force
// @Tags Organization
// @Produce  json
// @Param force query string false "?force"
// @Param organization body dto.CreateOrganizationRequestDto true "Organization information"
// @Success 200 {object} dto.CreateOrganizationResponseDto "ok"
// @Failure 400 "bad request"-
// @Router /createorganization [post]
// @Security OAuth2Password
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

	userId, _ := strconv.ParseUint(r.Header.Get(controller.XAuthUserId), 10, 32)
	user, _ := service.Db.GetUserByUserId(userId)
	if user == nil {
		log.Println("User", userId, "not found")
		InternalServerError(w)
		return
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
	org.GitSourceName = user.GitSourceName
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

	gitOrgExists := service.GitGateway.CheckOrganizationExists(gitSource, user, org.Name)
	log.Println("gitOrgExists:", gitOrgExists)
	if !gitOrgExists {
		log.Println("failed to find organization", org.Name, "from git")
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.GitOrganizationNotFoundError}
		JSONokResponse(w, response)
		return
	}

	isOwner, err := service.GitGateway.IsUserOwner(gitSource, user, org.Name)
	if err != nil {
		log.Println("Error from IsUserOwner:", err)
		InternalServerError(w)
		return
	}
	if !isOwner {
		log.Println("User", user.UserID, "is not owner")
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.UserNotOwnerError}
		JSONokResponse(w, response)
		return
	}

	if user.AgolaUserRef == nil { //Se diverso da nil l'utente è registrato su Agola
		userInfo, err := service.GitGateway.GetUserInfo(gitSource, user)
		if err != nil || userInfo == nil {
			log.Println("Error on getting user info from git:", err)
			InternalServerError(w)
			return
		}

		agolaUserRef := utils.GetAgolaUserRefByGitUsername(service.AgolaApi, gitSource.AgolaRemoteSource, userInfo.Login)
		if agolaUserRef == nil {
			log.Println("User not found in Agola")
			response := dto.CreateOrganizationResponseDto{ErrorCode: dto.UserAgolaRefNotFoundError}
			JSONokResponse(w, response)
			return
		}

		user.AgolaUserRef = agolaUserRef
		err = service.AgolaApi.CreateUserToken(user)
		if err != nil {
			log.Println("Error in CreateUserToken:", err)
			InternalServerError(w)
			return
		}

		service.Db.SaveUser(user)
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

	org.UserIDCreator = *user.UserID

	org.WebHookID, err = service.GitGateway.CreateWebHook(gitSource, user, org.Name, org.AgolaOrganizationRef)
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
			service.GitGateway.DeleteWebHook(gitSource, user, org.Name, org.WebHookID)
			response := dto.CreateOrganizationResponseDto{ErrorCode: dto.AgolaOrganizationExistsError}
			JSONokResponse(w, response)
			return
		}
		org.ID = agolaOrganizationID
	} else {
		org.ID, err = service.AgolaApi.CreateOrganization(org, org.Visibility)
		if err != nil {
			log.Println("failed to create organization", org.AgolaOrganizationRef, "in agola:", err)
			service.GitGateway.DeleteWebHook(gitSource, user, org.Name, org.WebHookID)
			InternalServerError(w)
			return
		}
	}

	log.Println("Organization created: ", org.AgolaOrganizationRef, " by:", org.UserIDCreator)
	log.Println("WebHook created: ", org.WebHookID)

	err = service.Db.SaveOrganization(org)
	if err != nil {
		log.Println("failed to save organization in db")
		InternalServerError(w)
		return
	}

	manager.StartOrganizationCheckout(service.Db, user, org, gitSource, service.AgolaApi, service.GitGateway)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(org.AgolaOrganizationRef, service.CommonMutex)
	locked = false

	response := dto.CreateOrganizationResponseDto{OrganizationURL: utils.GetOrganizationUrl(org), ErrorCode: dto.NoError}
	JSONokResponse(w, response)
}

// @Summary Delete Organization
// @Description Delete an organization in Papagaio and in Agola. Its possible to delete only in Papagaio using the parameter internalonly.
// @Tags Organization
// @Produce  json
// @Param organizationRef path string true "Organization Name"
// @Param internalonly query string false "?internalonly"
// @Success 200 {object} model.Organization "ok"
// @Failure 500 "Not found"
// @Router /deleteorganization{organizationRef} [delete]
// @Security OAuth2Password
func (service *OrganizationService) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

	log.Println("DeleteOrganization:", organizationRef)

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

	userIdRequest, _ := strconv.ParseUint(r.Header.Get(controller.XAuthUserId), 10, 32)
	userRequest, _ := service.Db.GetUserByUserId(userIdRequest)
	if userRequest == nil {
		log.Println("User", userRequest, "not found")
		InternalServerError(w)
		return
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

	isOwner, err := service.GitGateway.IsUserOwner(gitSource, userRequest, organization.Name)
	if err != nil {
		log.Println("Error from IsUserOwner:", err)
		InternalServerError(w)
		return
	}
	if !isOwner {
		log.Println("User", userRequest.UserID, "is not owner")
		response := dto.DeleteOrganizationResponseDto{ErrorCode: dto.UserNotOwnerError}
		JSONokResponse(w, response)
		return
	}

	userCreator, _ := service.Db.GetUserByUserId(organization.UserIDCreator)
	if userCreator == nil {
		log.Println("User creator", organization.UserIDCreator, "not found")
		InternalServerError(w)
		return
	}

	if !internalonly {
		err = service.AgolaApi.DeleteOrganization(organization, userCreator)
		if err != nil {
			InternalServerError(w)
			return
		}
	}

	service.GitGateway.DeleteWebHook(gitSource, userCreator, organization.Name, organization.WebHookID)
	err = service.Db.DeleteOrganization(organization.AgolaOrganizationRef)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
	locked = false

	if err != nil {
		InternalServerError(w)
		return
	}
	log.Println("Organization deleted:", organization.AgolaOrganizationRef, " by:", userIdRequest)

	response := dto.DeleteOrganizationResponseDto{ErrorCode: dto.NoError}
	JSONokResponse(w, response)
}

// @Summary Add External User
// @Description Add an external user
// @Tags Organization
// @Produce  json
// @Param organizationRef path string true "Organization name"
// @Param email body string true "external user email"
// @Success 200 {object} model.Organization "ok"
// @Failure 404 "not found"
// @Router /addexternaluser/{organizationRef} [post]
// @Security OAuth2Password
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

// @Summary Delete External User
// @Description Delete an external user
// @Tags Organization
// @Produce  json
// @Param organizationRef path string true "Organization name"
// @Param email body string true "external user email"
// @Success 200 {object} model.Organization "ok"
// @Failure 404 "not found"
// @Router /deleteexternaluser/{organizationRef} [delete]
// @Security OAuth2Password
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

// @Summary Get Report
// @Description Obtain a full report of all organizations
// @Tags Organization
// @Produce  json
// @Success 200 {array} dto.OrganizationDto "ok"
// @Router /report [get]
// @Security OAuth2Password
func (service *OrganizationService) GetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	organizations, _ := service.Db.GetOrganizations()

	retVal := make([]dto.OrganizationDto, 0)
	for _, organization := range *organizations {
		gitsource, _ := service.Db.GetGitSourceByName(organization.GitSourceName)
		user, _ := service.Db.GetUserByUserId(organization.UserIDCreator)
		retVal = append(retVal, manager.GetOrganizationDto(user, &organization, gitsource, service.GitGateway))
	}

	JSONokResponse(w, retVal)
}

// @Summary Get Report from a specific organization
// @Description Obtain a report of a specific organization
// @Tags Organization
// @Produce  json
// @Success 200 {array} dto.OrganizationDto "ok"
// @Failure 404 "not found"
// @Router /report/{organizationRef} [get]
// @Security OAuth2Password
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
	user, _ := service.Db.GetUserByUserId(organization.UserIDCreator)

	JSONokResponse(w, manager.GetOrganizationDto(user, organization, gitsource, service.GitGateway))
}

// @Summary Get Report from a specific organization/project
// @Description Obtain a report of a specific organization/project
// @Tags Organization
// @Produce  json
// @Success 200 {array} dto.OrganizationDto "ok"
// @Failure 404 "not found"
// @Router /report/{organizationRef}/{projectName} [get]
// @Security OAuth2Password
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
// @Failure 400 "bad request"
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
			agolaRefList = append(agolaRefList, agolaOrganization.Name)
		}
	}

	JSONokResponse(w, agolaRefList)
}
