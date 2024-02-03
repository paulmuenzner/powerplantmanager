package validate

// ValueEvaluator is a struct for evaluating input values
type ValueEvaluator struct {
	value  interface{}
	errors []string
}

// Validate creates a new ValueEvaluator instance with options.
func Validate(value interface{}) *ValueEvaluator {
	return &ValueEvaluator{value: value, errors: []string{}}
}

// GetResult returns the evaluation result.
func (ve *ValueEvaluator) GetResult() []string {
	if !ve.IsValid() {
		return ve.errors
	}
	return []string{}
}

func (ve *ValueEvaluator) IsValid() bool {
	return len(ve.errors) == 0
}

// CustomErrorMessage returns a formatted custom error message or an empty array if not provided.
func (ve *ValueEvaluator) CustomErrorMessage(customError []string, defaultMessage string) []string {
	if len(customError) > 0 {
		return customError
	}

	return []string{defaultMessage}
}
