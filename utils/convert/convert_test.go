package convert

import (
	"testing"
	"time"
)

// //////////////////////////////////////////////////
// //////////////////////////////////////////////////
// String to int test
// //////////////////
func TestConvertStringToInt(t *testing.T) {
	// Define test cases
	testCases := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"-456", -456},
		{"invalid", 0}, // error
		{"4l5", 0},     // error
		{"4,,5", 0},    // error
		{"42l", 0},     // error
	}

	// Create a channel to collect results from goroutines
	results := make(chan struct {
		index    int
		actual   int
		expected int
		err      error
	}, len(testCases))

	// Run each test case in a goroutine
	for index, tc := range testCases {
		go func(index int, tc struct {
			input    string
			expected int
		}) {
			actual, err := StringToInt(tc.input)
			results <- struct {
				index    int
				actual   int
				expected int
				err      error
			}{index, actual, tc.expected, err}
		}(index, tc)
	}

	// Close the results channel after all goroutines have completed
	go func() {
		time.Sleep(time.Millisecond) // Ensure all goroutines have a chance to execute
		close(results)
	}()

	// Collect and verify results
	for res := range results {
		if res.err != nil {
			// Expecting an error
			if res.expected != 0 {
				t.Errorf("Test case %d: Expected error, but got %v", res.index, res.actual)
			}
		} else {
			// Expecting a valid result
			if res.actual != res.expected {
				t.Errorf("Test case %d: Expected %d, but got %d", res.index, res.expected, res.actual)
			}
		}
	}
}
