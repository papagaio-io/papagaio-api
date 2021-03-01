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
	GetOrganization(organizationName string, userName string) (*model.Organization, error)
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
	db.SaveOrganization(&model.Organization{Name: "Sorint", UserName: "Ale", Type: "gitea", URL: "www.wecode.it"})
	db.SaveOrganization(&model.Organization{Name: "SorintDeb", UserName: "Simone", Type: "gitea", URL: "www.wecode.it"})
	db.SaveOrganization(&model.Organization{Name: "UatProjects", UserName: "Usernameexample", Type: "gitea", URL: "www.wecode.it"})

	organizations, err := db.GetOrganizations()
	if err != nil {
		fmt.Println("GetOrganizations error:", err)
	} else {
		for _, o := range *organizations {
			fmt.Println("organization :", o.Name, o.URL, o.Type)
		}
	}

	myOrg, _ := db.GetOrganization("Sorint", "ALE")
	if myOrg != nil {
		fmt.Println("myOrg name:", myOrg.URL)
	}
}
