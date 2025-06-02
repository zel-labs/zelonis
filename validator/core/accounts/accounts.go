package accounts

import (
	"fmt"
	"log"
	"strconv"
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

func (m *Manager) GetAccountTxHeight(pubkey []byte) (int, error) {
	key := []byte(fmt.Sprintf("%x-h", pubkey))
	status, err := m.db.Has(key)
	if err != nil {
		return 0, err
	}
	if !status {
		return 0, nil
	}
	heightByte, err := m.db.Get(key)
	if err != nil {
		return 0, err
	}
	height, err := strconv.Atoi(string(heightByte))
	if err != nil {
		return 0, err
	}
	return height, nil

}

func (m *Manager) UpdateAccountTxHeight(pubkey []byte) (int, error) {
	key := []byte(fmt.Sprintf("%x-h", pubkey))
	height, err := m.GetAccountTxHeight(pubkey)
	if err != nil {
		return 0, err
	}
	height = height + 1

	err = m.db.Set(key, []byte(strconv.Itoa(height)))
	if err != nil {
		return 0, err
	}
	return height, nil
}

func (m *Manager) AddTxToAccountHistory(pubkey []byte, txHash []byte) error {
	ch, err := m.UpdateAccountTxHeight(pubkey)
	if err != nil {
		return err
	}
	key := []byte(fmt.Sprintf("%x-%v", pubkey, ch))
	err = m.db.Set(key, txHash)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) getAccountTxByKey(pubkey []byte, height int) ([]byte, error) {
	key := []byte(fmt.Sprintf("%x-%v", pubkey, height))
	txByte, err := m.db.Get(key)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func (m *Manager) GetTxHistoryList(pubkey []byte, from, limit int) ([]string, error) {
	height, err := m.GetAccountTxHeight(pubkey)
	if err != nil {
		return nil, err
	}
	if height == 0 {
		return nil, nil
	}
	txs := make([]string, 0)
	from = height - from
	if limit > 100 {
		limit = 100
	} else if limit < 10 {
		limit = 10
	}
	counter := 0
	for height > 0 {
		counter++
		tx, err := m.getAccountTxByKey(pubkey, height)
		if err != nil {
			return nil, err
		}

		txs = append(txs, fmt.Sprintf("%x", tx))
		if counter == limit {
			break
		}
		height--

	}

	return txs, nil
}
