package membersManager

import (
	"errors"
	"log"

	"wecode.sorint.it/opensource/papagaio-api/model"
)

func SynkMembers(org *model.Organization, gitSource *model.GitSource) error {
	if gitSource != nil {
		if gitSource.GitType == model.Gitea {
			SyncMembersForGitea(org, gitSource)
		} else {
			SyncMembersForGithub(org, gitSource)
		}
	} else {
		log.Println("Warning!!! Found gitSource null: ", org.ID)
		return errors.New("gitsource empty")
	}

	return nil
}
