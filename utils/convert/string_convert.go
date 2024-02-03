package convert

import (
	"fmt"
	"strconv"
)

// String to int
func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to int: %v", err)
	}
	return i, nil
}

// Int to string
func IntToString(value int) string {
	str := strconv.Itoa(value)
	return str
}
