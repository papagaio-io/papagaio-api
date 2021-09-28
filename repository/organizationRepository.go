package repository

import (
	"encoding/json"
	"log"
	"strings"

	badger "github.com/dgraph-io/badger/v3"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func (db *AppDb) GetOrganizationsRef() ([]string, error) {
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
			if err != nil {
				return err
			}

			var organization model.Organization
			err = json.Unmarshal(value, &organization)

			if err != nil {
				return err
			}

			retVal = append(retVal, organization)
		}
		return nil
	})

	return &retVal, err
}

func (db *AppDb) SaveOrganization(organization *model.Organization) error {
	key := "org/" + organization.AgolaOrganizationRef
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
			err := json.Unmarshal(dst, &localOrganization)
			if err != nil {
				return err
			}

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

func (db *AppDb) GetOrganizationByAgolaRef(agolaOrganizationRef string) (*model.Organization, error) {
	var organization *model.Organization = nil

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("org/" + agolaOrganizationRef)
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var localOrganization model.Organization
			dst, _ = item.ValueCopy(dst)
			err := json.Unmarshal(dst, &localOrganization)
			if err != nil {
				return err
			}

			if strings.Compare(localOrganization.AgolaOrganizationRef, agolaOrganizationRef) != 0 {
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

func (db *AppDb) GetOrganizationsByGitSource(gitSourceName string) (*[]model.Organization, error) {
	organizations := make([]model.Organization, 0)

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
			err := json.Unmarshal(dst, &localOrganization)
			if err != nil {
				return err
			}

			if strings.Compare(localOrganization.GitSourceName, gitSourceName) != 0 {
				continue
			}

			organizations = append(organizations, localOrganization)
		}

		return nil
	})

	return &organizations, err
}
