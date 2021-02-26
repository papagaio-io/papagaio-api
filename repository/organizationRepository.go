package repository

import (
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

			if err != nil {
				log.Println("GetOrganizations read value error:", err)
				return err
			}

			retVal = append(retVal, model.Organization{OrganizationName: string(item.Key()), OrganizationURL: string(value), OrganizationType: "TODO"})
		}
		return nil
	})

	if err != nil {
		log.Println("repository error:", err)
	}

	return &retVal, err
}

func (db *AppDb) SaveOrganization(organization *model.Organization) error {
	err := db.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(organization.OrganizationName), []byte(organization.OrganizationURL))
		err := txn.SetEntry(e)

		if err != nil {
			log.Println("repository error:", err)
		}

		return err
	})

	return err
}

func (db *AppDb) GetOrganizationByName(organizationName string) (*model.Organization, error) {
	var organization *model.Organization = nil

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(organizationName)
		for it.Seek(prefix); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			if strings.Compare(key, organizationName) != 0 {
				continue
			}

			var err error
			dst, err = item.ValueCopy(dst)
			if err != nil {
				log.Println("repository error:", err)
				return err
			}

			organization = &model.Organization{OrganizationName: key, OrganizationURL: string(dst), OrganizationType: "TODO"}
		}

		return nil
	})

	return organization, err
}
