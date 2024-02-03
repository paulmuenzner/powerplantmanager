package env

import (
	"fmt"
	"os"
)

// GetEnvValue retrieves the value of an environment variable or returns an error and a default value if not available
func GetEnvValue(key, defaultValue string) (envValue string, err error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, fmt.Errorf("Unable to retrieve .env value for key '%s'. Default value used.", key)
	}
	return value, nil
}
