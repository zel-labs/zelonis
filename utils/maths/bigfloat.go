package maths

import "math/big"

func BytesToBigFloatString(b []byte) (*big.Float, error) {
	f, _, err := big.ParseFloat(string(b), 10, 256, big.ToNearestEven)
	return f, err
}
