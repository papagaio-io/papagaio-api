package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
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
const gitlabDefaultApiUrl = "https://gitlab.com"

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
// @Security ApiKeyToken
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
		login := config.Config.Server.ApiExposedURL + "/api/auth/login/" + v.Name
		gs = append(gs, dto.GitSourcesDto{Name: v.Name, GitAPIURL: v.GitAPIURL, LoginURL: login, GitType: v.GitType})
	}

	JSONokResponse(w, &gs)
}

// @Summary Add a GitSource
// @Description Add a GitSource with the data provided in the body
// @Tags GitSources
// @Produce  json
// @Param gitSource body dto.CreateGitSourceRequestDto true "Git Source information"
// @Success 200 "ok"
// @Failure 422 "Already exists"
// @Router /gitsource [post]
// @Security ApiKeyToken
func (service *GitSourceService) AddGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var gitSourceDto dto.CreateGitSourceRequestDto
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &gitSourceDto)
	if err != nil {
		log.Println("unmarshal error:", err)
		InternalServerError(w)
		return
	}

	oldGitSource, _ := service.Db.GetGitSourceByName(gitSourceDto.Name)
	if oldGitSource != nil {
		UnprocessableEntityResponse(w, "Gitsource "+gitSourceDto.Name+" already exists")
		return
	}

	err = gitSourceDto.IsValid()
	if err != nil {
		log.Println("request is not valid:", err)
		UnprocessableEntityResponse(w, err.Error())
		return
	}

	if gitSourceDto.GitAPIURL == nil {
		if gitSourceDto.GitType == types.Github {
			gitUrl := githubDefaultApiUrl
			gitSourceDto.GitAPIURL = &gitUrl
		} else if gitSourceDto.GitType == types.Gitlab {
			gitUrl := gitlabDefaultApiUrl
			gitSourceDto.GitAPIURL = &gitUrl
		}
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
		if gsList != nil {
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

	err = service.Db.SaveGitSource(&gitSource)
	if err != nil {
		log.Println("error in SaveGitSource:", err)
		InternalServerError(w)
		return
	}

	JSONokResponse(w, gitSource.ID)
}

// @Summary Remove a GitSource
// @Description Remove a GitSource
// @Tags GitSources
// @Produce  json
// @Param gitSourceName path string true "Git Source Name"
// @Success 200 "ok"
// @Failure 422 "Not found"
// @Router /gitsource/{gitSourceName} [delete]
// @Security ApiKeyToken
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

	deleteRemotesource, err := getBoolParameter(r, "deleteremotesource")
	if err != nil {
		UnprocessableEntityResponse(w, err.Error())
		return
	}

	service.deleteOrganizationsAndMembersByGitsourceRef(gitSourceName)

	if deleteRemotesource {
		err := service.AgolaApi.DeleteRemotesource(gitSource.AgolaRemoteSource)
		if err != nil {
			log.Println("DeleteRemotesource error:", err)
			InternalServerError(w)
			return
		}
	}

	error := service.Db.DeleteGitSource(gitSourceName)

	if error != nil {
		InternalServerError(w)
		return
	}
}

func (service *GitSourceService) deleteOrganizationsAndMembersByGitsourceRef(gitsourceRef string) {
	orgs, _ := service.Db.GetOrganizationsByGitSource(gitsourceRef)

	if orgs != nil {
		for _, org := range *orgs {
			err := service.Db.DeleteOrganization(org.GitPath)
			if err != nil {
				log.Println("error on deleting organization", org.GitPath, ":", err)
			}
		}
	}

	users, _ := service.Db.GetUsersIDByGitSourceName(gitsourceRef)
	for _, userId := range users {
		err := service.Db.DeleteUser(userId)
		if err != nil {
			log.Println("error deleting user", userId, ":", err)
		}
	}

}

// @Summary Update a GitSource
// @Description Update GitSource information
// @Tags GitSources
// @Produce  json
// @Param gitSourceName path string true "Git Source Name"
// @Param gitSource body dto.UpdateGitSourceRequestDto true "Git Source information"
// @Success 200 "ok"
// @Failure 404 "not found"
// @Router /gitsource/{gitSourceName} [put]
// @Security ApiKeyToken
func (service *GitSourceService) UpdateGitSource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	gitSourceName := vars["gitSourceName"]

	var req dto.UpdateGitSourceRequestDto
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println("unmarshal err:", err)
		InternalServerError(w)
		return
	}

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

	err = service.Db.SaveGitSource(oldGitSource)
	if err != nil {
		log.Println("error on saving gitSource:", err)
		InternalServerError(w)
	}
}

// @Summary List Git Organizations
// @Description Return a list of all Organizations
// @Tags GitSources
// @Produce  json
// @Success 200 "ok"
// @Failure 404 "not found"
// @Router /gitorganizations [get]
// @Security ApiKeyToken
func (service *GitSourceService) GetGitOrganizations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userId, _ := r.Context().Value(controller.UserIdParameter).(uint64)
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
		InternalServerError(w)
		return
	}

	organizations, err := service.GitGateway.GetOrganizations(gitSource, user)
	if err != nil {
		log.Println("GitGateway GetOrganizations error:", err.Error())
		InternalServerError(w)
		return
	}

	if organizations != nil {
		retVal := *organizations

		sort.SliceStable(retVal, func(i, j int) bool {
			return strings.Compare(strings.ToLower(retVal[i].Path), strings.ToLower(retVal[j].Path)) < 0
		})

		JSONokResponse(w, retVal)
	}
}
