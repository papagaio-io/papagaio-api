package dto

type WebHookDto struct {
	Secret     string        `json:"secret"`
	Action     string        `json:"action"`
	Repository RepositoryDto `json:"repository"`
}

type RepositoryDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
