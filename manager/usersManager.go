package manager

import "strings"

func ConvertGitToAgolaUsername(gitUserName string) string {
	return strings.ReplaceAll(".", gitUserName, "")
}
