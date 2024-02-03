package string

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomNumericString(t *testing.T) {
	// Define test cases
	testCases := []struct {
		input    int
		expected int
	}{
		{2, 2},
		{3, 3},
		{4, 4},
		{78, 78},
		{777, 777},
		// Add more test cases as needed
	}

	// Run tests concurrently using goroutines
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Stringify(%v)", tc.input), func(t *testing.T) {
			t.Parallel() // indicates that this test can be run concurrently

			result := len(GenerateRandomNumericString(tc.input))
			assert.Equal(t, tc.expected, result, "Unexpected result")
		})
	}
}
