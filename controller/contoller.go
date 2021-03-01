package controller

import (
	"net/http"
)

type OrganizationController interface {
	GetOrganizations(w http.ResponseWriter, r *http.Request)
	CreateOrganization(w http.ResponseWriter, r *http.Request)
}
