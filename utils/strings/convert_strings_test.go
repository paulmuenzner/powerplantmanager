package string

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringify(t *testing.T) {
	// Define test cases
	testCases := []struct {
		input    interface{}
		expected string
	}{
		{"test", "test"},
		{42, "42"},
		{3.14, "3.14"},
		{uint(100), "100"},
		{int32(-42), "-42"},
		{float64(2.71828), "2.71828"},
		// Add more test cases as needed
	}

	// Run tests concurrently using goroutines
	for _, tc := range testCases {
		tc := tc // create a local copy of the loop variable for each iteration
		t.Run(fmt.Sprintf("Stringify(%v)", tc.input), func(t *testing.T) {
			t.Parallel() // indicates that this test can be run concurrently

			result := Stringify(tc.input)
			assert.Equal(t, tc.expected, result, "Unexpected result")
		})
	}
}
