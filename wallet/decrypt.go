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
