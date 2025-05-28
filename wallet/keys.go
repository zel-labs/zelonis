package wallet

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

const LastVersion = 1

type File struct {
	Version            uint32
	EncryptedMnemonics []*EncryptedInfo
	EncrytionKey       []*EncryptedInfo
	ExtendedAddress    []*PublicWalletAddress
	CosignerIndex      uint32
	path               string
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil

	}

	return false, err
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return "", errors.WithStack(err)
	}

	return strings.TrimSpace(string(line)), nil
}

func defaultKeysFile(dir string) string {
	return filepath.Join(dir, "keys.json")
}

func (d *File) SetPath(path string, forceOverride bool) error {

	path = defaultKeysFile(path)

	if !forceOverride {
		exists, err := pathExists(path)
		if err != nil {
			return err
		}

		if exists {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("The file %s already exists. Are you sure you want to override it (type 'y' to approve)? ", d.path)
			line, err := readLine(reader)
			if err != nil {
				return err
			}

			if string(line) != "y" {
				return errors.Errorf("aborted setting the file path to %s", path)
			}
		}
	}
	d.path = path
	return nil
}

func DecryptedInfo(numThreads uint8, encryptedInfo *EncryptedInfo, password []byte) (string, error) {
	aead, err := getAEAD(numThreads, password, encryptedInfo.salt)
	if err != nil {
		return "", err
	}

	if len(encryptedInfo.cipher) < aead.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedInfo.cipher[:aead.NonceSize()], encryptedInfo.cipher[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	decrypted, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
