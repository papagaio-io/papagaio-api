package service

import (
	"encoding/json"
	"net/http"

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
