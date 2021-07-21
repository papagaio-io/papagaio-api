package service

import (
	"encoding/json"
	"net/http"

	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"

	triggerDto "wecode.sorint.it/opensource/papagaio-api/trigger/dto"
)

type TriggersService struct {
	Db                    repository.Database
	Tr                    utils.ConfigUtils
	ChanOrganizationSynk  chan string
	ChanDiscoveryRunFails chan string
	ChanUserSynk          chan string
	RtDto                 *triggerDto.TriggersRunTimeDto
}

const RESTART_ALL = "restartAll"
const RESTART_ORGANIZATION_SYNK_TRIGGER = "restartorganizationsynktrigger"
const RESTART_RUNS_FAILED_DISCOVERY_TRIGGER = "restartRunsFailedDiscoveryTrigger"
const RESTART_USERS_SYNK_TRIGGER = "restartUsersSynkTrigger"

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
	dto.OrganizationsTriggerLastRun = service.RtDto.OrganizationsTriggerLastRun
	dto.DiscoveryRunFailedTriggerLastRun = service.RtDto.DiscoveryRunFailedTriggerLastRun
	dto.UsersTriggerLastRun = service.RtDto.UsersTriggerLastRun

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

	var req dto.ConfigTriggersDto
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

// @Summary restart triggers
// @Description Restartr timers
// @Tags Triggers
// @Produce  json
// @Success 200
// @Router /restarttriggers [post]
// @Security OAuth2Password
func (service *TriggersService) RestartTriggers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	restartAll, err := getBoolParameter(r, RESTART_ALL)
	if err != nil {
		UnprocessableEntityResponse(w, err.Error())
		return
	}
	if restartAll {
		service.ChanDiscoveryRunFails <- "resume from TriggersService"
		service.ChanOrganizationSynk <- "resume from TriggersService"
		service.ChanUserSynk <- "resume from TriggersService"

		return
	}

	restartOrganizationSynkTrigger, err := getBoolParameter(r, RESTART_ORGANIZATION_SYNK_TRIGGER)
	if err != nil {
		UnprocessableEntityResponse(w, err.Error())
		return
	}
	if restartOrganizationSynkTrigger {
		service.ChanOrganizationSynk <- "resume from TriggersService"
	}

	restartRunsFailedDiscoveryTrigger, err := getBoolParameter(r, RESTART_RUNS_FAILED_DISCOVERY_TRIGGER)
	if err != nil {
		UnprocessableEntityResponse(w, err.Error())
		return
	}
	if restartRunsFailedDiscoveryTrigger {
		service.ChanDiscoveryRunFails <- "resume from TriggersService"
	}

	restartUsersSynkTrigger, err := getBoolParameter(r, RESTART_USERS_SYNK_TRIGGER)
	if err != nil {
		UnprocessableEntityResponse(w, err.Error())
		return
	}
	if restartUsersSynkTrigger {
		service.ChanUserSynk <- "resume from TriggersService"
	}
}
