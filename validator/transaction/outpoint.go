package transaction

import (
	"math/big"
	"zelonis/external"
	"zelonis/utils/maths"
)

func (m *Manager) outpointsTotal(outpoints []*external.Outpoint) (*big.Float, error) {
	totalVal := big.NewFloat(0)
	for _, outpoint := range outpoints {

		outVal, err := maths.BytesToBigFloatString(outpoint.Value)
		if err != nil {
			return nil, err
			break
		}
		totalVal = new(big.Float).Add(totalVal, outVal)
	}
	return totalVal, nil
}
