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
	"encoding/hex"
)

func (d *File) fromJSON(fileJSON *keysFileJSON) error {
	d.Version = fileJSON.Version

	winfo := make([]*PublicWalletAddress, len(fileJSON.ExtendedPublicKeys))
	for i, WalletPubkey := range fileJSON.ExtendedPublicKeys {
		winfo[i] = &PublicWalletAddress{
			Address: WalletPubkey.Address,
		}

	}
	d.ExtendedAddress = winfo

	EncryptedMnemonics, err := d.decryptAll(fileJSON.EncryptedPrivateKeys)
	if err != nil {
		return err
	}
	d.EncryptedMnemonics = EncryptedMnemonics

	EncryptedTimestamp, err := d.decryptAll(fileJSON.EncryptionKey)
	if err != nil {
		return err
	}
	d.EncrytionKey = EncryptedTimestamp

	ExtendedPublicKeys := d.decryptPublicKey(fileJSON.ExtendedPublicKeys)

	d.ExtendedAddress = ExtendedPublicKeys

	return nil
}

func (d *File) decryptAll(info []*encryptedInfoJson) ([]*EncryptedInfo, error) {
	rinfo := make([]*EncryptedInfo, len(info))
	for i, encryptedInfoJson := range info {
		cipher, err := hex.DecodeString(encryptedInfoJson.Cipher)
		if err != nil {
			return nil, err
		}

		salt, err := hex.DecodeString(encryptedInfoJson.Salt)
		if err != nil {
			return nil, err
		}

		rinfo[i] = &EncryptedInfo{
			cipher: cipher,
			salt:   salt,
		}
	}
	return rinfo, nil
}

func (d *File) decryptPublicKey(info []*publicWalletAddress) []*PublicWalletAddress {
	rinfo := make([]*PublicWalletAddress, len(info))
	for i, extendedInfoJson := range info {
		rinfo[i] = &PublicWalletAddress{
			Address: extendedInfoJson.Address,
		}
	}
	return rinfo
}
