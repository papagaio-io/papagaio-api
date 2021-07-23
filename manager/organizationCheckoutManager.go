package manager

import (
	"log"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

func StartOrganizationCheckout(db repository.Database, user *model.User, organization *model.Organization, gitSource *model.GitSource, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) {
	organizationCheckout(db, user, organization, gitSource, agolaApi, gitGateway)
}

func organizationCheckout(db repository.Database, user *model.User, organization *model.Organization, gitSource *model.GitSource, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) {
	log.Println("Start organization synk")

	err := membersManager.SynkMembers(organization, gitSource, agolaApi, gitGateway, user)
	if err != nil {
		log.Println("SynkMembers error:", err)
	}

	repositoryManager.CheckoutAllGitRepository(db, user, organization, gitSource, agolaApi, gitGateway)
}
