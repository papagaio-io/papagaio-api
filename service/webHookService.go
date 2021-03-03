package service

import (
	"net/http"

	"wecode.sorint.it/opensource/papagaio-be/repository"
)

type WebHookService struct {
	Db repository.Database
}

func (service *WebHookService) WebHookOrganization(w http.ResponseWriter, r *http.Request) {

}
