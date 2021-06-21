package membersManager

import (
	"errors"
	"log"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

func SynkMembers(org *model.Organization, gitSource *model.GitSource, agolaApi agolaApi.AgolaApiInterface, gitGateway *git.GitGateway, user *model.User) error {
	log.Println("SynkMembers", org.AgolaOrganizationRef, org.Name, "start")

	if gitSource != nil {
		if gitSource.GitType == types.Gitea {
			SyncMembersForGitea(org, gitSource, agolaApi, gitGateway, user)
		} else {
			SyncMembersForGithub(org, gitSource, agolaApi, gitGateway, user)
		}
	} else {
		log.Println("Warning!!! Found gitSource null: ", org.AgolaOrganizationRef)
		return errors.New("gitsource empty")
	}

	log.Println("SynkMembers", org.AgolaOrganizationRef, "end")

	return nil
}
