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
	"math/big"
	"zelonis/external"
	"zelonis/utils/maths"
)

func (m *Manager) calculateExpense(tx *external.Transaction) *big.Float {
	sentAmount, err := maths.BytesToBigFloatString(tx.Inpoint.Value)
	if err != nil {
		return nil
	}
	feeAmount, err := maths.BytesToBigFloatString(tx.Fee)
	if err != nil {
		return nil
	}
	totalExpense := new(big.Float).Add(feeAmount, sentAmount)
	return totalExpense

}
