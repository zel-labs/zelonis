package external

import (
	"encoding/json"
	"math/big"
	"zelonis/utils/maths"
)

type Account struct {
	Balance []byte `json:"balance"`
	Stake   []byte `json:"stake"`
	Reward  []byte `json:"reward"`
}

func (account *Account) SerializedToValidatorAccount(infoBytes []byte) error {
	return json.Unmarshal(infoBytes, account)
}
func (account *Account) ValidatorAccountToSerialized() ([]byte, error) {
	return json.Marshal(account)
}

func (account *Account) AccountBalanceBigFloat() *big.Float {
	balance, _ := maths.BytesToBigFloatString(account.Balance)
	return balance
}

func (account *Account) ReduceBalance(val []byte, fee []byte) *big.Float {
	valBig, _ := maths.BytesToBigFloatString(val)
	feeBig, _ := maths.BytesToBigFloatString(val)
	totalValue := new(big.Float).Add(valBig, feeBig)
	accountVal := account.AccountBalanceBigFloat()
	accountVal = accountVal.Sub(accountVal, totalValue)
	account.Balance = []byte(accountVal.String())
	return accountVal
}

func (account *Account) AddBalance(val []byte) *big.Float {
	valBig, _ := maths.BytesToBigFloatString(val)

	accountVal := account.AccountBalanceBigFloat()
	accountVal = accountVal.Add(accountVal, valBig)
	account.Balance = []byte(accountVal.String())
	return accountVal
}
