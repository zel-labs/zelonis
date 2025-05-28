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
