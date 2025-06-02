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

package wallet

import (
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

var flockMap = make(map[string]*flock.Flock)

func (d *File) TryLock() error {
	if _, ok := flockMap[d.path]; ok {
		return errors.Errorf("file %s is already locked", d.path)
	}

	lockFile := flock.New(d.path + ".lock")
	err := createFileDirectoryIfDoesntExist(lockFile.Path())
	if err != nil {
		return err
	}

	flockMap[d.path] = lockFile

	success, err := lockFile.TryLock()
	if err != nil {
		return err
	}

	if !success {
		return errors.Errorf("%s is locked and cannot be used. Make sure that no other active wallet command is using it.", d.path)
	}
	return nil
}

func createFileDirectoryIfDoesntExist(path string) error {
	dir := filepath.Dir(path)
	exists, err := pathExists(dir)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return os.MkdirAll(dir, 0700)
}
