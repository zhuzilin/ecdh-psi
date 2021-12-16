package main

import (
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	ecdh "github.com/zhuzilin/ecdh-psi/ecdh"
)

const (
	alice = "alice"
	bob   = "bob"
)

const NUM_ELEMENT = 10
const UPPER_BOUND = 20

func createValues(n int) []*big.Int {
	elements := make([]*big.Int, 0)
	for i := 0; i < n; i++ {
		element := big.NewInt(rand.Int63n(int64(UPPER_BOUND)))
		elements = append(elements, element)
	}
	return elements
}

func elementsToBytes(elements []*big.Int) [][]byte {
	elementBytes := make([][]byte, 0)
	for _, element := range elements {
		elementBytes = append(elementBytes, element.Bytes())
	}
	return elementBytes
}

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
	// Create Connection.
	conn := NewConn(serverRole, myAddr, peerAddr)
	defer conn.Close()

	// Create values to do PSI.
	elements := createValues(NUM_ELEMENT)
	fmt.Printf("role: %v\n", *role)
	fmt.Printf("elements of %v: %v\n", *role, elements)
	elementBytes := elementsToBytes(elements)

	// Step 1 for ECDH PSI.
	hashes, exponent, xs, ys := ecdh.Step1(elementBytes)

	// Send xs, ys to peer and receive the peerXs, peerYs.
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

	// Step 2 for ECDH PSI.
	peerHashes := ecdh.Step2(peerXs, peerYs, exponent)

	// Send the peerHashes back to peer.
	for i, peerHash := range peerHashes {
		conn.SendReceiveSameLength(peerHash, buf)
		copy(hashes[i], buf)
	}

	// Get the intersection index.
	intersection := ecdh.Intersect(hashes, peerHashes)

	intersectElements := make([]string, 0)
	for _, i := range intersection {
		intersectElements = append(intersectElements, elements[i].String())
	}

	fmt.Printf("intersection: %v\n", intersectElements)
}
