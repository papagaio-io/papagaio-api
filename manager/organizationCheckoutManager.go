package manager

import (
	"log"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

func StartOrganizationCheckout(db repository.Database, organization *model.Organization, gitSource *model.GitSource, agolaApi agola.AgolaApiInterface) {
	go organizationCheckout(db, organization, gitSource, agolaApi)
}

func organizationCheckout(db repository.Database, organization *model.Organization, gitSource *model.GitSource, agolaApi agola.AgolaApiInterface) error {
	log.Println("Start organization synk")

	membersManager.SynkMembers(organization, gitSource, agolaApi)

	err := repositoryManager.CheckoutAllGitRepository(db, organization, gitSource, agolaApi)
	if err != nil {
		return err
	}

	return nil
}
