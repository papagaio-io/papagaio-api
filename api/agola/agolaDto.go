package agola

type AgolaCreateORGDto struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
}

type RemoteSourcesDto struct {
	Name string `json:"name"`
}

type CreateProjectRequestDto struct {
	Name             string `json:"name"`
	ParentRef        string `json:"parent_ref"`
	Visibility       string `json:"visibility"`
	RemoteSourceName string `json:"remote_source_name"`
	RepoPath         string `json:"repo_path"`
}

type CreateProjectResponseDto struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Path             string `json:"path"`
	ParentPath       string `json:"parent_path"`
	Visibility       string `json:"visibility"`
	GlobalVisibility string `json:"global_visibility"`
}
