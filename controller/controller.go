package controller

import (
	"net/http"
)

type OrganizationController interface {
	GetOrganizations(w http.ResponseWriter, r *http.Request)
	CreateOrganization(w http.ResponseWriter, r *http.Request)
	DeleteOrganization(w http.ResponseWriter, r *http.Request)
	AddExternalUser(w http.ResponseWriter, r *http.Request)
	RemoveExternalUser(w http.ResponseWriter, r *http.Request)
	GetReport(w http.ResponseWriter, r *http.Request)
	GetOrganizationReport(w http.ResponseWriter, r *http.Request)
	GetProjectReport(w http.ResponseWriter, r *http.Request)
	GetAgolaOrganizations(w http.ResponseWriter, r *http.Request)
}

type WebHookController interface {
	WebHookOrganization(w http.ResponseWriter, r *http.Request)
}

type GitSourceController interface {
	GetGitSources(w http.ResponseWriter, r *http.Request)
	AddGitSource(w http.ResponseWriter, r *http.Request)
	RemoveGitSource(w http.ResponseWriter, r *http.Request)
	UpdateGitSource(w http.ResponseWriter, r *http.Request)
	GetGitOrganizations(w http.ResponseWriter, r *http.Request)
}

type TriggersController interface {
	GetTriggersConfig(w http.ResponseWriter, r *http.Request)
	SaveTriggersConfig(w http.ResponseWriter, r *http.Request)
}

type Oauth2Controller interface {
	Login(w http.ResponseWriter, r *http.Request)
	Callback(w http.ResponseWriter, r *http.Request)
}
