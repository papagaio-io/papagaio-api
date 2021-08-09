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
// @Param force query bool false "?force"
// @Param organization body dto.CreateOrganizationRequestDto true "Organization information"
// @Success 200 {object} dto.CreateOrganizationResponseDto "ok"
// @Failure 400 "bad request"-
// @Router /createorganization [post]
// @Security ApiKeyToken
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

	userId := r.Context().Value(controller.UserIdParameter).(uint64)
	user, _ := service.Db.GetUserByUserId(userId)
	if user == nil {
		log.Println("User", userId, "not found")
		InternalServerError(w)
		return
	}

	var req *dto.CreateOrganizationRequestDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("parsing error:", err)
		InternalServerError(w)
		return
	}

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
	org.GitPath = req.GitPath
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

	gitOrganization, _ := service.GitGateway.GetOrganization(gitSource, user, org.GitPath)
	log.Println("gitOrgExists:", gitOrganization != nil)
	if gitOrganization == nil {
		log.Println("failed to find organization", org.GitPath, "from git")
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.GitOrganizationNotFoundError}
		JSONokResponse(w, response)
		return
	}
	org.GitOrganizationID = gitOrganization.ID
	if len(gitOrganization.Name) > 0 {
		org.GitName = gitOrganization.Name
	} else {
		org.GitName = req.GitPath
	}

	isOwner, _ := service.GitGateway.IsUserOwner(gitSource, user, org.GitPath)
	if !isOwner {
		log.Println("User", user.UserID, "is not owner")
		response := dto.CreateOrganizationResponseDto{ErrorCode: dto.UserNotOwnerError}
		JSONokResponse(w, response)
		return
	}

	if user.AgolaUserRef == nil { //Se diverso da nil l'utente è registrato su Agola
		agolaUserRef := utils.GetAgolaUserRefByGitUsername(service.AgolaApi, gitSource.AgolaRemoteSource, user.Login)
		if agolaUserRef == nil {
			log.Println("User not found in Agola")
			response := dto.CreateOrganizationResponseDto{ErrorCode: dto.UserAgolaRefNotFoundError}
			JSONokResponse(w, response)
			return
		}

		user.AgolaUserRef = agolaUserRef

		if user.AgolaToken == nil {
			err = service.AgolaApi.CreateUserToken(user)
			if err != nil {
				log.Println("Error in CreateUserToken:", err)
				InternalServerError(w)
				return
			}
		}

		err := service.Db.SaveUser(user)
		if err != nil {
			log.Println("Error in SaveUser:", err)
			InternalServerError(w)
			return
		}
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
			if strings.Compare(organization.GitPath, org.GitPath) == 0 {
				log.Println("organization name", org.GitPath, "just present in papagaio with gitSource", org.GitSourceName)

				response := dto.CreateOrganizationResponseDto{ErrorCode: dto.PapagaioOrganizationExistsError}
				JSONokResponse(w, response)
				return
			}
		}
	}

	org.UserIDCreator = *user.UserID
	org.UserIDConnected = *user.UserID

	org.WebHookID, err = service.GitGateway.CreateWebHook(gitSource, user, org.GitPath, org.AgolaOrganizationRef)
	if err != nil {
		log.Println("failed to creare webhook:", err)
		InternalServerError(w)
		return
	}

	agolaOrganizationExists, agolaOrganizationID, err := service.AgolaApi.CheckOrganizationExists(org)
	if err != nil {
		log.Println("Agola CheckOrganizationExists error:", err)
		InternalServerError(w)
		return
	}

	log.Println("agolaOrganizationExists:", agolaOrganizationExists)
	if agolaOrganizationExists {
		log.Println("organization", org.AgolaOrganizationRef, "just exists in Agola")
		if !forceCreate {
			err = service.GitGateway.DeleteWebHook(gitSource, user, org.GitPath, org.WebHookID)
			if err != nil {
				log.Println("DeleteWebHook error:", err)
			}

			response := dto.CreateOrganizationResponseDto{ErrorCode: dto.AgolaOrganizationExistsError}
			JSONokResponse(w, response)
			return
		}
		org.ID = agolaOrganizationID
	} else {
		org.ID, err = service.AgolaApi.CreateOrganization(org, org.Visibility)
		if err != nil {
			log.Println("failed to create organization", org.AgolaOrganizationRef, "in agola:", err)
			err := service.GitGateway.DeleteWebHook(gitSource, user, org.GitPath, org.WebHookID)
			if err != nil {
				log.Println("DeleteWebHook error:", err)
			}

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
// @Success 200 {object} dto.DeleteOrganizationResponseDto "ok"
// @Failure 500 "Not found"
// @Router /deleteorganization{organizationRef} [delete]
// @Security ApiKeyToken
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

	userIdRequest, _ := r.Context().Value(controller.UserIdParameter).(uint64)
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

	isOwner, _ := service.GitGateway.IsUserOwner(gitSource, userRequest, organization.GitPath)
	if !isOwner {
		log.Println("User", userRequest.UserID, "is not owner")
		response := dto.DeleteOrganizationResponseDto{ErrorCode: dto.UserNotOwnerError}
		JSONokResponse(w, response)
		return
	}

	userCreator, _ := service.Db.GetUserByUserId(organization.UserIDConnected)
	if userCreator == nil {
		log.Println("User ", organization.UserIDConnected, "not found")
		InternalServerError(w)
		return
	}

	if !internalonly {
		err = service.AgolaApi.DeleteOrganization(organization, userCreator)
		if err != nil {
			log.Println("error in agola DeleteOrganization:", err)
			InternalServerError(w)
			return
		}
	}

	err = service.GitGateway.DeleteWebHook(gitSource, userCreator, organization.GitPath, organization.WebHookID)
	if err != nil {
		log.Println("DeleteWebHook error:", err)
		InternalServerError(w)
		return
	}

	err = service.Db.DeleteOrganization(organization.AgolaOrganizationRef)
	if err != nil {
		log.Println("db DeleteOrganization error:", err)
		InternalServerError(w)
		return
	}

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
	locked = false

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
// @Success 200 "ok"
// @Failure 404 "not found"
// @Router /addexternaluser/{organizationRef} [post]
// @Security ApiKeyToken
func (service *OrganizationService) AddExternalUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

	userId := r.Context().Value(controller.UserIdParameter).(uint64)
	user, _ := service.Db.GetUserByUserId(userId)
	if user == nil {
		log.Println("User", userId, "not found")
		InternalServerError(w)
		return
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
		log.Println("gitSource not found err:", err)
	}

	isOwner, err := service.GitGateway.IsUserOwner(gitSource, user, organization.GitPath)
	if err != nil {
		log.Println("Error from IsUserOwner:", err)
		InternalServerError(w)
		return
	}
	if !isOwner {
		log.Println("User", user.UserID, "is not owner")

		response := dto.OrganizationResponseDto{ErrorCode: dto.UserNotOwnerError}
		JSONokResponse(w, response)
		return
	}

	var req *dto.ExternalUserDto
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("decode error:", err)
	}

	if organization.ExternalUsers == nil {
		organization.ExternalUsers = make(map[string]bool)
	}

	organization.ExternalUsers[req.Email] = true
	err = service.Db.SaveOrganization(organization)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
	locked = false

	if err != nil {
		log.Println("SaveOrganization error:", err)
		InternalServerError(w)
	}
}

// @Summary Delete External User
// @Description Delete an external user
// @Tags Organization
// @Produce  json
// @Param organizationRef path string true "Organization name"
// @Param email body string true "external user email"
// @Success 200 "ok"
// @Failure 404 "not found"
// @Router /deleteexternaluser/{organizationRef} [delete]
// @Security ApiKeyToken
func (service *OrganizationService) RemoveExternalUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

	userId := r.Context().Value(controller.UserIdParameter).(uint64)
	user, _ := service.Db.GetUserByUserId(userId)
	if user == nil {
		log.Println("User", userId, "not found")
		InternalServerError(w)
		return
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
		log.Println("gitSource not found err:", err)
	}

	isOwner, err := service.GitGateway.IsUserOwner(gitSource, user, organization.GitPath)
	if err != nil {
		log.Println("Error from IsUserOwner:", err)
		InternalServerError(w)
		return
	}
	if !isOwner {
		log.Println("User", user.UserID, "is not owner")

		response := dto.OrganizationResponseDto{ErrorCode: dto.UserNotOwnerError}
		JSONokResponse(w, response)
		return
	}

	var req *dto.ExternalUserDto
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("parsing error:", err)
		InternalServerError(w)
		return
	}

	delete(organization.ExternalUsers, req.Email)
	err = service.Db.SaveOrganization(organization)

	mutex.Unlock()
	utils.ReleaseOrganizationMutex(organizationRef, service.CommonMutex)
	locked = false

	if err != nil {
		log.Println("SaveOrganization error:", err)
		InternalServerError(w)
	}
}

// @Summary Get Report
// @Description Obtain a full report of all organizations
// @Tags Organization
// @Produce  json
// @Success 200 {array} dto.OrganizationDto "ok"
// @Router /report [get]
// @Security ApiKeyToken
func (service *OrganizationService) GetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userId, _ := r.Context().Value(controller.UserIdParameter).(uint64)
	user, _ := service.Db.GetUserByUserId(userId)
	if user == nil {
		log.Println("User", userId, "not found")
		InternalServerError(w)
		return
	}

	gitsource, _ := service.Db.GetGitSourceByName(user.GitSourceName)
	if gitsource == nil {
		log.Println("gitSource", user.GitSourceName, "not found")
		InternalServerError(w)
		return
	}

	organizations, _ := service.Db.GetOrganizations()

	retVal := make([]dto.OrganizationDto, 0)
	for _, organization := range *organizations {
		if strings.Compare(organization.GitSourceName, user.GitSourceName) == 0 {
			retVal = append(retVal, manager.GetOrganizationDto(user, &organization, gitsource, service.GitGateway))
		}
	}

	JSONokResponse(w, retVal)
}

// @Summary Get Report from a specific organization
// @Description Obtain a report of a specific organization
// @Tags Organization
// @Produce  json
// @Param organizationRef path string true "Organization Name"
// @Success 200 {object} dto.OrganizationDto "ok"
// @Failure 404 "not found"
// @Router /report/{organizationRef} [get]
// @Security ApiKeyToken
func (service *OrganizationService) GetOrganizationReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]

	userId, _ := r.Context().Value(controller.UserIdParameter).(uint64)
	user, _ := service.Db.GetUserByUserId(userId)
	if user == nil {
		log.Println("User", userId, "not found")
		InternalServerError(w)
		return
	}

	gitsource, _ := service.Db.GetGitSourceByName(user.GitSourceName)
	if gitsource == nil {
		log.Println("gitSource", user.GitSourceName, "not found")
		InternalServerError(w)
		return
	}

	organization, _ := service.Db.GetOrganizationByAgolaRef(organizationRef)
	if organization == nil {
		NotFoundResponse(w)
		return
	}

	if strings.Compare(organization.GitSourceName, user.GitSourceName) != 0 {
		log.Println("user not authorized to get report of organizarion", organizationRef)

		InternalServerError(w)
		return
	}

	JSONokResponse(w, manager.GetOrganizationDto(user, organization, gitsource, service.GitGateway))
}

// @Summary Get Report from a specific organization/project
// @Description Obtain a report of a specific organization/project
// @Tags Organization
// @Produce  json
// @Param organizationRef path string true "Organization Name"
// @Param projectName path string true "Project Name"
// @Success 200 {object} dto.ProjectDto "ok"
// @Failure 404 "not found"
// @Router /report/{organizationRef}/{projectName} [get]
// @Security ApiKeyToken
func (service *OrganizationService) GetProjectReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	organizationRef := vars["organizationRef"]
	projectName := vars["projectName"]

	userId, _ := r.Context().Value(controller.UserIdParameter).(uint64)
	user, _ := service.Db.GetUserByUserId(userId)
	if user == nil {
		log.Println("User", userId, "not found")
		InternalServerError(w)
		return
	}

	gitsource, _ := service.Db.GetGitSourceByName(user.GitSourceName)
	if gitsource == nil {
		log.Println("gitSource", user.GitSourceName, "not found")
		InternalServerError(w)
		return
	}

	organization, _ := service.Db.GetOrganizationByAgolaRef(organizationRef)
	if organization == nil {
		NotFoundResponse(w)
		return
	}

	if strings.Compare(organization.GitSourceName, user.GitSourceName) != 0 {
		log.Println("user not authorized to get report of organizarion", organizationRef)

		InternalServerError(w)
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
// @Security ApiKeyToken
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
