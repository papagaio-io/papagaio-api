package repository

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
	"wecode.sorint.it/opensource/papagaio-be/config"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

type Database interface {
	GetOrganizations() (*[]model.Organization, error)
	SaveOrganization(organization *model.Organization) error
	GetOrganizationByName(organizationName string) (*model.Organization, error)
}

type AppDb struct {
	DB *badger.DB
}

func NewAppDb(config config.Configuration) AppDb {
	db := AppDb{}
	db.Init(config)

	databaseDataTest(&db) //TODO remove only for test

	return db
}

func (AppDb *AppDb) Init(config config.Configuration) {
	var err error
	AppDb.DB, err = badger.Open(badger.DefaultOptions("/badger/papagaio-be").WithSyncWrites(true).WithTruncate(true))
	if err != nil {
		log.Fatal(err)
	}
}

func databaseDataTest(db *AppDb) {
	db.SaveOrganization(&model.Organization{OrganizationName: "ORG/ALE/Sorint", OrganizationType: "gitea", OrganizationURL: "www.wecode.it"})
	db.SaveOrganization(&model.Organization{OrganizationName: "ORG/SIMONE/SorintDeb", OrganizationType: "gitea", OrganizationURL: "www.wecode.it"})
	db.SaveOrganization(&model.Organization{OrganizationName: "UatProjects", OrganizationType: "gitea", OrganizationURL: "www.wecode.it"})

	organizations, err := db.GetOrganizations()
	if err != nil {
		fmt.Println("GetOrganizations error:", err)
	} else {
		for _, o := range *organizations {
			fmt.Println("organization :", o.OrganizationName, o.OrganizationURL, o.OrganizationType)
		}
	}

	myOrg, _ := db.GetOrganizationByName("ORG/ALE/Sorint")
	if myOrg != nil {
		fmt.Println("myOrg name:", myOrg.OrganizationURL)
	}
}
