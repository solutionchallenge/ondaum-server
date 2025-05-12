package utils

func GroupBy[T comparable, S any](slice []S, getter func(val S) T) map[T][]S {
	if len(slice) == 0 {
		return make(map[T][]S)
	}
	results := make(map[T][]S, len(slice)/2)
	groupCounts := make(map[T]int, len(slice)/2)
	for _, val := range slice {
		key := getter(val)
		groupCounts[key]++
	}
	for key, count := range groupCounts {
		results[key] = make([]S, 0, count)
	}
	for _, val := range slice {
		key := getter(val)
		results[key] = append(results[key], val)
	}
	return results
}

func Reduce[T any, R any](slice []T, reducer func(acc R, val T) R, initial R) R {
	if len(slice) == 0 {
		return initial
	}
	acc := initial
	for _, val := range slice {
		acc = reducer(acc, val)
	}
	return acc
}
