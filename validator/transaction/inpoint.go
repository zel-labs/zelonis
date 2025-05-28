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
