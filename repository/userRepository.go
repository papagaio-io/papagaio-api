package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

func (db *AppDb) GetUsersID() ([]uint64, error) {
	var retVal []uint64 = make([]uint64, 0)

	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("user/")
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			key := string(item.Key())
			userID, _ := strconv.ParseUint(strings.Split(key, "/")[1], 10, 64)
			retVal = append(retVal, userID)
		}
		return nil
	})

	return retVal, err
}

func (db *AppDb) GetUsersIDByGitSourceName(gitSourceName string) ([]uint64, error) {
	var retVal []uint64 = make([]uint64, 0)

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
			err := json.Unmarshal(dst, &localUser)
			if err != nil {
				return err
			}

			if strings.Compare(localUser.GitSourceName, gitSourceName) != 0 {
				continue
			}

			key := string(item.Key())
			userID, _ := strconv.ParseUint(strings.Split(key, "/")[1], 10, 64)
			retVal = append(retVal, userID)
		}
		return nil
	})

	return retVal, err
}

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
			err := json.Unmarshal(dst, &localUser)
			if err != nil {
				return err
			}

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
			err := json.Unmarshal(dst, &localUser)
			if err != nil {
				return err
			}

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

func (db *AppDb) SaveUser(user *model.User) error {
	if user.UserID == nil {
		seq, err := db.DB.GetSequence([]byte("sequence/user"), 100000)
		if err != nil {
			return err
		}

		id, _ := seq.Next()
		if id == 0 {
			id, _ = seq.Next()
		}
		err = seq.Release()
		if err != nil {
			return err
		}
		user.UserID = &id
	}

	key := "user/" + fmt.Sprint(*user.UserID)
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

func (db *AppDb) DeleteUser(userId uint64) error {
	return db.DB.DropPrefix([]byte("user/" + fmt.Sprint(userId)))
}
