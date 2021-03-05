package repository

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/dgraph-io/badger"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func (db *AppDb) SaveUser(user *model.User) error {
	key := "user/" + string(user.Email)
	value, err := json.Marshal(user)
	if err != nil {
		log.Println("SaveUser error in json marshal", err)
		return err
	}

	err = db.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value)
		err := txn.SetEntry(e)

		return err
	})

	return err
}

//TODO
func (db *AppDb) DeleteUser(email string) error {
	var err error
	return err
}

func (db *AppDb) GetUserByEmail(email string) (*model.User, error) {
	var user model.User

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
			if strings.Compare(localUser.Email, email) != 0 {
				continue
			}

			user = localUser

			break
		}

		return nil
	})

	return &user, err
}
