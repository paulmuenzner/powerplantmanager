package files

import "fmt"

// Get file size in bytes of any file type by using Go file descriptor as parameter type
func GetSizeOfByteSlice(data interface{}) (int, error) {
	// This check safeguards against operating on a nil value (absence of a value), which could lead to runtime errors or unexpected behavior later in the function.
	if data == nil {
		return 0, fmt.Errorf("input data is nil")
	}

	// Validate data type
	byteSlice, ok := data.([]byte)
	if !ok {
		return 0, fmt.Errorf("input data is not a []byte slice")
	}

	size := len(byteSlice)
	return size, nil
}
