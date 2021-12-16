package ecdh

import (
	"crypto/elliptic"
	"fmt"
	"math/big"
)

// from crypto/elliptic/elliptic.go
// polynomial returns x³ - 3x + b.
func polynomial(curve elliptic.Curve, x *big.Int) *big.Int {
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)

	threeX := new(big.Int).Lsh(x, 1)
	threeX.Add(threeX, x)

	x3.Sub(x3, threeX)
	x3.Add(x3, curve.Params().B)
	x3.Mod(x3, curve.Params().P)

	return x3
}

func calcY(curve elliptic.Curve, x *big.Int, y *big.Int) bool {
	res := y.ModSqrt(polynomial(curve, x), curve.Params().P)
	return res != nil
}

func GetPoint(curve elliptic.Curve, x *big.Int) (*big.Int, *big.Int) {
	px := new(big.Int).Set(x)
	py := new(big.Int)
	one := big.NewInt(1)
	// 256 is used in https://github.com/encryptogroup/PSI
	px.Lsh(px, 8)
	px.Mod(px, curve.Params().P)
	for i := 0; i < (2 << 8); i++ {
		if calcY(curve, px, py) {
			if !curve.IsOnCurve(px, py) {
				panic("")
			}
			return px, py
		}
		px.Add(px, one)
	}
	panic(fmt.Sprintf("failed to find point for %v\n", x))
}
