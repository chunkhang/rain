package main

import (
	"math/rand"
)

// RandInt generates a random integer inclusive of min and max
func RandInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
