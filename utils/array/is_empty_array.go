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
	default:
		return false
	}
}
