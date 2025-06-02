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
