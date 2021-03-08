package controller

import (
	"net/http"
)

type OrganizationController interface {
	GetOrganizations(w http.ResponseWriter, r *http.Request)
	CreateOrganization(w http.ResponseWriter, r *http.Request)
	GetRemoteSources(w http.ResponseWriter, r *http.Request)
}

type WebHookController interface {
	WebHookOrganization(w http.ResponseWriter, r *http.Request)
}

type MemberController interface {
	AddOrganizationMember(w http.ResponseWriter, r *http.Request)
	RemoveOrganizationMember(w http.ResponseWriter, r *http.Request)
}

type GitSourceController interface {
	GetGitSources(w http.ResponseWriter, r *http.Request)
	AddGitSource(w http.ResponseWriter, r *http.Request)
	RemoveGitSource(w http.ResponseWriter, r *http.Request)
	UpdateGitSource(w http.ResponseWriter, r *http.Request)
}
