package string

import (
	"fmt"
)

// stringify converts a value to its string representation.
func Stringify(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", v)
	default:
		return ""
	}
}
