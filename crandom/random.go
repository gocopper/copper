// Package crandom provides utility random functions to generate random numbers and strings.
package crandom

import (
	"math"
	"math/rand"
	"time"
)

// Seed can be used with fx.Invoke to seed the randomizer.
func Seed() {
	rand.Seed(time.Now().Unix())
}

// GenerateRandomNumericalCode generates a positive random number of the given length.
func GenerateRandomNumericalCode(len uint) uint64 {
	if len == 0 {
		return 0
	}
	start := uint(math.Pow10(int(len - 1)))
	end := uint(math.Pow10(int(len)) - 1)
	return GenerateRandomNumberBetween(start, end)
}

// GenerateRandomNumberBetween generates a random positive number between the given range.
func GenerateRandomNumberBetween(start, end uint) uint64 {
	return (rand.Uint64() % uint64(end-start)) + uint64(start)
}

// GenerateRandomString generates a random string of the given length.
func GenerateRandomString(n uint) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
