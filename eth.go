package main

import (
	"fmt"
	"math/big"
)

func parseBigFloat(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}

func weiToEth(val *big.Float) *big.Float {
	return val.Quo(val, big.NewFloat(1e18))
}
