package controller

import (
	"net/http"

	"wecode.sorint.it/opensource/papagaio-be/service"
)

type ControllerInterface interface {
	GetOrganizations(w http.ResponseWriter, r *http.Request)
}

type Controller struct {
	Service service.ServiceInterface
}
