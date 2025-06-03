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
package domain

import (
	"crypto/ed25519"
	"encoding/json"
	"golang.org/x/crypto/blake2b"
	"time"
	"zelonis/external"
	"zelonis/stats"
	"zelonis/validator/core/accounts"
	"zelonis/validator/core/block"
	"zelonis/validator/core/transaction"
	"zelonis/zeldb"
)

type Domain struct {
	accountManager *accounts.Manager
	blockManager   *block.Manager

	txManager    *transaction.Manager
	statsManager *stats.Manager
}

func NewDomain(dir string) *Domain {
	statsManager := stats.NewManager(zeldb.NewDb("stats", dir))
	accountManger := accounts.NewManager(zeldb.NewDb("accounts", dir), statsManager)
	blockManager := block.NewManager(zeldb.NewDb("blocks", dir), statsManager)

	txManager := transaction.NewManager(zeldb.NewDb("tx", dir), accountManger, statsManager)

	domain := &Domain{
		accountManager: accountManger,
		blockManager:   blockManager,

		txManager:    txManager,
		statsManager: statsManager,
	}

	return domain
}

func (d *Domain) VerifyInsertBlockAndTransaction(block *external.Block) (bool, error) {
	// Check TxHash and add if missing

	block = d.checkTXHash(block)
	err := d.txManager.VerifyTxs(block.Transactions, block.Header.BlockHeight)
	if err != nil {
		return false, err
	}

	status, err := d.blockManager.VerifyAndAddBlock(block)
	if err != nil {
		return false, err
	}
	if status {
		return true, nil
	}
	return true, nil
}

func (d *Domain) checkTXHash(block *external.Block) *external.Block {
	for key, tx := range block.Transactions {
		if tx.TxHash == nil {
			tempTx := &external.Transaction{
				Inpoint:   tx.Inpoint,
				Outpoints: tx.Outpoints,
			}
			txHash := blake2b.Sum256(tempTx.TxSerialize())
			block.Transactions[key].TxHash = txHash[:]
		}
	}
	return block
}

func (d *Domain) GetAccountBalance(account []byte) (*external.Account, bool) {
	accountInfo, status := d.accountManager.GetAccount(account)
	if !status {
		return nil, false
	}
	return accountInfo, true
}

func (d *Domain) GetHighestBlockHash() ([]byte, error) {
	return d.blockManager.GetHighestBlockHash()
}

func (d *Domain) GetBlockByHash(hash []byte) (*external.Block, error) {
	return d.blockManager.GetBlockByHash(hash)
}

func (d *Domain) BlockManager() *block.Manager {
	return d.blockManager
}
func (d *Domain) AccountManager() *accounts.Manager {
	return d.accountManager
}

func (d *Domain) TxManager() *transaction.Manager {
	return d.txManager
}
func (d *Domain) StatsManager() *stats.Manager {
	return d.statsManager
}

func (d *Domain) StartValidatorMode(priv ed25519.PrivateKey, wallet []byte) *external.Block {
	unsignedBlock := d.buildUnsignedBlock(wallet)
	serialBlock, _ := json.Marshal(unsignedBlock)
	hash := blake2b.Sum256(serialBlock)

	unsignedBlock.Header.BlockHash = hash[:]
	unsignedBlock.Signature = d.signBlock(hash[:], priv)
	return unsignedBlock
	//Sign block
}

func (d *Domain) signBlock(hash []byte, priKey ed25519.PrivateKey) []byte {

	sig := ed25519.Sign(priKey, hash)
	if ed25519.Verify(priKey.Public().(ed25519.PublicKey), hash, sig) {

		return sig

	}
	return nil
}

func (d *Domain) buildUnsignedBlock(wallet []byte) *external.Block {
	header, valdiator := d.buildProposeBlockHeader(wallet)
	return &external.Block{
		Header:       header,
		Transactions: d.txManager.Mempool().GetMempoolTxs(),
		Validator:    valdiator,
	}
}

func (d *Domain) buildProposeBlockHeader(wallet []byte) (*external.Header, *external.ValidatorInfo) {
	blockHeight, err := d.statsManager.GetHighestBlockHeight()
	if err != nil {
		//This should not happen so panic
		panic(err)
	}
	recentBlockHash, err := d.GetHighestBlockHash()
	if err != nil {
		panic(err)
	}

	return &external.Header{
			BlockHeight: blockHeight + 1,
			BlockTime:   time.Now().UnixMilli(),
			ParentSlot:  0, //Get current Epoch
			ParentHash:  recentBlockHash,
			Version:     1,
		}, &external.ValidatorInfo{
			Pubkey:            wallet,
			PreviousBlockHash: recentBlockHash,
		}
}
