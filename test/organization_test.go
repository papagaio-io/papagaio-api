package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/service"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
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

	//assertEventDtoEquals(t, eventDto, dto.CreateFullEventDto(eventMock))
}

func TestCreateOrganizationOK(t *testing.T) {
	/*ctl := gomock.NewController(t)
		defer ctl.Finish()

		db := mock_repository.NewMockDatabase(ctl)
		agolaApi := mock_agola.NewMockAgolaApiInterface(ctl)
		giteaApi := mock_gitea.NewMockGiteaInterface(ctl)

		//TODO expect

		serviceOrganization := service.OrganizationService{
			Db:         db,
			AgolaApi:   agolaApi,
			GitGateway: &git.GitGateway{GiteaApi: giteaApi},
		}

		ts := httptest.NewServer(http.HandlerFunc(serviceOrganization.CreateOrganization))
		defer ts.Close()

		client := ts.Client()

		requestBody := strings.NewReader(`{"Uid":"6fcb514b-b878-4c9d-95b7-8dc3a7ce6fd8", "Slots": [2], "User":
	        {"Name": "Nuovonome", "Surname": "Nuovocognome", "Email": "nuovaemail@email.com"}, "privacyPolicy": true}}`)
		resp, err := client.Post(ts.URL, "application/json", requestBody)

		assert.Equal(t, err, nil)*/
}
