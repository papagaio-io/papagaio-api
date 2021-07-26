package dto

import (
	gitDto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
)

type OauthCallbackResponseDto struct {
	Token       string             `json:"token"`
	UserID      uint64             `json:"userId"`
	IsAdmin     bool               `json:"isAdmin"`
	GitUserInfo gitDto.UserInfoDto `json:"gitUserInfo"`
}
