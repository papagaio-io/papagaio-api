package service

import (
	"net/http"
	"net/http/httptest"
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
	chanOrganizationSynk := make(chan string)
	chanDiscoveryRunFails := make(chan string)
	chanUserSynk := make(chan string)
	rtDto := triggerDto.TriggersRunTimeDto{}

	serviceTrigger = TriggersService{
		Db:                    db,
		Tr:                    tr,
		ChanOrganizationSynk:  chanOrganizationSynk,
		ChanDiscoveryRunFails: chanDiscoveryRunFails,
		ChanUserSynk:          chanUserSynk,
		RtDto:                 &rtDto,
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
