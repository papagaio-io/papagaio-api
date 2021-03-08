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
	GetOrganizationByID(organizationID string) (*model.Organization, error)
	DeleteOrganization(organizationID string) error

	GetGitSources() (*[]model.GitSource, error)
	SaveGitSource(gitSource *model.GitSource) error
	GetGitSourceByName(name string) (*model.GitSource, error)
	DeleteGitSource(id string) error

	SaveUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	DeleteUser(email string) error
}

type AppDb struct {
	DB *badger.DB
}

func NewAppDb(config config.Configuration) AppDb {
	db := AppDb{}
	db.Init(config)

	//databaseDataTest(&db) //TODO remove only for test

	return db
}

func (db *AppDb) Init(config config.Configuration) {
	var err error
	db.DB, err = badger.Open(badger.DefaultOptions("/badger/papagaio-be").WithSyncWrites(true).WithTruncate(true))
	if err != nil {
		log.Fatal(err)
	}
}

func databaseDataTest(db *AppDb) {
	db.SaveOrganization(&model.Organization{ID: "123", Name: "Sorint", UserEmailOwner: "Ale"})
	db.SaveOrganization(&model.Organization{ID: "abc", Name: "SorintDeb", UserEmailOwner: "Simone"})
	db.SaveOrganization(&model.Organization{ID: "ddd", Name: "UatProjects", UserEmailOwner: "Usernameexample"})

	organizations, err := db.GetOrganizations()
	if err != nil {
		fmt.Println("GetOrganizations error:", err)
	} else {
		for _, o := range *organizations {
			fmt.Println("organization :", o)
		}
	}

	myOrg, _ := db.GetOrganizationByName("Sorint")
	if myOrg != nil {
		fmt.Println("myOrg name:", myOrg)
	}

	//////////

	db.SaveGitSource(&model.GitSource{Name: "Test1"})
	db.SaveGitSource(&model.GitSource{Name: "Test2"})
	db.SaveGitSource(&model.GitSource{Name: "Test3"})

	gs, _ := db.GetGitSources()
	for _, g := range *gs {
		fmt.Println("gitSource :", g)
	}

	mygs, _ := db.GetGitSourceByName("Test1")
	fmt.Println("mygs: ", mygs)
}
