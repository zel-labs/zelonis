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
	Balance             []byte `json:"balance"`
	Stake               []byte `json:"stake"`
	ActivatingStake     []byte `json:"activating_stake"`
	DeactivatingStake   []byte `json:"deactivating_stake"`
	PendingActivation   []byte `json:"pending_activation"`
	PendingDeactivation []byte `json:"pending_deactivation"`
	WarmupStake         []byte `json:"warmup_stake"`
	CoolingDownStake    []byte `json:"cooling_down_stake"`
	Reward              []byte `json:"reward"`
}

func (ac *Account) SerializedToValidatorAccount(infoBytes []byte) error {
	return json.Unmarshal(infoBytes, ac)
}
func (ac *Account) ValidatorAccountToSerialized() ([]byte, error) {
	return json.Marshal(ac)
}

func (ac *Account) AccountBalanceBigFloat() *big.Float {
	balance, _ := maths.BytesToBigFloatString(ac.Balance)
	return balance
}
func (ac *Account) AccountStakeBigFloat() *big.Float {
	balance, _ := maths.BytesToBigFloatString(ac.Stake)
	return balance
}
func (ac *Account) AccountRewardBigFloat() *big.Float {
	balance, _ := maths.BytesToBigFloatString(ac.Reward)
	return balance
}

func (ac *Account) ReduceBalance(val []byte, fee []byte) (*big.Float, bool) {
	valBig, _ := maths.BytesToBigFloatString(val)
	feeBig, _ := maths.BytesToBigFloatString(fee)

	totalValue := new(big.Float).Add(valBig, feeBig)
	accountVal := ac.AccountBalanceBigFloat()
	accountVal = accountVal.Sub(accountVal, totalValue)
	if accountVal.Cmp(big.NewFloat(0)) == -1 {
		return accountVal, false
	}
	ac.Balance = []byte(accountVal.String())
	return accountVal, true
}

func (ac *Account) TestReduceBalance(val []byte, fee []byte) bool {
	valBig, _ := maths.ByteTomZel(val)
	feeBig, _ := maths.ByteTomZel(fee)
	totalVal := new(big.Int).Add(valBig, feeBig)
	accountVal, _ := maths.ByteTomZel(ac.Balance)
	accountVal = accountVal.Sub(accountVal, totalVal)

	if accountVal.Cmp(big.NewInt(0)) == -1 {
		return false
	}
	newVal := maths.MZelToZelByte(accountVal)
	ac.Balance = []byte(newVal.String())
	return true

}
func (ac *Account) TestAddBalance(val []byte) bool {
	valBig, _ := maths.ByteTomZel(val)

	accountVal, _ := maths.ByteTomZel(ac.Balance)
	accountVal = accountVal.Add(accountVal, valBig)
	newVal := maths.MZelToZelByte(accountVal)
	ac.Balance = []byte(newVal.String())
	return true
}

func (ac *Account) AddBalance(val []byte) *big.Float {
	valBig, _ := maths.BytesToBigFloatString(val)

	accountVal := ac.AccountBalanceBigFloat()
	accountVal = accountVal.Add(accountVal, valBig)
	ac.Balance = []byte(accountVal.String())
	return accountVal
}

func (ac *Account) AddStake(val []byte) *big.Float {
	valBig, _ := maths.BytesToBigFloatString(val)
	accountVal := ac.AccountStakeBigFloat()
	accountVal = accountVal.Add(accountVal, valBig)
	ac.Stake = []byte(accountVal.String())
	return accountVal

}
