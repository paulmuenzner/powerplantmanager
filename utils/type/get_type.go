package typepackage

import "reflect"

// Get type of value x
func GetType(x interface{}) string {
	return reflect.TypeOf(x).String()
}
