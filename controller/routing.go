package controller

import (
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-be/config"
)

const WebHookPath string = "/webhook"
const WenHookPathParam string = "/{gitOrgRef}"

func SetupHTTPClient() {
	if config.Config.DisableSSLCertificateValidation {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func SetupRouter(router *mux.Router, ctrlOrganization OrganizationController, ctrlGitSource GitSourceController, ctrlMember MemberController, ctrlWebHook WebHookController) {
	setupPingRouter(router)
	setupGetOrganizationsRouter(router.PathPrefix("/organizations").Subrouter(), ctrlOrganization)
	setupCreateOrganizationEndpoint(router.PathPrefix("/createorganization").Subrouter(), ctrlOrganization)
	setupGetRemoteSourcesEndpoint(router.PathPrefix("/remotesources").Subrouter(), ctrlOrganization)

	setupGetGitSourcesEndpoint(router.PathPrefix("/gitsources").Subrouter(), ctrlGitSource)

	setupWebHookEndpoint(router.PathPrefix(WebHookPath).Subrouter(), ctrlWebHook)
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

func setupWebHookEndpoint(router *mux.Router, ctrl WebHookController) {
	router.HandleFunc(WenHookPathParam, ctrl.WebHookOrganization).Methods("POST")
}
