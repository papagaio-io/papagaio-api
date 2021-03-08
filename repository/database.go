package repository

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	"wecode.sorint.it/opensource/papagaio-be/config"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

type Database interface {
	GetOrganizations() (*[]model.Organization, error)
	SaveOrganization(organization *model.Organization) error
	GetOrganizationByName(organizationName string) (*model.Organization, error)
	GetOrganizationById(organizationID string) (*model.Organization, error)
	DeleteOrganization(organizationID string) error

	GetGitSources() (*[]model.GitSource, error)
	SaveGitSource(gitSource *model.GitSource) error
	GetGitSourceById(id string) (*model.GitSource, error)
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

	databaseDataTest(&db) //TODO remove only for test

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

	//db.SaveGitSource(&model.GitSource{Name: "gitSourceProva", GitType: "gitea", GitAPIURL: "https://wecode.sorintdev.it", GitToken: "9d600ae52773076680e5d14ba9a7ec8f6c2a5374"})
	db.SaveGitSource(&model.GitSource{Name: "gitSourceProva", GitType: "gitea", GitAPIURL: "https://wecode.sorintdev.it", GitToken: "d5e630f316de7132d4f840c305853865b2470cf2"})
	db.SaveUser(&model.User{Email: "test@sorint.it", AgolaUsersRef: []string{"test"}, AgolaUserToken: "d8fe9258aab60bb3dd192a7726cbf128747cfb0e"})
}

func getNewUid() string {
	uid := uuid.New()
	base64Uid := base64.RawURLEncoding.EncodeToString([]byte(uid.String()))
	uidResult := fmt.Sprintf("%s%s", "Stroodle", base64Uid[8:])

	return uidResult
}
