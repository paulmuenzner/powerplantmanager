package data

// flattenMap recursively flattens a nested map
func FlattenMap(inputMap map[string]interface{}) map[string]interface{} {
	flatMap := make(map[string]interface{})

	for key, value := range inputMap {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			// If the value is a nested map, recursively flatten it
			nestedFlatMap := FlattenMap(nestedMap)
			// Append the keys of the nested flat map with a prefix
			for nestedKey, nestedValue := range nestedFlatMap {
				flatMap[key+"."+nestedKey] = nestedValue
			}
		} else {
			// If the value is not a nested map, add it to the flat map
			flatMap[key] = value
		}
	}

	return flatMap
}
