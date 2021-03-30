package controller

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2/jwt"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

const XAuthEmail string = "X-Auth-Email"

const apiPath string = "/api"
const WebHookPath string = "/webhook"
const WenHookPathParam string = "/{gitOrgRef}"

var db repository.Database

type Claim struct {
	Expiry  *jwt.NumericDate `json:"exp,omitempty"`
	Email   string           `json:"email"`
	Name    string           `json:"given_name"`
	Surname string           `json:"family_name"`
}

func GetWebHookPath() string {
	return apiPath + WebHookPath
}

func SetupHTTPClient() {
	if config.Config.DisableSSLCertificateValidation {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func SetupRouter(database repository.Database, router *mux.Router, ctrlOrganization OrganizationController, ctrlGitSource GitSourceController, ctrlMember MemberController, ctrlWebHook WebHookController, ctrlUser UserController) {
	db = database

	apirouter := mux.NewRouter().PathPrefix("/api").Subrouter().UseEncodedPath()
	router.PathPrefix("/api").Handler(apirouter)

	setupPingRouter(router)

	setupGetOrganizationsRouter(apirouter.PathPrefix("/organizations").Subrouter(), ctrlOrganization)
	setupCreateOrganizationEndpoint(apirouter.PathPrefix("/createorganization").Subrouter(), ctrlOrganization)
	setupDeleteOrganizationEndpoint(apirouter.PathPrefix("/deleteorganization").Subrouter(), ctrlOrganization)
	setupAddOrganizationExternalUserEndpoint(apirouter.PathPrefix("/addexternaluser").Subrouter(), ctrlOrganization)
	setupDeleteOrganizationExternalUserEndpoint(apirouter.PathPrefix("/deleteexternaluser").Subrouter(), ctrlOrganization)

	setupGetGitSourcesEndpoint(apirouter.PathPrefix("/gitsources").Subrouter(), ctrlGitSource)
	setupAddGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)
	setupUpdateGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)
	setupDeleteGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)

	setupWebHookEndpoint(apirouter.PathPrefix(WebHookPath).Subrouter(), ctrlWebHook)

	setupAddUserEndpoint(apirouter.PathPrefix("/adduser").Subrouter(), ctrlUser)
}

func setupPingRouter(router *mux.Router) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})
}

func setupGetOrganizationsRouter(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleRestrictedAllRoutes)
	router.HandleFunc("", ctrl.GetOrganizations).Methods("GET")
}

func setupCreateOrganizationEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleRestrictedUserRoutes)
	router.HandleFunc("", ctrl.CreateOrganization).Methods("POST")
}

func setupDeleteOrganizationEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleRestrictedUserRoutes)
	router.HandleFunc("/{id}", ctrl.DeleteOrganization).Methods("DELETE")
}

func setupAddOrganizationExternalUserEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleRestrictedUserRoutes)
	router.HandleFunc("/{id}", ctrl.DeleteOrganization).Methods("POST")
}

func setupDeleteOrganizationExternalUserEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleRestrictedUserRoutes)
	router.HandleFunc("/{id}", ctrl.DeleteOrganization).Methods("DELETE")
}

func setupGetRemoteSourcesEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleRestrictedAllRoutes)
	router.HandleFunc("", ctrl.GetRemoteSources).Methods("GET")
}

func setupGetGitSourcesEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.Use(handleRestrictedAllRoutes)
	router.HandleFunc("", ctrl.GetGitSources).Methods("GET")
}

func setupAddGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.Use(handleRestrictedAdminRoutes)
	router.HandleFunc("", ctrl.AddGitSource).Methods("POST")
}

func setupUpdateGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.Use(handleRestrictedAdminRoutes)
	router.HandleFunc("", ctrl.UpdateGitSource).Methods("PUT")
}

func setupDeleteGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.Use(handleRestrictedAdminRoutes)
	router.HandleFunc("/{name}", ctrl.RemoveGitSource).Methods("DELETE")
}

func setupWebHookEndpoint(router *mux.Router, ctrl WebHookController) {
	router.HandleFunc("/{gitOrgRef}", ctrl.WebHookOrganization).Methods("POST")
}

func setupAddUserEndpoint(router *mux.Router, ctrl UserController) {
	router.Use(handleRestrictedAdminRoutes)
	router.HandleFunc("", ctrl.AddUser).Methods("POST")
}

func handleRestrictedUserRoutes(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		if tokenString == "" {
			log.Println("Undefined Authorization")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		parsedToken, err := jwt.ParseSigned(tokenString)
		if err != nil {
			log.Println("Failed to parse the token: ", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claim := Claim{}
		err = parsedToken.Claims(config.KeycloakPubKey, &claim)
		if err != nil {
			log.Println("Failed to claim JWT token: ", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if claim.Expiry.Time().Before(time.Now().Add(time.Duration(-config.Config.Keycloak.TokenValidity) * time.Second)) {
			log.Println("Your token was expired at: ", claim.Expiry.Time())
			log.Println("The actual time is: ", time.Now())
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, err := db.GetUserByEmail(claim.Email)
		if user == nil {
			log.Println("User", claim.Email, "not found")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		r.Header.Set(XAuthEmail, claim.Email)

		h.ServeHTTP(w, r)
	})
}

func handleRestrictedAdminRoutes(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !checkIsAdminUser(r.Header.Get("Authorization")) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func handleRestrictedAllRoutes(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if checkIsAdminUser(r.Header.Get("Authorization")) {
			h.ServeHTTP(w, r)
		} else {
			tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			if tokenString == "" {
				log.Println("Undefined Authorization")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			parsedToken, err := jwt.ParseSigned(tokenString)
			if err != nil {
				log.Println("Failed to parse the token: ", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claim := Claim{}
			err = parsedToken.Claims(config.KeycloakPubKey, &claim)
			if err != nil {
				log.Println("Failed to claim JWT token: ", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if claim.Expiry.Time().Before(time.Now().Add(time.Duration(-config.Config.Keycloak.TokenValidity) * time.Second)) {
				log.Println("Your token was expired at: ", claim.Expiry.Time())
				log.Println("The actual time is: ", time.Now())
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			/*user, err := db.GetUserByEmail(claim.Email)
			if user == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}*/

			h.ServeHTTP(w, r)
		}
	})
}

func checkIsAdminUser(authorization string) bool {
	token := strings.TrimPrefix(authorization, "token ")
	return strings.Compare(token, config.Config.AdminToken) == 0
}
