package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/test"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	triggerDto "wecode.sorint.it/opensource/papagaio-api/trigger/dto"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

var serviceTrigger TriggersService

func setupTriggerMock(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db = mock_repository.NewMockDatabase(ctl)
	tr := utils.ConfigUtils{Db: db}
	/*chanOrganizationSynk := make(chan string)
	chanDiscoveryRunFails := make(chan string)
	chanUserSynk := make(chan string)*/

	serviceTrigger = TriggersService{
		Db: db,
		Tr: tr,
		RtDtoOrganizationSynk: &triggerDto.TriggerRunTimeDto{
			Chan: make(chan triggerDto.TriggerMessage),
		},
		RtDtoDiscoveryRunFails: &triggerDto.TriggerRunTimeDto{
			Chan: make(chan triggerDto.TriggerMessage),
		},
		RtDtoUserSynk: &triggerDto.TriggerRunTimeDto{
			Chan: make(chan triggerDto.TriggerMessage),
		},
	}
}

func TestGetTriggetsConfigOK(t *testing.T) {
	setupTriggerMock(t)

	db.EXPECT().GetOrganizationsTriggerTime().Return(1)
	db.EXPECT().GetRunFailedTriggerTime().Return(2)
	db.EXPECT().GetUsersTriggerTime().Return(3)

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/gettriggersconfig", serviceTrigger.GetTriggersConfig)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/gettriggersconfig")

	var responseDto dto.ConfigTriggersDto
	test.ParseBody(resp, &responseDto)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
	assert.Equal(t, responseDto.OrganizationsTriggerTime, uint(1))
	assert.Equal(t, responseDto.RunFailedTriggerTime, uint(2))
	assert.Equal(t, responseDto.UsersTriggerTime, uint(3))
}

func TestSaveTriggetsConfigOK(t *testing.T) {
	setupTriggerMock(t)

	reqDto := dto.ConfigTriggersDto{}
	reqDto.OrganizationsTriggerTime = 4
	reqDto.RunFailedTriggerTime = 5
	reqDto.UsersTriggerTime = 6

	db.EXPECT().SaveOrganizationsTriggerTime(int(reqDto.OrganizationsTriggerTime))
	db.EXPECT().SaveRunFailedTriggerTime(int(reqDto.RunFailedTriggerTime))
	db.EXPECT().SaveUsersTriggerTime(int(reqDto.UsersTriggerTime))

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/savetriggersconfig", serviceTrigger.SaveTriggersConfig)
	ts := httptest.NewServer(router)
	defer ts.Close()

	data, _ := json.Marshal(reqDto)
	requestBody := strings.NewReader(string(data))

	client := ts.Client()
	resp, err := client.Post(ts.URL+"/savetriggersconfig", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func getChannel(c chan triggerDto.TriggerMessage) {
	<-c
}

func TestRestartTriggersOK(t *testing.T) {
	setupTriggerMock(t)

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/restarttriggers", serviceTrigger.RestartTriggers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	go getChannel(serviceTrigger.RtDtoOrganizationSynk.Chan)
	go getChannel(serviceTrigger.RtDtoDiscoveryRunFails.Chan)
	go getChannel(serviceTrigger.RtDtoUserSynk.Chan)

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/restarttriggers?restartorganizationsynktrigger&restartRunsFailedDiscoveryTrigger&restartUsersSynkTrigger")

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestRestartTriggersAllOK(t *testing.T) {
	setupTriggerMock(t)

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/restarttriggers", serviceTrigger.RestartTriggers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	go getChannel(serviceTrigger.RtDtoOrganizationSynk.Chan)
	go getChannel(serviceTrigger.RtDtoDiscoveryRunFails.Chan)
	go getChannel(serviceTrigger.RtDtoUserSynk.Chan)

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/restarttriggers?restartAll")

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}
