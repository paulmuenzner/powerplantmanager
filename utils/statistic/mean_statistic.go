package statistic

import (
	"fmt"
	typepackage "github.com/paulmuenzner/powerplantmanager/utils/type"
	"reflect"
)

// Mean calculates the mean value of a slice of numeric types
func Mean(data []float64) (float64, error) {
	// Check if the provided data is a slice
	isSlice := typepackage.IsSlice(data)
	if !isSlice {
		return 0, fmt.Errorf("error in 'Mean()'. input must be a slice")
	}

	// Check if the slice is not empty
	isSliceEmpty := typepackage.IsSliceEmpty(data)
	if isSliceEmpty {
		return 0, fmt.Errorf("error in 'Mean()'. slice must not be empty. slice length: %d", len(data))
	}

	// Calculate the sum of elements in the slice
	val := reflect.ValueOf(data)
	length := val.Len()
	sum := 0.0
	for i := 0; i < length; i++ {
		element := val.Index(i).Interface()
		switch element := element.(type) {
		case int:
			sum += float64(element)
		case int8:
			sum += float64(element)
		case int16:
			sum += float64(element)
		case int32:
			sum += float64(element)
		case int64:
			sum += float64(element)
		case float32:
			sum += float64(element)
		case float64:
			sum += element
		default:
			return 0, fmt.Errorf("error in 'Mean()'. unsupported type in slice: %T", element)
		}
	}

	// Calculate the mean value
	mean := sum / float64(length)
	return mean, nil
}
