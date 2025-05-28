package transaction

import (
	"crypto/ed25519"
	"errors"
	"log"
	"reflect"
	"strings"
	"zelonis/external"
	"zelonis/utils/maths"
)

var coreSender = []byte{
	0x5a, 0x65, 0x6c, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31,
	0x31, 0x31, 0x31, 0x31,
}
var genesisParentHash = strings.Repeat("0", 32)

func (m *Manager) VerifyTxs(txs []*external.Transaction, blockHeight uint64) error {
	for _, tx := range txs {

		if m.checkIfTxExists(tx) {

			continue
		}
		//Add verify signature here...
		m.verifyTx(tx, blockHeight)

	}
	return nil
}

func (m *Manager) verifyTx(tx *external.Transaction, blockHeight uint64) error {
	if status, _ := m.db.Has(tx.TxHash); status {

		return errors.New("tx already exists")
	}
	m.addTxToAccount(tx, blockHeight)
	err := m.db.Set(tx.TxHash, tx.TxSerialize())
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) addTxToAccount(tx *external.Transaction, blockHeight uint64) {
	log.Println("addTxToAccount tx:", tx.TxHash, "blockHeight:", blockHeight)
	sender := tx.Inpoint.PubKey
	senderVal, _ := maths.BytesToBigFloatString(tx.Inpoint.Value)

	m.accountManger.AddAccountTransaction(tx.Inpoint.PubKey, tx.TxHash)
	status := m.txValueSanity(tx)
	if !status {
		tx.Status = external.TXRejectedDueBalanceMutation
		return
	}

	if !m.coreSenderSanity(sender) {
		//check if tranaction already credited to user db

		//Verify Transaction
		if !m.signatureSanity(tx) {
			tx.Status = external.TxRejctedDueToSignatureMismatch
			return
		}
		//Get account

		sa, status := m.accountManger.GetAccount(sender)
		if !status {
			tx.Status = external.TXRejectedDueToBalance
			return
		}

		if sa.AccountBalanceBigFloat().Cmp(senderVal) <= 0 {
			//Reduce balance from sender
			//Add Balance to reciver
			sa.ReduceBalance(tx.Inpoint.Value, tx.Fee)
			m.accountManger.UpdateAccount(sa, tx.Inpoint.PubKey)
			status = m.updateReciverAccount(tx)
			if status {
				return
			}
		}

	}
	if m.coreSenderSanity(sender) && blockHeight == 0 {

		status = m.updateReciverAccount(tx)
		if status {
			return
		}

	}

}

func (m *Manager) updateReciverAccount(tx *external.Transaction) bool {
	for _, Outpoint := range tx.Outpoints {

		ra, status := m.accountManger.GetAccount(Outpoint.PubKey)

		if !status {
			ra = &external.Account{
				Balance: []byte("0"),
			}
		}
		ra.AddBalance(tx.Inpoint.Value)
		m.accountManger.UpdateAccount(ra, Outpoint.PubKey)

	}

	return true
}

func (m *Manager) checkIfTxExists(tx *external.Transaction) bool {
	info, err := m.db.Has(tx.TxHash)
	if err != nil {
		log.Println(err)
	}
	return info
}

func (m *Manager) getTransaction(txByte []byte) *external.Transaction {
	info, err := m.db.Get(txByte)
	if err != nil {
		log.Println(err)
	}
	tx := new(external.Transaction)
	tx.DbTxToDomainTX(info)
	return tx
}

func (m *Manager) signatureSanity(tx *external.Transaction) bool {

	pubkey := ed25519.PublicKey(tx.Inpoint.PubKey)
	if ed25519.Verify(pubkey, tx.TxHash, tx.Signature) {

		log.Printf("Verify Hash %x", tx.TxHash)
		return true
	}
	return false
}

func (m *Manager) coreSenderSanity(sender []byte) bool {
	if reflect.DeepEqual(sender, coreSender) {
		return true
	}
	return false
}

func (m *Manager) txValueSanity(tx *external.Transaction) bool {
	inpointTotalB, err := maths.BytesToBigFloatString(tx.Inpoint.Value)
	if err != nil {
		log.Println(err)
		return false
	}
	outpointTotal, err := m.outpointsTotal(tx.Outpoints)
	if err != nil {
		log.Println("outpointsTotal error:", err)
		return false
	}
	cmp := inpointTotalB.Cmp(outpointTotal)
	if cmp != 0 {
		log.Println("Inpoint total does not match outpoint total")
		return false
	}
	return true
}
