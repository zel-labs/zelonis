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
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const defaultNumThreads = 8

type EncryptedInfo struct {
	cipher []byte

	salt []byte
}

type encryptedInfoJson struct {
	Cipher string `json:"cipher"`
	Salt   string `json:"salt"`
}

type publicWalletAddress struct {
	Address []byte `json:"address"`
}

type keysFileJSON struct {
	Version              uint32                 `json:"version"`
	EncryptedPrivateKeys []*encryptedInfoJson   `json:"encryptedMnemonics"`
	ExtendedPublicKeys   []*publicWalletAddress `json:"publicKeys"`
	EncryptionKey        []*encryptedInfoJson   `json:"walletTimestamp"`
}

func encryptInfo(info string, password []byte) (*EncryptedInfo, error) {
	infoBytes := []byte(info)

	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	aead, err := getAEAD(defaultNumThreads, password, salt)
	if err != nil {
		return nil, err
	}

	// Select a random nonce, and leave capacity for the ciphertext.
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(infoBytes)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the message and append the ciphertext to the nonce.
	cipher := aead.Seal(nonce, nonce, []byte(infoBytes), nil)

	return &EncryptedInfo{
		cipher: cipher,
		salt:   salt,
	}, nil
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func getAEAD(threads uint8, password, salt []byte) (cipher.AEAD, error) {
	key := argon2.IDKey(password, salt, 1, 64*1024, threads, 32)
	return chacha20poly1305.NewX(key)
}

func (d *File) toJSON() *keysFileJSON {
	encryptedSeed := d.encryptAll(d.EncryptedMnemonics)
	encryptedAddress := d.convertAddress(d.ExtendedAddress)
	encryptedTimestamp := d.encryptAll(d.EncrytionKey)

	return &keysFileJSON{
		Version:              d.Version,
		EncryptedPrivateKeys: encryptedSeed,
		ExtendedPublicKeys:   encryptedAddress,
		EncryptionKey:        encryptedTimestamp,
	}
}

func (d *File) encryptAll(info []*EncryptedInfo) []*encryptedInfoJson {
	encryptedInfo := make([]*encryptedInfoJson, len(info))
	for i, encryptedPrivateKey := range info {
		encryptedInfo[i] = &encryptedInfoJson{
			Cipher: hex.EncodeToString(encryptedPrivateKey.cipher),
			Salt:   hex.EncodeToString(encryptedPrivateKey.salt),
		}
	}
	return encryptedInfo
}

func (d *File) convertAddress(info []*PublicWalletAddress) []*publicWalletAddress {
	encryptedInfo := make([]*publicWalletAddress, len(info))
	for i, publickey := range info {
		encryptedInfo[i] = &publicWalletAddress{
			Address: publickey.Address,
		}
	}
	return encryptedInfo
}
