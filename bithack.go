package main

import (
	"math/bits"
)

// https://graphics.stanford.edu/~seander/bithacks.html

func hasZero(v uint64) uint64 {
	// #define haszero(v) (((v) - 0x01010101UL) & ~(v) & 0x80808080UL)
	// adapted to 64-bit
	return (v - 0x0101010101010101) & ^v & 0x8080808080808080
}

func findPosition(v uint64, mask uint64) int {
	if hz := hasZero(v ^ mask); hz != 0 {
		return bits.TrailingZeros64(hz) >> 3
	}
	return -1
}
