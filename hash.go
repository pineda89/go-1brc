package main

import "github.com/cespare/xxhash"

func hash(input []byte) uint64 {
	return xxhash.Sum64(input)
}
