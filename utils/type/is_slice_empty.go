package typepackage

import (
	"reflect"
)

// Check if slice is empty
func IsSliceEmpty[T any](data []T) bool {
	// Check if the provided data is a slice
	val := reflect.ValueOf(data)
	length := val.Len()
	return length == 0
}
