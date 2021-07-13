package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

const githubDefaultApiUrl = "https://api.github.com"

type GitSourceService struct {
	Db         repository.Database
	GitGateway *git.GitGateway
	AgolaApi   agolaApi.AgolaApiInterface
}

// @Summary Return a list of gitsources
// @Description Return a list of gitsources
// @Tags GitSources
// @Produce  json
// @Success 200 {array} dto.GitSourcesDto "ok"
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
		login := config.Config.Server.LocalHostAddress + "/api/auth/login/" + v.Name
		gs = append(gs, dto.GitSourcesDto{Name: v.Name, GitAPIURL: v.GitAPIURL, LoginURL: login, GitType: v.GitType})
	}

	JSONokResponse(w, &gs)
}

// @Summary Add a GitSource
// @Description Add a GitSource with the data provided in the body
// @Tags GitSources
// @Produce  json
// @Param gitSource body model.GitSource true "Git Source information"
// @Success 200 {object} model.GitSource "ok"
// @Failure 422 "Already exists"
// @Router /gitsource [post]
// @Security OAuth2Password
func (service *GitSourceService) AddGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var gitSourceDto dto.CreateGitSourceRequestDto
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &gitSourceDto)

	oldGitSource, _ := service.Db.GetGitSourceByName(gitSourceDto.Name)
	if oldGitSource != nil {
		UnprocessableEntityResponse(w, "Gitsource "+gitSourceDto.Name+" already exists")
		return
	}

	err := gitSourceDto.IsValid()
	if err != nil {
		log.Println("request is not valid:", err)
		UnprocessableEntityResponse(w, err.Error())
		return
	}

	if gitSourceDto.GitAPIURL == nil && gitSourceDto.GitType == types.Github {
		gitUrl := githubDefaultApiUrl
		gitSourceDto.GitAPIURL = &gitUrl
	}

	gitSource := model.GitSource{
		Name:        gitSourceDto.Name,
		GitType:     gitSourceDto.GitType,
		GitAPIURL:   *gitSourceDto.GitAPIURL,
		GitClientID: gitSourceDto.GitClientID,
		GitSecret:   gitSourceDto.GitClientSecret,
	}

	if gitSourceDto.AgolaRemoteSourceName == nil || len(*gitSourceDto.AgolaRemoteSourceName) == 0 {
		gsList, err := service.AgolaApi.GetRemoteSources()
		if err != nil {
			log.Println("Error in GetRemoteSources:", err)
			InternalServerError(w)
			return
		}

		findRemoteSourceName := gitSourceDto.Name
		search := true
		i := 0
		for search {
			found := false
			for _, gs := range *gsList {
				if strings.Compare(gs.Name, findRemoteSourceName) == 0 {
					found = true
					break
				}
			}
			if !found {
				search = false
			} else {
				findRemoteSourceName = gitSourceDto.Name + fmt.Sprint(i)
				i++
			}
		}

		gitSourceDto.AgolaRemoteSourceName = &findRemoteSourceName

		err = service.AgolaApi.CreateRemoteSource(*gitSourceDto.AgolaRemoteSourceName, string(gitSourceDto.GitType), *gitSourceDto.GitAPIURL, *gitSourceDto.AgolaClientID, *gitSourceDto.AgolaClientSecret)
		if err != nil {
			log.Println("Error in CreateRemoteSource:", err)
			InternalServerError(w)
			return
		}
	}

	gitSource.AgolaRemoteSource = *gitSourceDto.AgolaRemoteSourceName

	service.Db.SaveGitSource(&gitSource)
	JSONokResponse(w, gitSource.ID)
}

// @Summary Remove a GitSource
// @Description Remove a GitSource
// @Tags GitSources
// @Produce  json
// @Param gitSourceName path string true "Git Source Name"
// @Success 200 {object} model.GitSource "ok"
// @Failure 422 "Not found"
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

	deleteRemotesourceQuery, ok := r.URL.Query()["deleteremotesource"]
	deleteRemotesource := false
	if ok {
		if len(deleteRemotesourceQuery[0]) == 0 {
			deleteRemotesource = true
		} else {
			var parsError error
			deleteRemotesource, parsError = strconv.ParseBool(deleteRemotesourceQuery[0])
			if parsError != nil {
				UnprocessableEntityResponse(w, "forceCreate param value is not valid")
				return
			}
		}
	}

	service.deleteOrganizationsAndMembersByGitsourceRef(gitSourceName)

	if deleteRemotesource {
		err := service.AgolaApi.DeleteRemotesource(gitSource.AgolaRemoteSource)
		if err != nil {
			log.Println("DeleteRemotesource error:", err)
			InternalServerError(w)
		}
	}

	error := service.Db.DeleteGitSource(gitSourceName)

	if error != nil {
		UnprocessableEntityResponse(w, error.Error())
		return
	}
}

func (service *GitSourceService) deleteOrganizationsAndMembersByGitsourceRef(gitsourceRef string) {
	orgs, _ := service.Db.GetOrganizationsByGitSource(gitsourceRef)

	if orgs != nil {
		for _, org := range *orgs {
			service.Db.DeleteOrganization(org.Name)
		}
	}

	users, _ := service.Db.GetUsersIDByGitSourceName(gitsourceRef)
	if users != nil {
		for _, userId := range users {
			service.Db.DeleteUser(userId)
		}
	}

}

// @Summary Update a GitSource
// @Description Update GitSource information
// @Tags GitSources
// @Produce  json
// @Param gitSourceName path string true "Git Source Name"
// @Success 200 {object} model.GitSource "ok"
// @Failure 404 "not found"
// @Router /gitsource/{gitSourceName} [put]
// @Security OAuth2Password
func (service *GitSourceService) UpdateGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	gitSourceName := vars["gitSourceName"]

	var req dto.UpdateGitSourceRequestDto
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
	if req.GitAPIURL != nil {
		oldGitSource.GitAPIURL = *req.GitAPIURL
	}
	if req.GitType != nil {
		oldGitSource.GitType = *req.GitType
	}
	if req.GitClientID != nil {
		oldGitSource.GitClientID = *req.GitClientID
	}
	if req.GitClientSecret != nil {
		oldGitSource.GitSecret = *req.GitClientSecret
	}

	service.Db.SaveGitSource(oldGitSource)
}

// @Summary List Git Organizations
// @Description Return a list of all Organizations by GitSource
// @Tags GitSources
// @Produce  json
// @Param gitSourceName path string true "Git Source Name"
// @Success 200 {object} model.GitSource "ok"
// @Failure 404 "not found"
// @Router /gitorganizations/{gitSourceName} [get]
// @Security OAuth2Password
func (service *GitSourceService) GetGitOrganizations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userId, _ := r.Context().Value(controller.XAuthUserId).(uint64)
	user, _ := service.Db.GetUserByUserId(userId)
	if user == nil {
		log.Println("User", userId, "not found")
		InternalServerError(w)
		return
	}

	gitSource, _ := service.Db.GetGitSourceByName(user.GitSourceName)

	log.Println("gitSourceName:", user.GitSourceName)

	if gitSource == nil {
		log.Println("gitSource", user.GitSourceName, "non trovato")
		NotFoundResponse(w)
		return
	}

	organizations, err := service.GitGateway.GetOrganizations(gitSource, user)
	if err != nil {
		log.Println("GitGateway GetOrganizations error:", err.Error())
		InternalServerError(w)
		return
	}

	if organizations != nil {
		JSONokResponse(w, organizations)
	}
}
