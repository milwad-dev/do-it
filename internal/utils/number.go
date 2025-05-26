package utils

import "math/rand"

// NumberBetween => generate random number between two numbers
func NumberBetween(min, max int) int {
	return rand.Intn(max-min) + min
}
