package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type UserService struct {
	Db repository.Database
}

func (service *UserService) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req model.User
	json.NewDecoder(r.Body).Decode(&req)

	if err := req.IsValid(); err != nil {
		UnprocessableEntityResponse(w, err.Error())
		return
	}

	user, _ := service.Db.GetUserByEmail(req.Email)
	if user == nil {
		service.Db.SaveUser(&model.User{Email: req.Email})
	} else {
		UnprocessableEntityResponse(w, "User already exists")
	}
}

func (service *UserService) RemoveUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	email := vars["email"]

	user, err := service.Db.GetUserByEmail(email)
	if err != nil {
		InternalServerError(w)
		return
	}

	if user == nil {
		NotFoundResponse(w)
		return
	}

	err = service.Db.DeleteUser(email)
	if err != nil {
		InternalServerError(w)
		return
	}
}
