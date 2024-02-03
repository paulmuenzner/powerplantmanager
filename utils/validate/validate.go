package validate

import (
	"fmt"
	"github.com/paulmuenzner/powerplantmanager/config"
)

// /////////////////////////////////////////////////////////////////////////////
// ///////////////////////////// IsEmail ///////////////////////////////////////
//
// IsEmail checks if the value, when converted to a string, is in email format.
func (ve *ValueEvaluator) IsEmail(customError ...string) *ValueEvaluator {
	strValue, ok := ve.value.(string)
	// Validate if type of string
	if !ok {
		ve.errors = append(ve.errors, "Currently, request cannot be processed.")
		return ve
	}

	emailRegex := config.Regex.Email

	if !emailRegex.MatchString(strValue) {
		customErrorMessage := ve.CustomErrorMessage(customError, "Invalid email format.")
		ve.errors = append(ve.errors, customErrorMessage...)
	}
	return ve
}

// /////////////////////////////////////////////////////////////////////////////
// ///////////////////////////// HasMapExactKeys ///////////////////////////////
//
// HasMapExactKeys checks if a map has the exact number and name of keys in first hierarchy layer
func (ve *ValueEvaluator) HasMapExactKeys(expectedKeys []string, customError ...string) *ValueEvaluator {
	data, ok := ve.value.(map[string]interface{})
	// Validate if type correct
	if !ok {
		ve.errors = append(ve.errors, "Currently, request cannot be processed.")
		return ve
	}

	// Check if the number of keys matches the expected number
	if len(data) != len(expectedKeys) {
		customErrorMessage := ve.CustomErrorMessage(customError, "Data input not correct.")
		ve.errors = append(ve.errors, customErrorMessage...)
		return ve
	}

	// Check if all expected keys are present in the map
	for _, key := range expectedKeys {
		if _, ok := data[key]; !ok {
			customErrorMessage := ve.CustomErrorMessage(customError, "Data input not correct.")
			ve.errors = append(ve.errors, customErrorMessage...)
			return ve
		}
	}
	return ve
}

// /////////////////////////////////////////////////////////////////////////////
// ///////////////////////////// MaxLength /////////////////////////////////////
//
// MaxLength checks if the string representation of the value exceeds a maximum length.
func (ve *ValueEvaluator) MaxLength(maxLength int, customError ...string) *ValueEvaluator {
	strValue, ok := ve.value.(string)
	// Validate if type of string
	if !ok {
		ve.errors = append(ve.errors, "Currently, request cannot be processed.")
		return ve
	}

	if len(strValue) > maxLength {
		customErrorMessage := ve.CustomErrorMessage(customError, fmt.Sprintf("Exceeds maximum length of %d", maxLength))
		ve.errors = append(ve.errors, customErrorMessage...)
	}
	return ve
}

// /////////////////////////////////////////////////////////////////////////////
// ///////////////////////////// MinLength /////////////////////////////////////
//
// MinLength checks if the string representation of the value exceeds a minimum length.
func (ve *ValueEvaluator) MinLength(minLength int, customError ...string) *ValueEvaluator {
	strValue, ok := ve.value.(string)
	// Validate if type of string
	if !ok {
		ve.errors = append(ve.errors, "Currently, request cannot be processed.")
		return ve
	}

	if len(strValue) < minLength {
		customErrorMessage := ve.CustomErrorMessage(customError, fmt.Sprintf("Has not minimum length of %d", minLength))
		ve.errors = append(ve.errors, customErrorMessage...)
	}
	return ve
}

// /////////////////////////////////////////////////////////////////////////////
// ///////////////////////////// MaxIntValue ///////////////////////////////////
//
// MaxIntValue checks if the value is an integer and does not exceed a maximum value
func (ve *ValueEvaluator) MaxIntValue(maxValue int, customError ...string) *ValueEvaluator {
	switch v := ve.value.(type) {
	case int:
		if v > maxValue {
			customErrorMessage := ve.CustomErrorMessage(customError, fmt.Sprintf("Exceeds maximum value of %d", maxValue))
			ve.errors = append(ve.errors, customErrorMessage...)
			return ve
		}
	default:
		ve.errors = append(ve.errors, "Currently, request cannot be processed.")
	}
	return ve
}

// /////////////////////////////////////////////////////////////////////////////
// /////////////////////////// IsValidIPList ///////////////////////////////////
//
// IsValidIPList checks if the slice of strings contains valid IPv6 or IPv4 addresses.
func (ve *ValueEvaluator) IsValidIPList(customError ...string) *ValueEvaluator {
	ipv4Regex := config.Regex.Ipv4
	ipv6Regex := config.Regex.Ipv6
	switch ipList := ve.value.(type) {
	case []interface{}:
		for _, ipStr := range ipList {
			if !ipv4Regex.MatchString(ipStr.(string)) && !ipv6Regex.MatchString(ipStr.(string)) {
				customErrorMessage := ve.CustomErrorMessage(customError, fmt.Sprintf("Invalid IP address: %s", ipStr))
				ve.errors = append(ve.errors, customErrorMessage...)
				return ve // Return after adding error and stop chain processing
			}
		}
	default:
		ve.errors = append(ve.errors, "Currently, request cannot be processed.")
	}
	return ve
}

// /////////////////////////////////////////////////////////////////////////////
// ////////////////////// MaxStringArrayLength /////////////////////////////////
//
// Validate length of array of strings
func (ve *ValueEvaluator) MaxStringArrayLength(maxLength int, customError ...string) *ValueEvaluator {
	switch stringsArray := ve.value.(type) {
	case []string:
		if len(stringsArray) > maxLength {
			customErrorMessage := ve.CustomErrorMessage(customError, fmt.Sprintf("Exceeds maximum array length of %d", maxLength))
			ve.errors = append(ve.errors, customErrorMessage...)
		}
	case []interface{}:
		// Check if each element is a string and update length accordingly
		var length int
		for _, element := range stringsArray {
			if _, ok := element.(string); ok {
				length++
			}
		}
		if length > maxLength {
			customErrorMessage := ve.CustomErrorMessage(customError, fmt.Sprintf("Exceeds maximum array length of %d", maxLength))
			ve.errors = append(ve.errors, customErrorMessage...)
		}
	default:
		ve.errors = append(ve.errors, "Currently, request cannot be processed.")
	}
	return ve
}

// /////////////////////////////////////////////////////////////////////////////
// /////////////////////////// IsPasswordValid /////////////////////////////////
//
// Check valid password format
func (ve *ValueEvaluator) IsPasswordValid(customError ...string) *ValueEvaluator {
	strValue, ok := ve.value.(string)
	// Validate if type of string
	if !ok {
		ve.errors = append(ve.errors, "Currently, request cannot be processed.")
		return ve
	}

	passwordRegex := config.Regex.Password

	if !passwordRegex.MatchString(strValue) {
		customErrorMessage := ve.CustomErrorMessage(customError, "Invalid password format. Only alphanumeric values and special chars allowed.")
		ve.errors = append(ve.errors, customErrorMessage...)
	}
	return ve
}
