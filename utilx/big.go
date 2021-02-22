package utilx

import (
	"fmt"
	"math/big"
)

func stringToBigInt(s1, s2 string) (n1, n2 *big.Int, err error) {

	var (
		ok bool
	)

	if s1 == "" {
		n1 = big.NewInt(0)
	} else {
		n1, ok = new(big.Int).SetString(s1, 10)
		if !ok {
			err = fmt.Errorf("param 1:%s cast to big.int err", s1)
			return
		}
	}

	if s2 == "" {
		n2 = big.NewInt(0)
	} else {

		n2, ok = new(big.Int).SetString(s2, 10)
		if !ok {
			err = fmt.Errorf("param 2:%s cast to big.int err", s2)
			return
		}

	}

	return
}

func stringToBigFloat(s1, s2 string) (n1, n2 *big.Float, err error) {

	var (
		ok bool
	)

	if s1 == "" {
		n1 = big.NewFloat(0)
	} else {
		n1, ok = new(big.Float).SetString(s1)
		if !ok {
			err = fmt.Errorf("param 1:%s cast to big.float err", s1)
			return
		}
	}

	if s2 == "" {
		n2 = big.NewFloat(0)
	} else {

		n2, ok = new(big.Float).SetString(s2)
		if !ok {
			err = fmt.Errorf("param 2:%s cast to big.float err", s2)
			return
		}

	}

	return
}

func BigIntAdd(s1, s2 string) (rs string, err error) {

	n1, n2, err := stringToBigInt(s1, s2)
	if err != nil {
		return
	}

	rs = n1.Add(n1, n2).String()

	return
}

func BigIntDiv(s1, s2 string) (rs string, err error) {
	if s1 == "" || s1 == "0" || s2 == "" || s2 == "0" {
		return
	}
	n1, n2, err := stringToBigInt(s1, s2)
	if err != nil {
		return
	}

	rs = n1.Div(n1, n2).String()

	return
}

func BigIntMul(s1, s2 string) (rs string, err error) {

	n1, n2, err := stringToBigInt(s1, s2)
	if err != nil {
		return
	}

	rs = n1.Mul(n1, n2).String()

	return
}

func BigFloatAdd(s1, s2 string) (rs string, err error) {

	n1, n2, err := stringToBigFloat(s1, s2)
	if err != nil {
		return
	}

	rs = n1.Add(n1, n2).String()

	return
}

func BigFloatMul(s1, s2 string) (rs string, err error) {

	n1, n2, err := stringToBigFloat(s1, s2)
	if err != nil {
		return
	}

	rs = n1.Mul(n1, n2).String()

	return
}

func BigFloatDiv(s1, s2 string) (rs string, err error) {

	n1, n2, err := stringToBigFloat(s1, s2)
	if err != nil {
		return
	}

	rs = n1.Quo(n1, n2).String()

	return
}
