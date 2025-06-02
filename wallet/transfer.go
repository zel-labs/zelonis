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
	"strconv"

	"golang.org/x/crypto/blake2b"

	"zelonis/external"
)

func CreateSignedTransaction(seed string, key string, to string, val string) *external.Transaction {
	keyInt, _ := strconv.ParseInt(key, 10, 64)
	wallet := createAndRecoverWallet(keyInt, true, seed)
	return wallet.buildTx(val, to)
}

func (w *walletInfo) pubKey() []byte {
	return []byte(w.Address)
}

func (w *walletInfo) buildTx(val string, to string) *external.Transaction {

	inPoint := &external.Inpoint{
		PubKey:        w.pubKey(),
		Value:         []byte(val),
		PrevBlockHash: make([]byte, 0),
	}
	outPoint := &external.Outpoint{
		PubKey: []byte(to),
		Value:  []byte(val),
		TxType: 1,
	}

	tx := &external.Transaction{
		Inpoint:   inPoint,
		Outpoints: []*external.Outpoint{outPoint},
	}

	hash := blake2b.Sum256(tx.TxSerialize())
	tx.TxHash = hash[:]
	sig := ed25519.Sign(w.PrivateKey, hash[:])
	tx.Signature = sig
	return tx
}
