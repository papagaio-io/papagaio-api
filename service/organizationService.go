package service

import (
	"encoding/json"
	"net/http"

	"wecode.sorint.it/opensource/papagaio-be/repository"
)

type OrganizationService struct {
	Db repository.Database
}

func (service *OrganizationService) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	organization := service.Db.GetOrganizations()
	jsonResponse, _ := json.Marshal(organization)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResponse)
}

func (service *OrganizationService) CreateOrganizationEndpoint(w http.ResponseWriter, r *http.Request) {
	/*w.Header().Add("content-type", "application/json")
	var organization model.Organization
	_ = json.NewDecoder(r.Body).Decode(&organization)
	collection := service.Db.Database("papagaioFirstDatabase").Collection("organization")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, organization)
	json.NewEncoder(r).Encode(result)*/
}
