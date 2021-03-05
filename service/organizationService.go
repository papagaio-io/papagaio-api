package service

import (
	"encoding/json"
	"net/http"

	agolaApi "wecode.sorint.it/opensource/papagaio-be/api/agola"
	gitApi "wecode.sorint.it/opensource/papagaio-be/api/git"
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

	emailUserLogged := "pippo@sorint.it" //ONLY FOR TEST

	user, err := service.Db.GetUserByEmail(emailUserLogged)
	if user == nil || err != nil {
		UnprocessableEntityResponse(w, "User not authorized!")
		return
	}

	var org *model.Organization
	json.NewDecoder(r.Body).Decode(&org)

	//Some checks
	gitSource, err := service.Db.GetGitSourceByID(org.GitSourceID)
	if gitSource == nil || err != nil {
		UnprocessableEntityResponse(w, "Gitsource non found")
		return
	}

	gitOrgExists := gitApi.CheckOrganizationExists(gitSource, org.GitOrgRef)
	if gitOrgExists == false {
		UnprocessableEntityResponse(w, "Organization not found")
	}

	agolaOrg, err := service.Db.GetOrganizationByName(org.GitOrgRef)
	if agolaOrg != nil {
		UnprocessableEntityResponse(w, "Organization just present in Agola")
		return
	}

	if !contains(user.AgolaUsersRef, org.AgolaUserRefOwner) {
		UnprocessableEntityResponse(w, "AgolaUserRef not valid for user "+emailUserLogged)
		return
	}
	org.UserEmailOwner = emailUserLogged

	org.ID, err = agolaApi.CreateOrganization(org.GitOrgRef, org.Visibility)
	agolaApi.AddOrganizationMember(org.GitOrgRef, org.AgolaUserRefOwner, "owner")
	org.WebHookID, err = gitApi.CreateWebHook(gitSource, org.GitOrgRef, "*")

	//Save to db

	err = service.Db.SaveOrganization(org)
	if err != nil {
		InternalServerError(w)
		return
	}

	JSONokResponse(w, org)
}

//TODO
func (service *OrganizationService) GetGitOrganizations(w http.ResponseWriter, r *http.Request) {

}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
