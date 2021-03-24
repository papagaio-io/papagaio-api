package dto

type TeamResponseDto struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

type UserTeamResponseDto struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CommitMetadataDto struct {
	Sha    string            `json:"sha"`
	Author map[string]string `json:"author"`
}
