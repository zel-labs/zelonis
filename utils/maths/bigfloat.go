package maths

import (
	"fmt"
	"math/big"
)

func BytesToBigFloatString(b []byte) (*big.Float, error) {
	f, _, err := big.ParseFloat(string(b), 10, 21, big.ToNearestEven)

	return f, err
}

func ByteTomZel(v []byte) (*big.Int, error) {

	vfloat := fmt.Sprintf("%s", v)

	vbigf := new(big.Float)
	vbigf.SetString(vfloat)

	//zbigF := new(big.Float).Abs(vbigf)

	mbigf := new(big.Float).SetFloat64(1_000_000_000)
	vbigf.Mul(vbigf, mbigf)
	vbigint := new(big.Int)
	intval, _ := vbigf.Int64()
	vbigint.SetInt64(intval)
	return vbigint, nil
}

func MZelToZelByte(v *big.Int) *big.Float {
	var prec uint = 1024
	vbigF := new(big.Float)
	vbigF.SetPrec(prec).SetString(v.String())

	mbigf := new(big.Float).SetPrec(prec).SetFloat64(1_000_000_000)

	q := vbigF.Quo(vbigF, mbigf)

	return q
}
