package service

import (
	"encoding/json"
	"net/http"

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

	var org *model.Organization
	json.NewDecoder(r.Body).Decode(&org)

	err := service.Db.SaveOrganization(org)
	if err != nil {
		InternalServerError(w)
		return
	}

	JSONokResponse(w, org)
}

//TODO
func (service *OrganizationService) GetGitOrganizations(w http.ResponseWriter, r *http.Request) {

}
