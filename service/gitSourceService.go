package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type GitSourceService struct {
	Db         repository.Database
	GitGateway *git.GitGateway
}

// @Summary Return a list of gitsources
// @Description Return a list of gitsources
// @Tags GitSources
// @Produce  json
// @Success 200 {object} model.GitSource "ok"
// @Failure 400 "bad request"
// @Router /gitsources [get]
// @Security OAuth2Password
func (service *GitSourceService) GetGitSources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	gitSources, err := service.Db.GetGitSources()
	if err != nil {
		InternalServerError(w)
		return
	}
	gs := make([]dto.GitSourcesDto, 0)

	for _, v := range *gitSources {
		gs = append(gs, dto.GitSourcesDto{Name: v.Name, GitAPIURL: v.GitAPIURL})
	}

	JSONokResponse(w, &gs)
}

// @Summary Add a GitSource
// @Description Add a GitSource
// @Tags GitSources
// @Produce  json
// @Success 200 {object} model.GitSource "ok"
// @Failure 400 "bad request"
// @Router /gitsource [post]
// @Security OAuth2Password
func (service *GitSourceService) AddGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var gitGitSource model.GitSource
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &gitGitSource)

	oldGitSource, _ := service.Db.GetGitSourceByName(gitGitSource.Name)
	if oldGitSource != nil {
		UnprocessableEntityResponse(w, "Gitsource "+gitGitSource.Name+" already exists")
		return
	}

	service.Db.SaveGitSource(&gitGitSource)
	JSONokResponse(w, gitGitSource.ID)
}

// @Summary Remove a GitSource
// @Description Remove a GitSource
// @Tags GitSources
// @Produce  json
// @Param gitSourceName path int true "Git Source Name"
// @Success 200 {object} model.GitSource "ok"
// @Failure 400 "bad request"
// @Router /gitsource/{gitSourceName} [delete]
// @Security OAuth2Password
func (service *GitSourceService) RemoveGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	gitSourceName := vars["gitSourceName"]
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

// @Summary Update a GitSource
// @Description Update a GitSource
// @Tags GitSources
// @Produce  json
// @Param gitSourceName path int true "Git Source Name"
// @Success 200 {object} model.GitSource "ok"
// @Failure 400 "bad request"
// @Router /gitsource/{gitSourceName} [put]
// @Security OAuth2Password
func (service *GitSourceService) UpdateGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	gitSourceName := vars["gitSourceName"]

	var req dto.UpdateRemoteSourceRequestDto
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &req)

	oldGitSource, _ := service.Db.GetGitSourceByName(gitSourceName)
	if oldGitSource == nil {
		NotFoundResponse(w)
		return
	}

	if req.AgolaRemoteSource != nil {
		oldGitSource.AgolaRemoteSource = *req.AgolaRemoteSource
	}
	if req.AgolaToken != nil {
		oldGitSource.AgolaToken = *req.AgolaToken
	}
	if req.GitAPIURL != nil {
		oldGitSource.GitAPIURL = *req.GitAPIURL
	}
	if req.GitToken != nil {
		oldGitSource.GitToken = *req.GitToken
	}
	if req.GitType != nil {
		oldGitSource.GitType = *req.GitType
	}

	service.Db.SaveGitSource(oldGitSource)
}

// @Summary List Git Organizations
// @Description Return a list of Organizations by GitSource
// @Tags GitSources
// @Produce  json
// @Param gitSourceName path int true "Git Source Name"
// @Success 200 {object} model.GitSource "ok"
// @Failure 400 "bad request"
// @Router /gitsource/{gitSourceName} [get]
// @Security OAuth2Password
func (service *GitSourceService) GetGitOrganizations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	gitSourceName := vars["gitSourceName"]
	gitSource, _ := service.Db.GetGitSourceByName(gitSourceName)

	log.Println("gitSourceName:", gitSourceName)

	if gitSource == nil {
		log.Println("gitSource", gitSourceName, "non trovato")
		NotFoundResponse(w)
		return
	}

	organizations, err := service.GitGateway.GetOrganizations(gitSource)
	if err != nil {
		log.Println("GitGateway GetOrganizations error:", err.Error())
		InternalServerError(w)
		return
	}

	if organizations != nil {
		JSONokResponse(w, organizations)
	}
}
