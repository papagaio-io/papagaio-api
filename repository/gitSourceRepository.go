package repository

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/dgraph-io/badger"
	"wecode.sorint.it/opensource/papagaio-api/model"
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
			if err != nil {
				return err
			}

			var gitSource model.GitSource
			err = json.Unmarshal(value, &gitSource)
			if err != nil {
				log.Println("unmarshal error:", err)
				return err
			}

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
	if len(gitSource.ID) == 0 {
		gitSource.ID = getNewUid()
	}

	key := "gs/" + string(gitSource.Name)
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

func (db *AppDb) GetGitSourceByName(name string) (*model.GitSource, error) {
	var retVal *model.GitSource

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		prefix := "gs/" + string(name)

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

			var gitSource model.GitSource
			err = json.Unmarshal(dst, &gitSource)
			if err != nil {
				log.Println("unmarshal error:", err)
				return err
			}

			retVal = &gitSource

			break
		}

		return nil
	})

	return retVal, err
}

func (db *AppDb) GetGitSourceById(id string) (*model.GitSource, error) {
	var gitSource *model.GitSource

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		prefix := "gs/"

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var localGitSource model.GitSource
			dst, _ = item.ValueCopy(dst)
			err := json.Unmarshal(dst, &localGitSource)
			if err != nil {
				log.Println("GetGitSourceById unmarshal error:", err)
				return err
			}

			if strings.Compare(localGitSource.ID, id) != 0 {
				continue
			}

			gitSource = &localGitSource

			break
		}

		return nil
	})

	return gitSource, err
}

func (db *AppDb) DeleteGitSource(id string) error {
	return db.DB.DropPrefix([]byte("gs/" + id))
}
