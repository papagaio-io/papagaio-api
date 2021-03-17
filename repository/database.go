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
	GetOrganizations() (*[]model.Organization, error)
	SaveOrganization(organization *model.Organization) error
	GetOrganizationByName(organizationName string) (*model.Organization, error)
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
}

type AppDb struct {
	DB *badger.DB
}

func NewAppDb(config config.Configuration) AppDb {
	db := AppDb{}
	db.Init(config)

	databaseDataInit(&db) //TODO remove only for test

	return db
}

func (db *AppDb) Init(config config.Configuration) {
	var err error

	db.DB, err = badger.Open(badger.DefaultOptions(config.Database.DbPath + "/" + config.Database.DbName).WithSyncWrites(true).WithTruncate(true))
	if err != nil {
		log.Fatal(err)
	}
}

func databaseDataInit(db *AppDb) {
	db.DB.DropAll()
	/*db.SaveGitSource(&model.GitSource{Name: "gitSourceProva", GitType: model.Gitea, GitAPIURL: "https://wecode.sorintdev.it", GitToken: "d5e630f316de7132d4f840c305853865b2470cf2", AgolaToken: "aad79c015e46597a443d9018b7517c1c4b73c2d1", AgolaRemoteSource: "gitea"})
	db.SaveUser(&model.User{Email: "test@sorint.it"})*/

	/*gitSource, _ := db.GetGitSourceById("N2ItNWUwNy00YzMyLWI0YzQtMzI3YTcwZjIwNmE4")
	gitSource.AgolaRemoteSource = "gitea"
	gitSource.AgolaToken = "aad79c015e46597a443d9018b7517c1c4b73c2d1" //token prova di tullio
	db.SaveGitSource(gitSource)*/
}

func getNewUid() string {
	uid := uuid.New()
	base64Uid := base64.RawURLEncoding.EncodeToString([]byte(uid.String()))
	uidResult := base64Uid[8:]

	return uidResult
}
