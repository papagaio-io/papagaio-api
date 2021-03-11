package utils

import "strings"

func ConvertGitToAgolaUsername(gitUserName string) string {
	return strings.ReplaceAll(".", gitUserName, "")
}
