package controller

import (
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-api/config"
)

const WebHookPath string = "/webhook"
const WenHookPathParam string = "/{gitOrgRef}"

func SetupHTTPClient() {
	if config.Config.DisableSSLCertificateValidation {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func SetupRouter(router *mux.Router, ctrlOrganization OrganizationController, ctrlGitSource GitSourceController, ctrlMember MemberController, ctrlWebHook WebHookController) {
	apirouter := mux.NewRouter().PathPrefix("/api").Subrouter().UseEncodedPath()
	router.PathPrefix("/api").Handler(apirouter)

	setupPingRouter(router)

	setupGetOrganizationsRouter(apirouter.PathPrefix("/organizations").Subrouter(), ctrlOrganization)
	setupCreateOrganizationEndpoint(apirouter.PathPrefix("/createorganization").Subrouter(), ctrlOrganization)

	setupGetGitSourcesEndpoint(apirouter.PathPrefix("/gitsources").Subrouter(), ctrlGitSource)
	setupAddGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)
	setupUpdateGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)
	setupDeleteGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)

	setupWebHookEndpoint(apirouter.PathPrefix(WebHookPath).Subrouter(), ctrlWebHook)
}

func setupPingRouter(router *mux.Router) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})
}

func setupGetOrganizationsRouter(router *mux.Router, ctrl OrganizationController) {
	router.HandleFunc("", ctrl.GetOrganizations).Methods("GET")
}

func setupCreateOrganizationEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.HandleFunc("", ctrl.CreateOrganization).Methods("POST")
}

func setupGetRemoteSourcesEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.HandleFunc("", ctrl.GetRemoteSources).Methods("GET")
}

func setupGetGitSourcesEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.HandleFunc("", ctrl.GetGitSources).Methods("GET")
}

func setupAddGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.HandleFunc("", ctrl.AddGitSource).Methods("POST")
}

func setupUpdateGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.HandleFunc("", ctrl.UpdateGitSource).Methods("PUT")
}

func setupDeleteGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.HandleFunc("/{name}", ctrl.RemoveGitSource).Methods("DELETE")
}

func setupWebHookEndpoint(router *mux.Router, ctrl WebHookController) {
	router.HandleFunc(WenHookPathParam, ctrl.WebHookOrganization).Methods("POST")
}
