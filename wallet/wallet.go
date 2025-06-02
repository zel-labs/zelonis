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
	"crypto/ed25519"
	"fmt"
	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"log"
	"strconv"
	"time"
)

type walletInfo struct {
	Mnemonic   string             `json:"mnemonic"`
	PrivateKey ed25519.PrivateKey `json:"privateKey"`
	Address    string             `json:"address"`
	EncryptKey int64              `json:"encryptKey"`
}

func CreateWallet() walletInfo {
	encryptKey := time.Now().UnixNano()
	return createAndRecoverWallet(encryptKey, false, "")
}

func RecoverWallet(encryptKeyStr string, oldMnemonic string) walletInfo {
	encryptKey, _ := strconv.ParseInt(encryptKeyStr, 10, 64)

	return createAndRecoverWallet(encryptKey, true, oldMnemonic)
}

func createAndRecoverWallet(encryptKey int64, isRecovery bool, oldmnemonic string) walletInfo {
	var mnemonic string
	if !isRecovery {
		entropy, _ := bip39.NewEntropy(128)
		mnemonic, _ = bip39.NewMnemonic(entropy)
	} else {
		mnemonic = oldmnemonic
	}

	passPhrase := fmt.Sprintf("%d", encryptKey)
	// 2. Derive seed from mnemonic (no passphrase)
	seed := bip39.NewSeed(mnemonic, passPhrase)

	// 3. Use BIP-32 to get master keys
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatal(err)
	}

	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coinType, _ := purpose.NewChildKey(bip32.FirstHardenedChild + 501)
	account, _ := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	change, _ := account.NewChildKey(0)
	indexKey, _ := change.NewChildKey(0)

	var seed32 [32]byte
	copy(seed32[:], indexKey.Key[:32])
	priv := ed25519.NewKeyFromSeed(seed32[:])
	pub := priv.Public().(ed25519.PublicKey)

	address := base58.Encode(pub)

	return walletInfo{
		PrivateKey: priv,
		Mnemonic:   mnemonic,
		Address:    address,
		EncryptKey: encryptKey,
	}
}
