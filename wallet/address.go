package wallet

import (
	"github.com/pkg/errors"
)

func ListAddress(path string) error {
	keysFile, err := ReadKeysFile(path)
	if err != nil {
		return errors.Wrapf(err, "Error reading keys file %s", keysFile)
	}

	return nil
}
