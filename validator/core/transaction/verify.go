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
package transaction

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/mr-tron/base58"
	"log"
	"math/big"
	"reflect"
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

func (m *Manager) VerifyTxs(txs []*external.Transaction, blockHeight uint64) error {
	for _, tx := range txs {

		if m.checkIfTxExists(tx) {

			continue
		}

		m.verifyTx(tx, blockHeight)

		//Remove tx from mempool
		m.mempool.RemoveTxFromMempool(tx)
		m.statsManager.UpdateTotalTransactions()
		m.statsManager.UpdateUtxoBlockHeight(blockHeight)
		if tx.Status == external.TxAccepted {
			if tx.TxType == external.TxStakingSend {
				m.statsManager.UpdateTotalStake(tx.Inpoint.Value)
			}
			if reflect.DeepEqual(tx.Inpoint.PubKey, coreSender) {
				m.statsManager.UpdateTotalSupply(tx.Inpoint.Value)
			}
		}

		//txs[key] = ntx
	}
	return nil
}

func (m *Manager) verifyTx(tx *external.Transaction, blockHeight uint64) error {
	if status, _ := m.db.Has(tx.TxHash); status {

		return errors.New("tx already exists")
	}
	m.addTxToAccount(tx, blockHeight)
	txBlock := []byte(fmt.Sprintf("TxBlock-%x", tx.TxHash))
	err := m.db.Set(txBlock, []byte(fmt.Sprintf("%v", blockHeight)))
	if err != nil {
		return err
	}
	err = m.db.Set(tx.TxHash, tx.TxSerialize())
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) addTxToAccount(tx *external.Transaction, blockHeight uint64) {

	m.accountManger.AddAccountTransaction(tx.Inpoint.PubKey, tx.TxHash)
	status := m.txValueSanity(tx)
	if !status {
		tx.Status = external.TXRejectedDueBalanceMutation
		return
	}

	if tx.TxType == external.TxTransfer || (tx.TxType == 0 && blockHeight == 0) {
		m.transferTx(tx, blockHeight)
		log.Println(tx.Status)
		//os.Exit(12)
	} else if tx.TxType == external.TxStakingSend {
		m.stakingTx(tx, blockHeight)
	} else if tx.TxType == external.TxStakingRelease {
		//m.stakingRelease(tx, blockHeight)
	} else if tx.TxType == external.TxRewardSend {
		//m.rewardTx(tx, blockHeight)
	} else if tx.TxType == external.TxRewardRelease {
		//m.rewardReleaseTx(tx, blockHeight)
	}

}

func (m *Manager) stakingTx(tx *external.Transaction, blockHeight uint64) {
	sender := tx.Inpoint.PubKey
	senderVal, _ := maths.BytesToBigFloatString(tx.Inpoint.Value)

	if !m.coreSenderSanity(sender) {

		if !m.signatureSanity(tx) {
			tx.Status = external.TxRejctedDueToSignatureMismatch
			return
		}

		sa, acstatus := m.accountManger.GetAccount(sender)

		if !acstatus {
			tx.Status = external.TXRejectedDueToBalance
			return
		}

		if sa.AccountBalanceBigFloat().Cmp(senderVal) >= 0 {

			//Reduce balance from sender
			//Add Balance to reciver
			status := sa.TestReduceBalance(tx.Inpoint.Value, tx.Fee)

			sa.AddStake(tx.Inpoint.Value)
			m.accountManger.UpdateAccount(sa, tx.Inpoint.PubKey)
			m.accountManger.AddTxToAccountHistory(tx.Inpoint.PubKey, tx.TxHash)
			if !status {
				tx.Status = external.TXRejectedDueToBalance
				return
			}
			tx.Status = external.TxAccepted
		}

	}
}

func (m *Manager) transferTx(tx *external.Transaction, blockHeight uint64) {
	sender := tx.Inpoint.PubKey
	senderVal, _ := maths.BytesToBigFloatString(tx.Inpoint.Value)
	if !m.coreSenderSanity(sender) {

		//check if tranaction already credited to user db

		//Verify Transaction
		if !m.signatureSanity(tx) {
			tx.Status = external.TxRejctedDueToSignatureMismatch
			return
		}
		//Get account

		sa, acstatus := m.accountManger.GetAccount(sender)
		if !acstatus {
			tx.Status = external.TXRejectedDueToBalance
			return
		}

		if sa.AccountBalanceBigFloat().Cmp(senderVal) >= 0 {
			//Reduce balance from sender
			//Add Balance to reciver
			status := sa.TestReduceBalance(tx.Inpoint.Value, tx.Fee)
			if !status {
				tx.Status = external.TXRejectedDueToBalance
				return
			}
			m.accountManger.UpdateAccount(sa, tx.Inpoint.PubKey)
			m.accountManger.AddTxToAccountHistory(tx.Inpoint.PubKey, tx.TxHash)
			acstatus = m.updateReciverAccount(tx)
			if acstatus {
				tx.Status = external.TxAccepted

				return
			}
		}

	}
	if m.coreSenderSanity(sender) && blockHeight == 0 {

		status := m.updateReciverAccount(tx)
		if status {
			tx.Status = external.TxAccepted
			return
		}

	}
}

func (m *Manager) updateReciverAccount(tx *external.Transaction) bool {
	for _, Outpoint := range tx.Outpoints {

		ra, status := m.accountManger.GetAccount(Outpoint.PubKey)

		if !status {
			ra = &external.Account{
				Balance:             []byte("0"),
				Stake:               []byte("0"),
				ActivatingStake:     []byte("0"),
				DeactivatingStake:   []byte("0"),
				PendingActivation:   []byte("0"),
				PendingDeactivation: []byte("0"),
				WarmupStake:         []byte("0"),
				CoolingDownStake:    []byte("0"),
				Reward:              []byte("0"),
			}
		}
		ra.TestAddBalance(Outpoint.Value)
		m.accountManger.UpdateAccount(ra, Outpoint.PubKey)
		m.accountManger.AddTxToAccountHistory(Outpoint.PubKey, tx.TxHash)
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
	pubKey, _ := base58.Decode(fmt.Sprintf("%s", tx.Inpoint.PubKey))

	pubkey := ed25519.PublicKey(pubKey)
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
	cmp = inpointTotalB.Cmp(big.NewFloat(0))
	if cmp <= 0 {
		return false
	}
	return true
}

func (m *Manager) GetTransactionByHash(hash string) (*external.Transaction, error) {

	hashBytes, err := external.NewDomainHashFromString(hash)
	if err != nil {
		return nil, err
	}
	status, err := m.db.Has(hashBytes.ByteSlice())
	if err != nil {
		return nil, err
	}
	if !status {
		return nil, fmt.Errorf("tx not exist")
	}
	tx := m.getTransaction(hashBytes.ByteSlice())
	return tx, nil

}

func (m *Manager) CheckIfTxExists(tx *external.Transaction) bool {
	return m.checkIfTxExists(tx)
}

func (m *Manager) GetTransactionBlockHeight(hash string) (string, error) {

	txBlock := []byte(fmt.Sprintf("TxBlock-%s", hash))
	heightByte, err := m.db.Get(txBlock)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", heightByte), nil
}
