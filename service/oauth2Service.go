package service

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/common"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type Oauth2Service struct {
	Db         repository.Database
	Sd         *common.TokenSigningData
	GitGateway *git.GitGateway
}

func (service *Oauth2Service) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	gitSourceName := vars["gitSourceName"]
	gitSource, _ := service.Db.GetGitSourceByName(gitSourceName)

	if gitSource == nil {
		NotFoundResponse(w)
		return
	}

	token, err := common.GenerateOauth2JWTToken(service.Sd, gitSourceName)
	if err != nil {
		log.Println(err)
		InternalServerError(w)
	}

	redirectUrl := controller.GetRedirectUrl()

	response := dto.OauthLoginResponseDto{Oauth2RedirectURL: service.GitGateway.GetOauth2AuthorizePathUrl(gitSource, redirectUrl, token)}
	JSONokResponse(w, response)
}

func (service *Oauth2Service) Callback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		log.Println("code not present in query param")
		InternalServerError(w)
	}

	state := r.URL.Query().Get("state")
	if len(state) == 0 {
		log.Println("state not present in query param")
		InternalServerError(w)
	}

	token, err := common.ParseToken(service.Sd, state)
	if err != nil {
		log.Println("failed to parse jwt:", err)
		InternalServerError(w)
		return
	}
	if !token.Valid {
		log.Println("invalid token")
		InternalServerError(w)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	gitSourceName := claims["git_source_name"].(string)

	exp := int64(claims["exp"].(float64))
	//exp := claims["exp"].(int64)
	expTime := time.Unix(exp, 0)
	if common.IsAccessTokenExpired(expTime) {
		log.Println("Your token was expired at: ", expTime)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	gitSource, _ := service.Db.GetGitSourceByName(gitSourceName)
	if gitSource == nil {
		log.Println("gitSource", gitSourceName, "not found")
		InternalServerError(w)
		return
	}

	accessToken, err := service.GitGateway.GetOauth2AccessToken(gitSource, code)
	if err != nil || accessToken == nil {
		log.Println("error during access token request:", err)
		InternalServerError(w)
		return
	}

	tempUser := &model.User{
		Oauth2AccessToken:          accessToken.AccessToken,
		Oauth2AccessTokenExpiresAt: accessToken.ExpiryAt,
		Oauth2RefreshToken:         accessToken.RefreshToken,
	}

	userInfo, err := service.GitGateway.GetUserInfo(gitSource, tempUser)
	if err != nil || userInfo == nil {
		log.Println("Failed to get userinfo")
		InternalServerError(w)
		return
	}

	//Create new user if not exists and update token
	user, _ := service.Db.GetUserByGitSourceNameAndID(gitSourceName, uint64(userInfo.ID))
	if user == nil {
		user = &model.User{
			GitSourceName: gitSourceName,
			ID:            uint64(userInfo.ID),
			Email:         userInfo.Email,
			IsAdmin:       userInfo.IsAdmin,
			Login:         userInfo.Login,
		}
	}
	user.Oauth2AccessToken = accessToken.AccessToken
	user.Oauth2AccessTokenExpiresAt = accessToken.ExpiryAt
	user.Oauth2RefreshToken = accessToken.RefreshToken

	user.IsAdmin = userInfo.IsAdmin
	user.Login = userInfo.Login
	user.Email = userInfo.Email

	user, err = service.Db.SaveUser(user)

	if err != nil || user == nil {
		log.Println("Error on creating user:", err, user)
		InternalServerError(w)
		return
	}

	userToken, err := common.GenerateLoginJWTToken(service.Sd, *user.UserID)
	if err != nil {
		log.Println("Error during generate login token:", err)
		InternalServerError(w)
	}

	response := dto.OauthCallbackResponseDto{Token: userToken, UserID: *user.UserID, UserInfo: *userInfo}
	JSONokResponse(w, response)

	log.Println("Callback end for user:", *user.UserID)
}
