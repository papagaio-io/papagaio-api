package repository

import (
	"encoding/base64"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

type Database interface {
	GetOrganizationsRef() ([]string, error)
	GetOrganizations() (*[]model.Organization, error)
	SaveOrganization(organization *model.Organization) error
	GetOrganizationByName(organizationName string) (*model.Organization, error)
	GetOrganizationByAgolaRef(organizationName string) (*model.Organization, error)
	GetOrganizationById(organizationID string) (*model.Organization, error)
	DeleteOrganization(organizationID string) error

	GetGitSources() (*[]model.GitSource, error)
	SaveGitSource(gitSource *model.GitSource) error
	GetGitSourceById(id string) (*model.GitSource, error)
	GetGitSourceByName(name string) (*model.GitSource, error)
	DeleteGitSource(id string) error

	SaveUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	DeleteUser(email string) error

	GetOrganizationsTriggerTime() int
	SaveOrganizationsTriggerTime(val int)
	GetRunFailedTriggerTime() int
	SaveRunFailedTriggerTime(val int)
}

type AppDb struct {
	DB *badger.DB
}

func NewAppDb(config config.Configuration) AppDb {
	db := AppDb{}
	db.Init(config)

	orgList, _ := db.GetOrganizationsRef()
	for _, org := range orgList {
		db.DeleteOrganization(org)
	}

	return db
}

func (db *AppDb) Init(config config.Configuration) {
	var err error

	db.DB, err = badger.Open(badger.DefaultOptions(config.Database.DbPath + "/" + config.Database.DbName).WithSyncWrites(true).WithTruncate(true))
	if err != nil {
		log.Fatal(err)
	}
}

func getNewUid() string {
	uid := uuid.New()
	base64Uid := base64.RawURLEncoding.EncodeToString([]byte(uid.String()))
	uidResult := base64Uid[8:]

	return uidResult
}
