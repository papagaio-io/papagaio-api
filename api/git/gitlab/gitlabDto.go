package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

type GitlabUser struct {
	ID          int
	Username    string
	AccessLevel gitlab.AccessLevelValue
	Email       string
}

func (user *GitlabUser) HasOwnerPermission() bool {
	return user.AccessLevel == gitlab.OwnerPermissions
}
