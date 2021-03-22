package gitea

type CreateWebHookRequestDto struct {
	Active       bool                    `json:"active"`
	BranchFilter string                  `json:"branch_filter"`
	Config       WebHookConfigRequestDto `json:"config"`
	Events       []string                `json:"events"`
	Type         string                  `json:"type"`
}

type CreateWebHookResponseDto struct {
	ID        int                      `json:"id"`
	Type      string                   `json:"type"`
	Config    WebHookConfigResponseDto `json:"config"`
	Events    []string                 `json:"events"`
	Active    bool                     `json:"active"`
	UpdatedAt string                   `json:"updated_at"`
	CreatedAt string                   `json:"created_at"`
}

type WebHookConfigRequestDto struct {
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
	HTTPMethod  string `json:"http_method"`
}

type WebHookConfigResponseDto struct {
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

type RepositoryDto struct {
	Name string `json:"name"`
}

type MetadataResponseDto struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

type BranchResponseDto struct {
	Name string `json:"name"`
}
