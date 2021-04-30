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
}
