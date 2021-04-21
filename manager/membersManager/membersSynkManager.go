package membersManager

import (
	"errors"
	"log"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func SynkMembers(org *model.Organization, gitSource *model.GitSource, agolaApi agolaApi.AgolaApiInterface, gitGateway *git.GitGateway) error {
	if gitSource != nil {
		if gitSource.GitType == model.Gitea {
			SyncMembersForGitea(org, gitSource, agolaApi, gitGateway)
		} else {
			SyncMembersForGithub(org, gitSource, agolaApi, gitGateway)
		}
	} else {
		log.Println("Warning!!! Found gitSource null: ", org.ID)
		return errors.New("gitsource empty")
	}

	return nil
}
