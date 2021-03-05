package service

import (
	"fmt"
	"net/http"

	"wecode.sorint.it/opensource/papagaio-be/repository"
)

type WebHookService struct {
	Db repository.Database
}

//TODO
func (service *WebHookService) WebHookOrganization(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebHookOrganization start...")
}
