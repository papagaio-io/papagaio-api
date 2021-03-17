package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var gitGitSource *model.GitSource
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, gitGitSource)

	oldGitSource, _ := service.Db.GetGitSourceById(gitGitSource.ID)
	if oldGitSource != nil {
		UnprocessableEntityResponse(w, "Gitsource "+gitGitSource.Name+" already exists")
		return
	}

	service.Db.SaveGitSource(gitGitSource)
}

func (service *GitSourceService) RemoveGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	gitSourceName := vars["name"]
	gitSource, _ := service.Db.GetGitSourceByName(gitSourceName)

	if gitSource == nil {
		UnprocessableEntityResponse(w, "Gitsource "+gitSourceName+" not found")
		return
	}

	error := service.Db.DeleteGitSource(gitSource.ID)

	if error != nil {
		UnprocessableEntityResponse(w, error.Error())
		return
	}
}

func (service *GitSourceService) UpdateGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var gitGitSource *model.GitSource
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, gitGitSource)

	oldGitSource, _ := service.Db.GetGitSourceById(gitGitSource.ID)
	if oldGitSource == nil {
		NotFoundResponse(w)
		return
	}

	service.Db.SaveGitSource(gitGitSource)
}
