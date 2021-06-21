package controller

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/repository"

	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

const XAuthUserId string = "X-Auth-User-Id"

const apiPath string = "/api"
const WebHookPath string = "/webhook"
const WenHookPathParam string = "/{organizationRef}"

var db repository.Database
var sd *common.TokenSigningData

type Claim struct {
	Expiry int64  `json:"exp,omitempty"`
	Sub    string `json:"sub"`
}

func GetWebHookPath() string {
	return apiPath + WebHookPath
}

func SetupHTTPClient() {
	if config.Config.DisableSSLCertificateValidation {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func SetupRouter(signingData *common.TokenSigningData, database repository.Database, router *mux.Router, ctrlOrganization OrganizationController, ctrlGitSource GitSourceController, ctrlWebHook WebHookController, ctrlTrigger TriggersController, ctrlOauth2 Oauth2Controller) {
	db = database
	sd = signingData

	apirouter := mux.NewRouter().PathPrefix("/api").Subrouter().UseEncodedPath()
	router.PathPrefix("/api").Handler(apirouter)

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	setupPingRouter(router)

	setupGetOrganizationsRouter(apirouter.PathPrefix("/organizations").Subrouter(), ctrlOrganization) //USED FOR DEBUGGING
	setupCreateOrganizationEndpoint(apirouter.PathPrefix("/createorganization").Subrouter(), ctrlOrganization)
	setupDeleteOrganizationEndpoint(apirouter.PathPrefix("/deleteorganization").Subrouter(), ctrlOrganization)
	setupAddOrganizationExternalUserEndpoint(apirouter.PathPrefix("/addexternaluser").Subrouter(), ctrlOrganization)
	setupDeleteOrganizationExternalUserEndpoint(apirouter.PathPrefix("/deleteexternaluser").Subrouter(), ctrlOrganization)
	setupReportEndpoint(apirouter.PathPrefix("/report").Subrouter(), ctrlOrganization)
	setupOrganizationReportEndpoint(apirouter.PathPrefix("/report").Subrouter(), ctrlOrganization)
	setupProjectReportEndpoint(apirouter.PathPrefix("/report").Subrouter(), ctrlOrganization)
	setupGetAgolaRefs(apirouter.PathPrefix("/agolarefs").Subrouter(), ctrlOrganization)

	setupGetGitSourcesEndpoint(apirouter.PathPrefix("/gitsources").Subrouter(), ctrlGitSource)
	setupAddGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)
	setupUpdateGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)
	setupDeleteGitSourceEndpoint(apirouter.PathPrefix("/gitsource").Subrouter(), ctrlGitSource)
	setupGetGitOrganizations(apirouter.PathPrefix("/gitorganizations").Subrouter(), ctrlGitSource)

	setupWebHookEndpoint(apirouter.PathPrefix(WebHookPath).Subrouter(), ctrlWebHook)

	setupGetTriggersConfigEndpoint(apirouter.PathPrefix("/gettriggersconfig").Subrouter(), ctrlTrigger)
	setupSaveTriggersConfigEndpoint(apirouter.PathPrefix("/savetriggersconfig").Subrouter(), ctrlTrigger)

	setupOauth2Login(apirouter.PathPrefix("/auth/login").Subrouter(), ctrlOauth2)
	setupOauth2Callback(apirouter.PathPrefix("/auth/callback").Subrouter(), ctrlOauth2)

	//TODO SOLO PER I TEST. DA RIMUOVERE
	/*ferouter := mux.NewRouter().PathPrefix("").Subrouter().UseEncodedPath()
	router.PathPrefix("").Handler(ferouter)
	setupTestOauth2Callback(ferouter.PathPrefix("/auth/callback").Subrouter(), ctrlOauth2)*/
}

/*func setupTestOauth2Callback(router *mux.Router, ctrl Oauth2Controller) {
	router.HandleFunc("", ctrl.Callback).Methods("GET")
}*/

func setupPingRouter(router *mux.Router) {
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})
}

func setupGetOrganizationsRouter(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleRestrictedAdminRoutes)
	router.HandleFunc("", ctrl.GetOrganizations).Methods("GET")
}

func setupCreateOrganizationEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("", ctrl.CreateOrganization).Methods("POST")
}

func setupDeleteOrganizationEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("/{organizationRef}", ctrl.DeleteOrganization).Methods("DELETE")
}

func setupAddOrganizationExternalUserEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("/{organizationRef}", ctrl.AddExternalUser).Methods("POST")
}

func setupDeleteOrganizationExternalUserEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("/{organizationRef}", ctrl.RemoveExternalUser).Methods("DELETE")
}

func setupReportEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("", ctrl.GetReport).Methods("GET")
}

func setupOrganizationReportEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("/{organizationRef}", ctrl.GetOrganizationReport).Methods("GET")
}

func setupProjectReportEndpoint(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("/{organizationRef}/{projectName}", ctrl.GetProjectReport).Methods("GET")
}

func setupGetGitSourcesEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.HandleFunc("", ctrl.GetGitSources).Methods("GET")
}

func setupAddGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.Use(handleRestrictedAdminRoutes)
	router.HandleFunc("", ctrl.AddGitSource).Methods("POST")
}

func setupUpdateGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.Use(handleRestrictedAdminRoutes)
	router.HandleFunc("/{gitSourceName}", ctrl.UpdateGitSource).Methods("PUT")
}

func setupDeleteGitSourceEndpoint(router *mux.Router, ctrl GitSourceController) {
	router.Use(handleRestrictedAdminRoutes)
	router.HandleFunc("/{gitSourceName}", ctrl.RemoveGitSource).Methods("DELETE")
}

func setupGetGitOrganizations(router *mux.Router, ctrl GitSourceController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("", ctrl.GetGitOrganizations).Methods("GET")
}

func setupWebHookEndpoint(router *mux.Router, ctrl WebHookController) {
	router.HandleFunc("/{organizationRef}", ctrl.WebHookOrganization).Methods("POST")
}

func setupGetTriggersConfigEndpoint(router *mux.Router, ctrl TriggersController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("", ctrl.GetTriggersConfig).Methods("GET")
}

func setupSaveTriggersConfigEndpoint(router *mux.Router, ctrl TriggersController) {
	router.Use(handleLoggedUserRoutes)
	router.HandleFunc("", ctrl.SaveTriggersConfig).Methods("POST")
}

func setupGetAgolaRefs(router *mux.Router, ctrl OrganizationController) {
	router.Use(handleRestrictedAllRoutes)
	router.HandleFunc("", ctrl.GetAgolaOrganizations).Methods("GET")
}

func setupOauth2Login(router *mux.Router, ctrl Oauth2Controller) {
	router.HandleFunc("/{gitSourceName}", ctrl.Login).Methods("GET")
}

func setupOauth2Callback(router *mux.Router, ctrl Oauth2Controller) {
	router.HandleFunc("", ctrl.Callback).Methods("GET")
}

func handleLoggedUserRoutes(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenString == "" {
			log.Println("Undefined Authorization")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token, err := common.ParseToken(sd, tokenString)
		if err != nil {
			log.Println("failed to parse jwt:", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			log.Println("invalid token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		//userId := claims["sub"].(uint64)
		userId := uint64(claims["sub"].(float64))

		//exp := claims["exp"].(int64)
		exp := int64(claims["exp"].(float64))
		expTime := time.Unix(exp, 0)
		if common.IsAccessTokenExpired(expTime) {
			log.Println("Your token was expired at: ", expTime)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, _ := db.GetUserByUserId(userId)

		if user == nil {
			log.Println("user", userId, "not found")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		fmt.Println("http request user", userId)
		r.Header.Set(XAuthUserId, fmt.Sprint(userId))
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

			token, err := common.ParseToken(sd, tokenString)
			if err != nil {
				log.Println("failed to parse jwt:", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				log.Println("invalid token")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			//userId := claims["sub"].(uint64)
			userId := uint64(claims["sub"].(float64))

			//exp := claims["exp"].(int64)
			exp := int64(claims["exp"].(float64))
			expTime := time.Unix(exp, 0)
			if common.IsAccessTokenExpired(expTime) {
				log.Println("Your token was expired at: ", expTime)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			user, _ := db.GetUserByUserId(userId)

			if user == nil {
				log.Println("user", userId, "not found")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			fmt.Println("http request user", userId)
			r.Header.Set(XAuthUserId, fmt.Sprint(userId))
			h.ServeHTTP(w, r)
		}
	})
}

func checkIsAdminUser(authorization string) bool {
	token := strings.TrimPrefix(authorization, "token ")
	return strings.Compare(token, config.Config.AdminToken) == 0
}

const redirectPath string = "%s/auth/callback"

func GetRedirectUrl() string {
	return fmt.Sprintf(redirectPath, config.Config.Server.LocalHostAddress)
}
