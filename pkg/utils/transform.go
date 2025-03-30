package utils

func Map[IN any, OUT any](slice []IN, transformer func(val IN) OUT) []OUT {
	results := make([]OUT, len(slice))
	for idx, val := range slice {
		results[idx] = transformer(val)
	}
	return results
}
