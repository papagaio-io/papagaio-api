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
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

var serviceUser UserService

func setupUserMock(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db = mock_repository.NewMockDatabase(ctl)

	serviceUser = UserService{
		Db: db,
	}
}

func TestChangeUserRoleOK(t *testing.T) {
	setupUserMock(t)

	user := test.MakeUser()
	request := dto.ChangeUserRoleRequestDto{
		UserID:   utils.NewUint64(*user.UserID),
		UserRole: dto.Administrator,
	}

	db.EXPECT().GetUserByUserId(*user.UserID).Return(user, nil)
	db.EXPECT().SaveUser(gomock.Any()).Return(nil)

	data, _ := json.Marshal(request)
	requestBody := strings.NewReader(string(data))

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/changeuserrole", serviceUser.ChangeUserRole)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Post(ts.URL+"/changeuserrole", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestChangeUserRoleNotFound(t *testing.T) {
	setupUserMock(t)

	user := test.MakeUser()
	request := dto.ChangeUserRoleRequestDto{
		UserID:   utils.NewUint64(*user.UserID),
		UserRole: dto.Administrator,
	}

	db.EXPECT().GetUserByUserId(*user.UserID).Return(nil, nil)

	data, _ := json.Marshal(request)
	requestBody := strings.NewReader(string(data))

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/changeuserrole", serviceUser.ChangeUserRole)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Post(ts.URL+"/changeuserrole", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")
}

func TestChangeUserRoleRequestNotValid(t *testing.T) {
	setupUserMock(t)

	user := test.MakeUser()
	request := dto.ChangeUserRoleRequestDto{
		UserID:   utils.NewUint64(*user.UserID),
		UserRole: "guest",
	}

	data, _ := json.Marshal(request)
	requestBody := strings.NewReader(string(data))

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/changeuserrole", serviceUser.ChangeUserRole)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Post(ts.URL+"/changeuserrole", "application/json", requestBody)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")
}
