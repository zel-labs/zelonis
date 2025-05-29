package transaction

import (
	"crypto/ed25519"
	"golang.org/x/crypto/blake2b"
	"zelonis/external"
	"zelonis/stats"
	"zelonis/validator/accounts"
	"zelonis/validator/core/transaction/mempool"
	"zelonis/zeldb"
)

type Manager struct {
	db            *zeldb.ZelDB
	accountManger *accounts.Manager
	statsManager  *stats.Manager
	mempool       *mempool.TransactionsPool
}

func NewManager(db *zeldb.ZelDB, account *accounts.Manager, stats *stats.Manager) *Manager {
	tempMempool := mempool.NewMempool()

	return &Manager{
		db:            db,
		accountManger: account,
		statsManager:  stats,
		mempool:       tempMempool.NewTransactionsPool(),
	}
}

func (m *Manager) Mempool() *mempool.TransactionsPool {
	return m.mempool
}

func (m *Manager) BuildTxFromType(sender, reciver []byte, val []byte, prevBlock []byte, txType int8) *external.Transaction {
	// build inpoint
	inpoint := m.buildInpoint(sender, val, prevBlock)

	outpoint := []*external.Outpoint{
		m.buildOutpoint(reciver, val),
	}
	tx := &external.Transaction{
		Inpoint:   inpoint,
		Outpoints: outpoint,
		TxType:    txType,
	}

	//Filter fee from txType
	switch txType {
	case external.TxTransfer:
		tx.Fee = m.transaferTxFee()

	case external.TxStakingSend:
		tx.Fee = m.stakingTxFee()
	case external.TxStakingRelease:
		tx.Fee = m.stakingTxFee()
	case external.TxRewardSend:
		tx.Fee = m.rewardTxFee()
	case external.TxRewardRelease:
		tx.Fee = m.rewardTxFee()
	}

	//Build txHash
	txHash := m.buildTxHash(tx)
	tx.TxHash = txHash[:]
	return tx
}

func (m *Manager) SignTxAndVerify(tx *external.Transaction, priKey ed25519.PrivateKey) []byte {
	sig := ed25519.Sign(priKey, tx.TxHash)
	if ed25519.Verify(priKey.Public().(ed25519.PublicKey), tx.TxHash, sig) {

		return sig

	}
	return nil
}

func (m *Manager) buildTxHash(tx *external.Transaction) [32]byte {
	ntx := &external.Transaction{
		Inpoint:   tx.Inpoint,
		Outpoints: tx.Outpoints,
		TxType:    20,
		Fee:       m.stakingTxFee(),
	}
	return blake2b.Sum256(ntx.TxSerialize())

}
func (m *Manager) rewardTxFee() []byte {
	return []byte("0.00000")
}
func (m *Manager) transaferTxFee() []byte {
	return []byte("0.0000001")
}
func (m *Manager) stakingTxFee() []byte {
	return []byte("0.000005")
}

func (m *Manager) buildInpoint(addr []byte, val []byte, prevBlock []byte) *external.Inpoint {
	return &external.Inpoint{
		PubKey:        addr,
		Value:         val,
		PrevBlockHash: prevBlock,
	}
}

func (m *Manager) buildOutpoint(addr []byte, val []byte) *external.Outpoint {
	return &external.Outpoint{
		PubKey: addr,
		Value:  val,
	}
}
