package github

import "strings"

type GitHubUser struct {
	ID       int
	Username string
	Role     string
	Email    string
}

func (user *GitHubUser) HasOwnerPermission() bool {
	return strings.Compare(user.Role, "owner") == 0
}
