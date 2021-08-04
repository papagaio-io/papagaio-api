package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"golang.org/x/oauth2"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	gitDto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/test"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_gitea"
	"wecode.sorint.it/opensource/papagaio-api/test/mock/mock_repository"
)

var oauth2Service Oauth2Service

func setupOauth2Mock(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	db = mock_repository.NewMockDatabase(ctl)
	giteaApi = mock_gitea.NewMockGiteaInterface(ctl)

	oauth2Service = Oauth2Service{
		Db:         db,
		GitGateway: &git.GitGateway{GiteaApi: giteaApi},
		Sd: &common.TokenSigningData{
			Duration: time.Minute,
			Method:   jwt.SigningMethodHS256,
			Key:      []byte("keytest"),
		},
	}
}

func TestLoginOK(t *testing.T) {
	setupOauth2Mock(t)

	gitSource := (*test.MakeGitSourceMap())["gitea"]

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/login/{gitSourceName}", oauth2Service.Login)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/login/" + gitSource.Name)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestCallbackNewUserOK(t *testing.T) {
	setupOauth2Mock(t)

	gitSource := (*test.MakeGitSourceMap())["gitea"]
	token, _ := common.GenerateOauth2JWTToken(oauth2Service.Sd, gitSource.Name)
	code := "test"
	userInfo := gitDto.UserInfoDto{
		ID:       1,
		Login:    "test_login",
		Email:    "test_email",
		FullName: "test_fullname",
		IsAdmin:  false,
	}
	accessToken := common.Token{Token: oauth2.Token{AccessToken: "test", RefreshToken: "test", TokenType: "bearer", Expiry: time.Now()}}

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	giteaApi.EXPECT().GetOauth2AccessToken(gomock.Any(), code).Return(&accessToken, nil)
	giteaApi.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&userInfo, nil)
	db.EXPECT().GetUserByGitSourceNameAndID(gitSource.Name, uint64(userInfo.ID)).Return(nil, nil)
	db.EXPECT().SaveUser(gomock.Any()).Do(func(user *model.User) error {
		id := uint64(1)
		user.UserID = &id
		return nil
	})

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/callback", oauth2Service.Callback)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/callback?code=" + code + "&state=" + token)

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestCallbackUserJustExist(t *testing.T) {
	setupOauth2Mock(t)

	gitSource := (*test.MakeGitSourceMap())["gitea"]
	token, _ := common.GenerateOauth2JWTToken(oauth2Service.Sd, gitSource.Name)
	code := "test"
	userInfo := gitDto.UserInfoDto{
		ID:       1,
		Login:    "test_login",
		Email:    "test_email",
		FullName: "test_fullname",
		IsAdmin:  false,
	}
	accessToken := common.Token{Token: oauth2.Token{AccessToken: "test", RefreshToken: "test", TokenType: "bearer", Expiry: time.Now()}}

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	giteaApi.EXPECT().GetOauth2AccessToken(gomock.Any(), code).Return(&accessToken, nil)
	giteaApi.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&userInfo, nil)
	db.EXPECT().GetUserByGitSourceNameAndID(gitSource.Name, uint64(userInfo.ID)).Return(nil, nil)
	db.EXPECT().SaveUser(gomock.Any()).Do(func(user *model.User) error {
		id := uint64(1)
		user.UserID = &id
		return nil
	})

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/callback", oauth2Service.Callback)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	resp, err := client.Get(ts.URL + "/callback?code=" + code + "&state=" + token)

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "http StatusCode is not OK")
}

func TestCallbackWithErrors(t *testing.T) {
	setupOauth2Mock(t)

	gitSource := (*test.MakeGitSourceMap())["gitea"]
	token, _ := common.GenerateOauth2JWTToken(oauth2Service.Sd, gitSource.Name)
	code := "test"
	userInfo := gitDto.UserInfoDto{
		ID:       1,
		Login:    "test_login",
		Email:    "test_email",
		FullName: "test_fullname",
		IsAdmin:  false,
	}
	accessToken := common.Token{Token: oauth2.Token{AccessToken: "test", RefreshToken: "test", TokenType: "bearer", Expiry: time.Now()}}

	router := test.SetupBaseRouter(nil)
	router.HandleFunc("/callback", oauth2Service.Callback)
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()

	// code empty

	resp, err := client.Get(ts.URL + "/callback")

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")

	// state empty

	resp, err = client.Get(ts.URL + "/callback?code=" + code)

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")

	// ParseToken error

	resp, err = client.Get(ts.URL + "/callback?code=" + code + "&state=" + "123456")

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")

	// gitSource not found

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(nil, nil)

	resp, err = client.Get(ts.URL + "/callback?code=" + code + "&state=" + token)

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")

	// GetOauth2AccessToken error

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	giteaApi.EXPECT().GetOauth2AccessToken(gomock.Any(), code).Return(nil, errors.New("test error"))

	resp, err = client.Get(ts.URL + "/callback?code=" + code + "&state=" + token)

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")

	// GetUserInfo error

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	giteaApi.EXPECT().GetOauth2AccessToken(gomock.Any(), code).Return(&accessToken, nil)
	giteaApi.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, errors.New("test error"))

	resp, err = client.Get(ts.URL + "/callback?code=" + code + "&state=" + token)

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")

	// SaveUser error

	db.EXPECT().GetGitSourceByName(gitSource.Name).Return(&gitSource, nil)
	giteaApi.EXPECT().GetOauth2AccessToken(gomock.Any(), code).Return(&accessToken, nil)
	giteaApi.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&userInfo, nil)
	db.EXPECT().GetUserByGitSourceNameAndID(gitSource.Name, uint64(userInfo.ID)).Return(nil, nil)
	db.EXPECT().SaveUser(gomock.Any()).Return(errors.New("test error"))

	resp, err = client.Get(ts.URL + "/callback?code=" + code + "&state=" + token)

	fmt.Println("status:", resp.Status)

	assert.Equal(t, err, nil)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError, "http StatusCode is not correct")
}
