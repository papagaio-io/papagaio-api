package manager

import (
	"log"

	"wecode.sorint.it/opensource/papagaio-api/manager/membersManager"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

func StartOrganizationCheckout(db repository.Database, organization *model.Organization, gitSource *model.GitSource) {
	go organizationCheckout(db, organization, gitSource)
}

func organizationCheckout(db repository.Database, organization *model.Organization, gitSource *model.GitSource) error {
	log.Println("Start organization synk")

	membersManager.SynkMembers(organization, gitSource)

	err := repositoryManager.CheckoutAllGitRepository(db, organization, gitSource)
	if err != nil {
		return err
	}

	return nil
}
