package membersManager

import (
	"errors"
	"log"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func SynkMembers(org *model.Organization, gitSource *model.GitSource, agolaApi agolaApi.AgolaApiInterface) error {
	if gitSource != nil {
		if gitSource.GitType == model.Gitea {
			SyncMembersForGitea(org, gitSource, agolaApi)
		} else {
			SyncMembersForGithub(org, gitSource, agolaApi)
		}
	} else {
		log.Println("Warning!!! Found gitSource null: ", org.ID)
		return errors.New("gitsource empty")
	}

	return nil
}
