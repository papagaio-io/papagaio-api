package service

import (
	"net/http"

	"wecode.sorint.it/opensource/papagaio-be/repository"
)

type GitSourceService struct {
	Db repository.Database
}

func (service *GitSourceService) GetGitSources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	gitSources, err := service.Db.GetGitSources()
	if err != nil {
		InternalServerError(w)
		return
	}

	JSONokResponse(w, gitSources)
}

func (service *GitSourceService) AddGitSource(w http.ResponseWriter, r *http.Request) {

}

func (service *GitSourceService) RemoveGitSource(w http.ResponseWriter, r *http.Request) {

}
