package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-be/database"
	"wecode.sorint.it/opensource/papagaio-be/service"
)

func NewRouter() http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("ping"))

	}).Methods("GET")

	db := database.Database{}
	srv := service.Service{Db: &db}
	ctrl := Controller{Service: &srv}
	router.HandleFunc("/organization", ctrl.GetOrganizations).Methods("GET")

	return router
}
