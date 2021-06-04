package dto

import (
	gitDto "wecode.sorint.it/opensource/papagaio-api/api/git/dto"
)

type OauthCallbackResponseDto struct {
	Token    string             `json:"token"`
	UserID   uint               `json:"userId"`
	UserInfo gitDto.UserInfoDto `json:"userInfo"`
}
