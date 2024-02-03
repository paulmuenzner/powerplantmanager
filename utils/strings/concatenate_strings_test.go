package string

import (
	"sync"
	"testing"
)

func TestConcatenateStrings(t *testing.T) {
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{"Testing", " ", "goroutine", " ", "concurrency", " ", "now."}, "Testing goroutine concurrency now."},
		{[]string{"Separate", " ", "file", " ", "test"}, "Separate file test"},
		{[]string{"Go", " ", "testing"}, "Go testing"},
		{[]string{"Multiple", " ", "test", " ", "cases", " ", "are", " ", "possible."}, "Multiple test cases are possible."},
		{[]string{"Another", " ", "concurrent", " ", "execution", " ", "example"}, "Another concurrent execution example"},
		{[]string{"Simple", " ", "Golang", " ", "tests"}, "Simple Golang tests"},
		{[]string{"Easy", " ", "to", " ", "understand"}, "Easy to understand"},
		{[]string{"Short", " ", "test", "."}, "Short test."},
	}

	var wg sync.WaitGroup

	for _, testCase := range testCases {
		wg.Add(1)
		go func(tc struct {
			input    []string
			expected string
		}) {
			defer wg.Done()

			result := ConcatenateStrings(tc.input...)
			if result != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, result)
			}
		}(testCase)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
