package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/common"
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
