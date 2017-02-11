package main

import (
	"fmt"
	"math"
)

func main() {
	// NOTE: see http://www.vldb.org/pvldb/vol8/p1816-teller.pdf
	vals := []float64{12.0, 24.0}
	bits := float64SliceToBits(vals)
	fmt.Printf("%x %x\n", bits[0], bits[1])
	xor := bits[0] ^ bits[1]
	fmt.Printf("%016x\n", xor)
	fmt.Printf("leadingZeros=%d\n", numOfLeadingZeros(xor))
}

func float64SliceToBits(vals []float64) []uint64 {
	intVals := make([]uint64, len(vals))
	for i := 0; i < len(vals); i++ {
		intVals[i] = math.Float64bits(vals[i])
	}
	return intVals
}

func numOfLeadingZeros(v uint64) uint8 {
	mask := uint64(0x8000000000000000)
	c := uint8(0)
	for ; c < 64 && v&mask == 0; c++ {
		mask >>= 1
	}
	return c
}
