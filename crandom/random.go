// Package crandom provides methods to generate random numbers and strings
package crandom

import (
	"crypto/rand"
	"math"
	"math/big"
)

// GenerateRandomNumericalCode generates a positive random number of the given length.
func GenerateRandomNumericalCode(codeLen uint) uint64 {
	if codeLen < 1 {
		panic("codeLen must be >= 1")
	}

	start := uint(math.Pow10(int(codeLen - 1)))
	end := uint(math.Pow10(int(codeLen)) - 1)

	return GenerateRandomNumberBetween(start, end)
}

// GenerateRandomNumberBetween generates a random positive number between the given range.
func GenerateRandomNumberBetween(start, end uint) uint64 {
	return (randInt() % uint64(end-start)) + uint64(start)
}

// GenerateRandomString generates a random string of the given length.
func GenerateRandomString(n uint) string {
	var (
		letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
		b           = make([]rune, n)
	)

	for i := range b {
		b[i] = letterRunes[randIntWithMax(int64(len(letterRunes)))]
	}

	return string(b)
}

func randInt() uint64 {
	return randIntWithMax(math.MaxInt64)
}

func randIntWithMax(max int64) uint64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}

	return nBig.Uint64()
}
