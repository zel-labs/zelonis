package domain

import (
	"golang.org/x/crypto/blake2b"
	"zelonis/external"
	"zelonis/stats"
	"zelonis/validator/accounts"
	"zelonis/validator/core"
	"zelonis/validator/core/block"
	"zelonis/validator/core/transaction"
	"zelonis/zeldb"
)

type Domain struct {
	accountManager *accounts.Manager
	blockManager   *block.Manager

	txManager    *transaction.Manager
	statsManager *stats.Manager
	coreManager  *core.Core
}

func NewDomain(dir string) *Domain {
	statsManager := stats.NewManager(zeldb.NewDb("stats", dir))
	accountManger := accounts.NewManager(zeldb.NewDb("accounts", dir), statsManager)
	blockManager := block.NewManager(zeldb.NewDb("blocks", dir), statsManager)

	txManager := transaction.NewManager(zeldb.NewDb("tx", dir), accountManger, statsManager)
	coreManager := core.New(accountManger, txManager, blockManager, statsManager)
	domain := &Domain{
		accountManager: accountManger,
		blockManager:   blockManager,

		txManager:    txManager,
		statsManager: statsManager,
		coreManager:  coreManager,
	}

	return domain
}

func (d *Domain) VerifyInsertBlockAndTransaction(block *external.Block) (bool, error) {
	// Check TxHash and add if missing
	d.GetHighestBlockHash()
	block = d.checkTXHash(block)

	status, err := d.blockManager.VerifyAndAddBlock(block)
	if err != nil {
		return false, err
	}
	if status {
		return true, nil
	}
	if err = d.txManager.VerifyTxs(block.Transactions, block.Header.BlockHeight); err != nil {
		return false, err
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

func (d *Domain) CreateBlockWithTransaction(tx *external.Transaction) (bool, error) {
	txPool := []*external.Transaction{tx}
	if err := d.txManager.VerifyTxs(txPool, 0); err != nil {
		return false, err
	}
	return true, nil
}

func (d *Domain) GetHighestBlockHash() ([]byte, error) {
	return d.blockManager.GetHighestBlockHash()
}

func (d *Domain) GetBlockByHash(hash []byte) (*external.Block, error) {
	return d.blockManager.GetBlockByHash(hash)
}
