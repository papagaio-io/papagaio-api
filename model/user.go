package model

import "time"

type User struct {
	UserID        *uint64 `json:"userId"`
	GitSourceName string  `json:"gitSourceName"`
	IsAdmin       bool    `json:"isAdmin"`

	ID         uint64 `json:"id"`
	Email      string `json:"email"`
	Login      string `json:"login"`
	IsGitAdmin bool   `json:"isGitAdmin"`

	Oauth2AccessToken          string    `json:"oauth2_access_token"`
	Oauth2RefreshToken         string    `json:"oauth2_refresh_token"`
	Oauth2AccessTokenExpiresAt time.Time `json:"oauth_2_access_token_expires_at"`

	AgolaUserRef   *string `json:"agolaUserRef"`
	AgolaTokenName *string `json:"agolaTokenName"`
	AgolaToken     *string `json:"agolaToken"`
}
