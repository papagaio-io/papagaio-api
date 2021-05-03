package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/service"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_agola"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func TestGetOrganizationsOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	organizationsMock := MakeOrganizationList()

	db := mock_repository.NewMockDatabase(ctl)
	db.EXPECT().GetOrganizations().Return(organizationsMock, nil)

	serviceOrganization := service.OrganizationService{
		Db: db,
	}

	ts := httptest.NewServer(http.HandlerFunc(serviceOrganization.GetOrganizations))
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL)

	assert.Equal(t, err, nil)

	var organizations *[]model.Organization
	parseBody(resp, organizations)
}

func TestAddExternalUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	commonMutex := utils.NewEventMutex()
	db := mock_repository.NewMockDatabase(ctl)
	agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)

	serviceOrganization := service.OrganizationService{
		Db:          db,
		AgolaApi:    agolaApi,
		CommonMutex: &commonMutex,
	}
	user := model.User{Email: "user@email.com"}

	org := (*MakeOrganizationList())[0]

	db.EXPECT().GetOrganizationByAgolaRef(org.AgolaOrganizationRef).Return(&org, nil)
	db.EXPECT().SaveOrganization(gomock.Any()).Return(nil)
	router := mux.NewRouter()

	router.HandleFunc("/{organizationName}", serviceOrganization.AddExternalUser)
	ts := httptest.NewServer(router)

	client := ts.Client()

	data, _ := json.Marshal(user)
	requestBody := strings.NewReader(string(data))
	resp, err := client.Post(ts.URL+"/"+org.AgolaOrganizationRef, "application/json", requestBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")

}
