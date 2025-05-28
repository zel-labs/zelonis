package accounts

import (
	"fmt"
	"log"
	"time"
	"zelonis/external"
)

func (m *Manager) AddAccountTransaction(account []byte, txHash []byte) (bool, error) {
	//Build Transaction Key
	key := time.Now().UnixNano()
	txKey := []byte(fmt.Sprintf("%s:%v", account, key))

	err := m.db.Set(txKey, txHash)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *Manager) GetAccount(account []byte) (*external.Account, bool) {
	accountStatus, err := m.db.Has(account)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	if !accountStatus {
		
		return nil, false
	}
	//Account exists check balance
	accountBytes, err := m.db.Get(account)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	userAccount := &external.Account{}
	err = userAccount.SerializedToValidatorAccount(accountBytes)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return userAccount, true
}

func (m *Manager) UpdateAccount(account *external.Account, pubkey []byte) bool {
	acBytes, err := account.ValidatorAccountToSerialized()
	if err != nil {
		return false
	}
	err = m.db.Set(pubkey, acBytes)
	if err != nil {
		log.Println("Error updating account")
		return false
	}
	return true
}
