package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	"wecode.sorint.it/opensource/papagaio-be/database"
	"wecode.sorint.it/opensource/papagaio-be/dto"
	"wecode.sorint.it/opensource/papagaio-be/service"
)

func NewRouter() http.Handler {

	router := mux.NewRouter()

	var client *mongo.Client

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, "mongodb://localhost:27017")

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("ping"))

	}).Methods("GET")

	db := database.Database{}
	srv := service.Service{Db: &db}
	ctrl := Controller{Service: &srv}
	router.HandleFunc("/organizations", ctrl.GetOrganizations).Methods("GET")

	router.HandleFunc("/organizationdbtest", CreateOrganizationEndpoint).Methods("POST")

	return router
}
func CreateOrganizationEndpoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")
	var organization dto.Organization
	_ = json.NewDecoder(r.Body).Decode(&organization)
	collection := client.Database("papagaioFirstDatabase").Collection("organization")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, organization)
	json.NewEncoder(r).Encode(result)
}
