package repository

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/dgraph-io/badger"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func (db *AppDb) GetGitSources() (*[]model.GitSource, error) {
	var retVal []model.GitSource = make([]model.GitSource, 0)

	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("gs/")
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			dst := make([]byte, 0)
			value, err := item.ValueCopy(dst)

			var gitSource model.GitSource
			json.Unmarshal(value, &gitSource)

			if err != nil {
				log.Println("GetGitSources read value error:", err)
				return err
			}

			retVal = append(retVal, gitSource)
		}
		return nil
	})

	return &retVal, err
}

func (db *AppDb) SaveGitSource(gitSource *model.GitSource) error {
	gitSource.ID = getNewUid()

	key := "gs/" + string(gitSource.ID)
	value, err := json.Marshal(gitSource)
	if err != nil {
		log.Println("SaveGitSource error in json marshal", err)
		return err
	}

	err = db.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value)
		err := txn.SetEntry(e)

		return err
	})

	return err
}

func (db *AppDb) GetGitSourceById(id string) (*model.GitSource, error) {
	var gitSource model.GitSource

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		prefix := "gs/" + string(id)

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			if strings.Compare(key, prefix) != 0 {
				continue
			}

			var err error
			dst, err = item.ValueCopy(dst)
			if err != nil {
				log.Println("repository error:", err)
				return err
			}

			json.Unmarshal(dst, &gitSource)

			break
		}

		return nil
	})

	return &gitSource, err
}

func (db *AppDb) DeleteGitSource(id string) error {
	return db.DB.DropPrefix([]byte("gs/" + id))
}
