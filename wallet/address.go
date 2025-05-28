package wallet

import (
	"fmt"
)

func ListAddress(path string) error {
	keysFile, err := ReadKeysFile(path)
	if err != nil {
		return fmt.Errorf("%s Error reading keys file %s", err, keysFile)
	}

	return nil
}
