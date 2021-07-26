package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type UserService struct {
	Db repository.Database
}

func (service *UserService) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestDto dto.ChangeUserRoleRequestDto
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &requestDto)
	if err != nil {
		log.Println("unmarshal error:", err)
		InternalServerError(w)
		return
	}

	err = requestDto.IsValid()
	if err != nil {
		log.Println(err)
		InternalServerError(w)
		return
	}

	user, _ := service.Db.GetUserByUserId(*requestDto.UserID)
	if user == nil {
		log.Println("user", requestDto.UserID, "not found")
		InternalServerError(w)
		return
	}

	user.IsAdmin = requestDto.UserRole == dto.Administrator

	err = service.Db.SaveUser(user)
	if err != nil {
		log.Println("error in SaveUser:", err)
		InternalServerError(w)
	}
}
