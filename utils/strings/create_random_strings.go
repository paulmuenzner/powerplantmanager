package string

import (
	"math/rand"
)

func GenerateRandomNumericString(length int) string {

	const numbers = "0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = numbers[rand.Intn(len(numbers))]
	}

	return string(result)
}
