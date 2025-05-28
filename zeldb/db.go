package zeldb

import (
	"errors"
	"github.com/dgraph-io/badger/v4"
)

type ZelDB struct {
	Name string
	*badger.DB
}

func (db *ZelDB) Has(key []byte) (bool, error) {

	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err

	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (db *ZelDB) Get(key []byte) ([]byte, error) {
	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			value = append([]byte{}, val...) // deep copy
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (db *ZelDB) Set(key, value []byte) error {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
