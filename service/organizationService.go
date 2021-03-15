package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	agolaApi "wecode.sorint.it/opensource/papagaio-be/api/agola"
	gitApi "wecode.sorint.it/opensource/papagaio-be/api/git"
	"wecode.sorint.it/opensource/papagaio-be/dto"
	"wecode.sorint.it/opensource/papagaio-be/manager"
	"wecode.sorint.it/opensource/papagaio-be/model"
	"wecode.sorint.it/opensource/papagaio-be/repository"
)

type OrganizationService struct {
	Db repository.Database
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

	emailUserLogged := "test@sorint.it" //ONLY FOR TEST

	user, err := service.Db.GetUserByEmail(emailUserLogged)
	fmt.Println("user: ", user)
	if user == nil || err != nil {
		UnprocessableEntityResponse(w, "User not authorized!")
		return
	}

	var req *dto.CreateOrganizationDto
	json.NewDecoder(r.Body).Decode(&req)

	if req.IsValid() != nil {
		UnprocessableEntityResponse(w, "Parameters have no correct values")
		return
	}

	log.Println("Req CreateOrganizationDto: ", req)

	org := &model.Organization{}
	org.Name = req.Name
	org.GitSourceID = req.GitSourceId
	org.Visibility = req.Visibility
	org.BehaviourType = req.BehaviourType
	org.BehaviourInclude = req.BehaviourInclude
	org.BehaviourExclude = req.BehaviourExclude

	//Some checks
	gitSource, err := service.Db.GetGitSourceById(org.GitSourceID)
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

	org.UserEmailCreator = emailUserLogged

	org.WebHookID, err = gitApi.CreateWebHook(gitSource, org.Name)
	if err != nil {
		UnprocessableEntityResponse(w, err.Error())
		return
	}

	org.ID, err = agolaApi.CreateOrganization(org.Name, org.Visibility)

	log.Println("Organization created: ", org.ID)
	log.Println("WebHook created: ", org.WebHookID)

	err = service.Db.SaveOrganization(org)
	if err != nil {
		InternalServerError(w)
		return
	}

	manager.StartSynkOrganization(service.Db, org, gitSource)

	JSONokResponse(w, org.ID)
}

func (service *OrganizationService) GetRemoteSources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	remoteSources, _ := agolaApi.GetRemoteSources()

	JSONokResponse(w, remoteSources)
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
