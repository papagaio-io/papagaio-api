package repository

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/dgraph-io/badger"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func (db *AppDb) GetOrganizations() (*[]model.Organization, error) {
	var retVal []model.Organization = make([]model.Organization, 0)

	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			dst := make([]byte, 0)
			value, err := item.ValueCopy(dst)

			var organization model.Organization
			json.Unmarshal(value, &organization)

			if err != nil {
				log.Println("GetOrganizations read value error:", err)
				return err
			}

			retVal = append(retVal, organization)
		}
		return nil
	})

	return &retVal, err
}

func (db *AppDb) SaveOrganization(organization *model.Organization) error {
	key := "ORG/" + organization.ID
	value, err := json.Marshal(organization)
	if err != nil {
		log.Println("SaveOrganization erro in json marshal", err)
		return err
	}

	err = db.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value)
		err := txn.SetEntry(e)

		return err
	})

	return err
}

func (db *AppDb) GetOrganization(name string, userName string) (*model.Organization, error) {
	var organization *model.Organization

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte("ORG/" + userName + "/" + name)
		for it.Seek(prefix); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			if strings.Compare(key, name) != 0 {
				continue
			}

			var err error
			dst, err = item.ValueCopy(dst)
			if err != nil {
				log.Println("repository error:", err)
				return err
			}

			json.Unmarshal(dst, &organization)

			break
		}

		return nil
	})

	return organization, err
}
