package arrayhandler

// /////////////////////////////////////////////////////////////////////////////////////////////
// CHANGE PLANT CONFIGURATION VALIDATION
// ///////////////////////
func IsEmptyArray(value interface{}) bool {
	switch v := value.(type) {
	case []interface{}:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case []float64:
		return len(v) == 0
	// Add more cases as needed
	default:
		return false
	}
}
