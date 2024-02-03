package convert

import (
	"fmt"
)

// Error to string
func ErrorToString(err error) string {
	if err == nil {
		return "No error detected."
	}
	return fmt.Sprintf("%v", err)
}
