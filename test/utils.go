package test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func ParseBody(resp *http.Response, dto interface{}) {
	data, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal([]byte(string(data)), dto)
	if err != nil {
		log.Println("ParseBody error:", err)
	}
}

func SortProjectsDto(projects []dto.ProjectDto) []dto.ProjectDto {
	sort.SliceStable(projects[:], func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
	return projects
}

func SortBranchesDto(branches []dto.BranchDto) []dto.BranchDto {
	sort.SliceStable(branches[:], func(i, j int) bool {
		return branches[i].Name < branches[j].Name
	})
	return branches
}

func SetupBaseRouter(user *model.User) *mux.Router {
	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if user == nil {
				ctx = context.WithValue(ctx, controller.AdminUserParameter, true)
			} else {
				ctx = context.WithValue(ctx, controller.AdminUserParameter, false)
				ctx = context.WithValue(ctx, controller.UserIdParameter, user.ID)
			}

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	return router
}
