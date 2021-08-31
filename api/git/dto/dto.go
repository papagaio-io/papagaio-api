package dto

import "strings"

type TeamResponseDto struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

func (team *TeamResponseDto) HasOwnerPermission() bool {
	return strings.Compare(team.Permission, "owner") == 0
}

type UserTeamResponseDto struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CommitMetadataDto struct {
	Sha     string            `json:"sha"`
	Author  map[string]string `json:"author"`
	Parents []CommitParentDto `json:"parents"`
}

func (commitMetadata *CommitMetadataDto) GetAuthorEmail() *string {
	if commitMetadata.Author != nil {
		email, ok := commitMetadata.Author["email"]
		if ok {
			return &email
		}
	}
	return nil
}

type CommitParentDto struct {
	Sha string `json:"sha"`
}

type OrganizationDto struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
	ID        int64  `json:"id"`
}

type AccessTokenRequestDto struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Code         string `json:"code"`
	RedirectURL  string `json:"redirect_uri"`
}

type UserInfoDto struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	AvatarURL   string `json:"avatar_url"`
	IsAdmin     bool   `json:"is_admin"`
	UserPageURL string `json:"user_page_url"`
}
