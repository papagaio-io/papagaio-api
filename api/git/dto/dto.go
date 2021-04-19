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

type email string

type CommitMetadataDto struct {
	Sha     string            `json:"sha"`
	Author  map[string]string `json:"author"`
	Parents []CommitParentDto `json:"parents"`
}

func (commitMetadata *CommitMetadataDto) GetAuthorEmail() string {
	return commitMetadata.Author["email"]
}

type CommitParentDto struct {
	Sha string `json:"sha"`
}

type OrganizationDto struct {
	Name      string
	AvatarURL string
	ID        int
}
