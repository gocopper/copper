package crandom_test

import (
	"math"
	"testing"

	"github.com/gocopper/copper/crandom"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomNumericalCode(t *testing.T) {
	t.Parallel()

	for codeLen := 1; codeLen <= 10; codeLen++ {
		code := crandom.GenerateRandomNumericalCode(uint(codeLen))
		actualCodeLen := int(math.Log10(float64(code)) + 1)

		assert.Equal(t, codeLen, actualCodeLen)
	}
}

func TestGenerateRandomNumericalCode_InvalidCodeLen(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		crandom.GenerateRandomNumericalCode(0)
	})
}

func TestGenerateRandomNumberBetween(t *testing.T) {
	t.Parallel()

	num := crandom.GenerateRandomNumberBetween(100, 200)

	assert.GreaterOrEqual(t, num, uint64(100))
	assert.Less(t, num, uint64(200))
}

func TestGenerateRandomString(t *testing.T) {
	t.Parallel()

	for stringLen := 1; stringLen <= 10; stringLen++ {
		str := crandom.GenerateRandomString(uint(stringLen))
		actualStrLen := len(str)

		assert.Equal(t, stringLen, actualStrLen)
	}
}
