package ecdh

import (
	"math/big"
	"math/rand"
)

func Step1(elements [][]byte) ([][]byte, *big.Int, []*big.Int, []*big.Int) {
	// Hash each elements with sha256.
	hashes := Hash(elements)
	// Pick a random exponent.
	exponent := big.NewInt(rand.Int63n(10) + 5)
	// Get a points from elliptic curve P256 from the hashes.
	xs, ys := GetPoints(hashes, exponent)
	// Exp the points with the random exponent.
	Exp(xs, ys, exponent)
	return hashes, exponent, xs, ys
}

func Step2(peerXs []*big.Int, peerYs []*big.Int, exponent *big.Int) [][]byte {
	// Exp the peer points with the random exponent.
	Exp(peerXs, peerYs, exponent)
	// Hash the peer x.
	peerHashes := HashBigInt(peerXs)
	return peerHashes
}

// Return the intersection index
func Intersect(hashes [][]byte, peerHashes [][]byte) []int {
	hashSet := make(map[string]bool)
	for _, peerHash := range peerHashes {
		hashSet[string(peerHash)] = true
	}

	intersection := make([]int, 0)
	for i, hash := range hashes {
		if _, ok := hashSet[string(hash)]; ok {
			intersection = append(intersection, i)
		}
	}
	return intersection
}
