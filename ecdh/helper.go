package ecdh

import (
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"
)

var curve = elliptic.P256()

func Hash(elements [][]byte) [][]byte {
	hashes := make([]([]byte), 0)
	for _, element := range elements {
		sha := sha256.Sum256(element)
		hashes = append(hashes, sha[:])
	}
	return hashes
}

func HashBigInt(elements []*big.Int) [][]byte {
	hashes := make([]([]byte), 0)
	for _, element := range elements {
		sha := sha256.Sum256(element.Bytes())
		hashes = append(hashes, sha[:])
	}
	return hashes
}

func GetPoints(hashes [][]byte, exponent *big.Int) ([]*big.Int, []*big.Int) {
	xs := make([]*big.Int, 0)
	ys := make([]*big.Int, 0)
	for _, hash := range hashes {
		val := new(big.Int).SetBytes(hash)
		x, y := GetPoint(curve, val)

		xs = append(xs, x)
		ys = append(ys, y)
	}
	return xs, ys
}

func Exp(xs []*big.Int, ys []*big.Int, exponent *big.Int) {
	for i := 0; i < len(xs); i++ {
		x, y := curve.ScalarMult(xs[i], ys[i], exponent.Bytes())
		xs[i] = x
		ys[i] = y
	}
}
