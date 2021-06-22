package service

import (
	"encoding/json"
	"net/http"

	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

type TriggersService struct {
	Db repository.Database
	Tr utils.ConfigUtils
}

// @Summary Return time triggers
// @Description Get trigger timers
// @Tags Triggers
// @Produce  json
// @Success 200 {array} dto.ConfigTriggersDto "ok"
// @Router /gettriggersconfig [get]
// @Security OAuth2Password
func (service *TriggersService) GetTriggersConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dto := dto.ConfigTriggersDto{}
	dto.OrganizationsTriggerTime = service.Tr.GetOrganizationsTriggerTime()
	dto.RunFailedTriggerTime = service.Tr.GetRunFailedTriggerTime()
	dto.UsersTriggerTime = service.Tr.GetUsersTriggerTime()
	JSONokResponse(w, dto)
}

// @Summary Save time triggers
// @Description Save trigger timers
// @Tags Triggers
// @Produce  json
// @Success 200 {array} dto.ConfigTriggersDto "ok"
// @Router /savetriggersconfig [post]
// @Security OAuth2Password
func (service *TriggersService) SaveTriggersConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var req *dto.ConfigTriggersDto
	json.NewDecoder(r.Body).Decode(&req)
	if req.OrganizationsTriggerTime != 0 {
		service.Db.SaveOrganizationsTriggerTime(int(req.OrganizationsTriggerTime))
	}
	if req.RunFailedTriggerTime != 0 {
		service.Db.SaveRunFailedTriggerTime(int(req.RunFailedTriggerTime))
	}
	if req.UsersTriggerTime != 0 {
		service.Db.SaveUsersTriggerTime(int(req.UsersTriggerTime))
	}
}
