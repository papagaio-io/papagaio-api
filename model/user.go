package model

import "time"

type User struct {
	UserID        *uint  `json:"userId"`
	GitSourceName string `json:"gitSourceName"`

	ID      uint64 `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`

	Oauth2AccessToken          string    `json:"oauth2_access_token"`
	Oauth2RefreshToken         string    `json:"oauth2_refresh_token"`
	Oauth2AccessTokenExpiresAt time.Time `json:"oauth_2_access_token_expires_at"`

	AgolaUserRef   *string `json:"agolaUserRef"`
	AgolaTokenName *string `json:"agolaTokenName"`
	AgolaToken     *string `json:"agolaToken"`
}
