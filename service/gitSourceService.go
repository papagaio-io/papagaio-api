package service

import (
	"net/http"

	"wecode.sorint.it/opensource/papagaio-be/repository"
)

type GitSourceService struct {
	Db repository.Database
}

func (service *GitSourceService) GetGitSources(w http.ResponseWriter, r *http.Request) {

}

func (service *GitSourceService) AddGitSource(w http.ResponseWriter, r *http.Request) {

}

func (service *GitSourceService) RemoveGitSource(w http.ResponseWriter, r *http.Request) {

}
