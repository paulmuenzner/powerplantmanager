package typepackage

import "reflect"

func IsSlice(data interface{}) bool {

	// Check if the provided data is a slice
	val := reflect.ValueOf(data)
	return val.Kind() == reflect.Slice
}
