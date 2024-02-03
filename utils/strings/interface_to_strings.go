package string

import "fmt"

func InterfaceToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
