package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"wecode.sorint.it/opensource/papagaio-be/config"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

type Database interface {
	GetOrganizations() []model.Organization
}

type AppDb struct {
	DB *mongo.Client
}

func NewAppDb(config config.Configuration) AppDb {
	db := AppDb{}
	db.Init(config)
	return db
}

func (AppDb *AppDb) Init(config config.Configuration) {
	/*badgerDb, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	////////

	dbConnectionString := "/tmp/badger"

	var err error
	AppDb.DB, err = gorm.Open(config.Database.DbType, dbConnectionString)

	if err != nil {
		log.Println("error db connection")
		log.Fatal(err)
	}*/
}
