package controller

import (
	"net/http"
)

type OrganizationController interface {
	GetOrganizations(w http.ResponseWriter, r *http.Request)
	CreateOrganization(w http.ResponseWriter, r *http.Request)
}

type WebHookController interface {
	WebHookOrganization(w http.ResponseWriter, r *http.Request)
}

type MemberController interface {
	AddMember(w http.ResponseWriter, r *http.Request)
	RemoveMember(w http.ResponseWriter, r *http.Request)
}
