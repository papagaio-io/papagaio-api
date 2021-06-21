package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/dgraph-io/badger"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func (db *AppDb) GetUserByUserId(userId uint64) (*model.User, error) {
	var user *model.User

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("user/")
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var localUser model.User
			dst, _ = item.ValueCopy(dst)
			json.Unmarshal(dst, &localUser)
			if *localUser.UserID != userId {
				continue
			}

			user = &localUser

			break
		}

		return nil
	})

	return user, err
}

func (db *AppDb) GetUserByGitSourceNameAndID(gitSourceName string, id uint64) (*model.User, error) {
	var user *model.User

	dst := make([]byte, 0)
	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("user/")
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			var localUser model.User
			dst, _ = item.ValueCopy(dst)
			json.Unmarshal(dst, &localUser)
			if strings.Compare(localUser.GitSourceName, gitSourceName) != 0 || localUser.ID != id {
				continue
			}

			user = &localUser

			break
		}

		return nil
	})

	return user, err
}

func (db *AppDb) SaveUser(user *model.User) (*model.User, error) {
	if user.UserID == nil {
		seq, err := db.DB.GetSequence([]byte("sequence/user"), 100000)
		if err != nil {
			return nil, err
		}

		id, _ := seq.Next()
		if id == 0 {
			id, _ = seq.Next()
		}
		err = seq.Release()
		if err != nil {
			return nil, err
		}
		user.UserID = &id
	}

	key := "user/" + fmt.Sprint(*user.UserID)
	value, err := json.Marshal(user)
	if err != nil {
		log.Println("SaveUser error in json marshal", err)
		return nil, err
	}

	err = db.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value)
		err := txn.SetEntry(e)

		return err
	})

	return user, err
}

func (db *AppDb) DeleteUser(userId uint) error {
	return db.DB.DropPrefix([]byte("user/" + fmt.Sprint(userId)))
}
