package utils

import "strings"

func ConvertGiteaToAgolaUsername(gitUserName string) string {
	return strings.ReplaceAll(gitUserName, ".", "")
}

func ConvertGithubToAgolaUsername(gitUserName string) string {
	return gitUserName
}
