package agola

import "wecode.sorint.it/opensource/papagaio-be/model"

type AgolaCreateORGDto struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Visibility model.VisibilityType `json:"visibility"`
}

type RemoteSourcesDto struct {
	Name string `json:"name"`
}

type CreateProjectRequestDto struct {
	Name             string               `json:"name"`
	ParentRef        string               `json:"parent_ref"`
	Visibility       model.VisibilityType `json:"visibility"`
	RemoteSourceName string               `json:"remote_source_name"`
	RepoPath         string               `json:"repo_path"`
}

type CreateProjectResponseDto struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Path             string               `json:"path"`
	ParentPath       string               `json:"parent_path"`
	Visibility       model.VisibilityType `json:"visibility"`
	GlobalVisibility string               `json:"global_visibility"`
}

type OrganizationMembersResponseDto struct {
	Members []MemberDto `json:"members"`
}

type MemberDto struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Role     RoleType `json:"role"`
}

type RoleType string

const (
	Owner  RoleType = "owner"
	Member RoleType = "member"
)
