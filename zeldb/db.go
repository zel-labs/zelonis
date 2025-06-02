/*
Copyright (C) 2025 Zelonis Contributors

This file is part of Zelonis.

Zelonis is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Zelonis is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Zelonis. If not, see <https://www.gnu.org/licenses/>.
*/

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
