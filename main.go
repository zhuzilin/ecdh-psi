package main

import (
	"crypto/elliptic"
	"crypto/sha256"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

const (
	alice = "alice"
	bob   = "bob"
)

const NUM_ELEMENT = 10
const UPPER_BOUND = 20

func main() {
	rand.Seed(time.Now().UnixNano())

	role := flag.String("role", "", "Role of the party.")
	aliceAddr := flag.String("alice", "localhost:12346", "address of alice")
	bobAddr := flag.String("bob", "localhost:23455", "address of bob")

	flag.Parse()

	if *role == "" {
		panic("Must provide -role flag. Possible values are `bob` and `alice`")
	} else if *role != alice && *role != bob {
		panic(fmt.Sprintf("-role must be `bob` or `alice`, got %s", *role))
	}

	var myAddr, peerAddr, serverRole string
	if *role == alice {
		myAddr = *aliceAddr
		peerAddr = *bobAddr
		serverRole = SERVER
	} else {
		myAddr = *bobAddr
		peerAddr = *aliceAddr
		serverRole = CLIENT
	}
	conn := NewConn(serverRole, myAddr, peerAddr)
	defer conn.Close()

	// create values
	elements := make([]*big.Int, 0)
	for i := 0; i < NUM_ELEMENT; i++ {
		element := big.NewInt(rand.Int63n(int64(UPPER_BOUND)))
		elements = append(elements, element)
	}
	fmt.Printf("role: %v\n", *role)
	fmt.Printf("elements of %v: %v\n", *role, elements)

	// hash all values
	hashes := make([]([]byte), 0)
	for _, element := range elements {
		sha := sha256.Sum256(element.Bytes())
		hashes = append(hashes, sha[:])
	}

	// encript hashes
	exponent := big.NewInt(rand.Int63n(10) + 5)
	// fmt.Printf("exponent of %v: %v\n", *role, exponent)

	xs := make([]*big.Int, 0)
	ys := make([]*big.Int, 0)
	curve := elliptic.P256()
	for _, hash := range hashes {
		val := new(big.Int).SetBytes(hash)
		x, y := GetPoint(curve, val)

		x, y = curve.ScalarMult(x, y, exponent.Bytes())

		// fmt.Printf("x[%v]: %v\n", i, x.String())
		// fmt.Printf("y[%v]: %v\n", i, y.String())

		xs = append(xs, x)
		ys = append(ys, y)
	}

	peerXs := make([]*big.Int, 0)
	peerYs := make([]*big.Int, 0)
	numBytes := 256 / 8
	buf := make([]byte, numBytes)
	for i := 0; i < len(xs); i++ {
		conn.SendReceiveSameLength(xs[i].Bytes(), buf)
		x := new(big.Int).SetBytes(buf)
		// fmt.Printf("peer x[%v]: %v\n", i, x.String())
		peerXs = append(peerXs, x)

		conn.SendReceiveSameLength(ys[i].Bytes(), buf)
		y := new(big.Int).SetBytes(buf)
		// fmt.Printf("peer y[%v]: %v\n", i, y.String())
		peerYs = append(peerYs, y)
	}

	for i := 0; i < len(xs); i++ {
		x, y := curve.ScalarMult(peerXs[i], peerYs[i], exponent.Bytes())
		peerXs[i] = x
		peerYs[i] = y
	}

	peerHashes := make([]([]byte), 0)
	for _, x := range peerXs {
		sha := sha256.Sum256(x.Bytes())
		peerHashes = append(peerHashes, sha[:])
	}

	for i, peerHash := range peerHashes {
		conn.SendReceiveSameLength(peerHash, buf)
		copy(hashes[i], buf)
	}

	hashSet := make(map[string]bool)
	for _, peerHash := range peerHashes {
		hashSet[string(peerHash)] = true
	}

	intersection := make([]string, 0)
	for i, hash := range hashes {
		if _, ok := hashSet[string(hash)]; ok {
			intersection = append(intersection, elements[i].String())
		}
	}
	fmt.Printf("intersection: %v\n", intersection)
}
