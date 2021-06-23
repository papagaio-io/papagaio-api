package repository

import (
	"strconv"

	"github.com/dgraph-io/badger"
)

const organizationTriggerTime string = "organizationTriggerTime"
const runFailedTriggerTime string = "runFailedTriggerTime"
const usersTriggerTime string = "usersTriggerTime"

func (db *AppDb) GetOrganizationsTriggerTime() int {
	retVal := -1

	db.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(organizationTriggerTime))
		if err != nil {
			return err
		}

		dst := make([]byte, 0)
		dst, _ = item.ValueCopy(dst)
		retVal, _ = strconv.Atoi(string(dst))
		return nil
	})

	return retVal
}

func (db *AppDb) SaveOrganizationsTriggerTime(value int) {

	db.DB.Update(func(txn *badger.Txn) error {
		byteVal := []byte(strconv.Itoa(value))
		e := badger.NewEntry([]byte(organizationTriggerTime), byteVal)
		err := txn.SetEntry(e)

		return err
	})
}

func (db *AppDb) GetRunFailedTriggerTime() int {
	retVal := -1

	db.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(runFailedTriggerTime))
		if err != nil {
			return err
		}

		dst := make([]byte, 0)
		dst, _ = item.ValueCopy(dst)
		retVal, _ = strconv.Atoi(string(dst))
		return nil
	})
	return retVal
}

func (db *AppDb) SaveRunFailedTriggerTime(value int) {
	db.DB.Update(func(txn *badger.Txn) error {
		byteVal := []byte(strconv.Itoa(value))
		e := badger.NewEntry([]byte(runFailedTriggerTime), byteVal)
		err := txn.SetEntry(e)

		return err
	})
}

func (db *AppDb) GetUsersTriggerTime() int {
	retVal := -1

	db.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(usersTriggerTime))
		if err != nil {
			return err
		}

		dst := make([]byte, 0)
		dst, _ = item.ValueCopy(dst)
		retVal, _ = strconv.Atoi(string(dst))
		return nil
	})

	return retVal
}

func (db *AppDb) SaveUsersTriggerTime(value int) {
	db.DB.Update(func(txn *badger.Txn) error {
		byteVal := []byte(strconv.Itoa(value))
		e := badger.NewEntry([]byte(usersTriggerTime), byteVal)
		err := txn.SetEntry(e)

		return err
	})
}
