package repository

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/dgraph-io/badger"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func (db *AppDb) GetOrganizationsName() ([]string, error) {
	var retVal []string = make([]string, 0)

	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("org/")
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			key := string(item.Key())
			retVal = append(retVal, strings.Split(key, "/")[1])

		}
		return nil
	})

	return retVal, err
}

func (db *AppDb) GetOrganizations() (*[]model.Organization, error) {
	var retVal []model.Organization = make([]model.Organization, 0)

	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("org/")
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
	key := "org/" + organization.Name
	value, err := json.Marshal(organization)
	if err != nil {
		log.Println("SaveOrganization error in json marshal", err)
		return err
	}

	err = db.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value)
		err := txn.SetEntry(e)

		return err
	})

	return err
}

func (db *AppDb) GetOrganizationById(organizationID string) (*model.Organization, error) {
	var organization *model.Organization = nil

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("org/")
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var localOrganization model.Organization
			dst, _ = item.ValueCopy(dst)
			json.Unmarshal(dst, &localOrganization)
			if strings.Compare(localOrganization.ID, organizationID) != 0 {
				continue
			}

			organization = &localOrganization

			break
		}

		return nil
	})

	return organization, err
}

func (db *AppDb) GetOrganizationByName(organizationName string) (*model.Organization, error) {
	var organization *model.Organization = nil

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("org/" + organizationName)
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var localOrganization model.Organization
			dst, _ = item.ValueCopy(dst)
			json.Unmarshal(dst, &localOrganization)
			if strings.Compare(localOrganization.Name, organizationName) != 0 {
				continue
			}

			organization = &localOrganization

			break
		}

		return nil
	})

	return organization, err
}

func (db *AppDb) DeleteOrganization(organizationName string) error {
	return db.DB.DropPrefix([]byte("org/" + organizationName))
}
