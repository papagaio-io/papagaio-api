package controller

import (
	"encoding/json"
	"net/http"
)

func (ctrl *Controller) GetOrganizations(w http.ResponseWriter, r *http.Request) {

	organization := ctrl.Service.GetOrganizations()
	jsonResponse, _ := json.Marshal(organization)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResponse)

}
