package utils

func Intersect[T comparable](first, second []T) []T {
	if len(first) == 0 || len(second) == 0 {
		return nil
	}
	firstSet := make(map[T]struct{}, len(first))
	for _, v := range first {
		firstSet[v] = struct{}{}
	}
	result := make([]T, 0, min(len(first), len(second)))
	for _, v := range second {
		if _, exists := firstSet[v]; exists {
			result = append(result, v)
			delete(firstSet, v)
		}
	}
	return result
}

func Deduplicate[T comparable](input []T) []T {
	if len(input) == 0 {
		return nil
	}
	seen := make(map[T]struct{}, len(input))
	unique := make([]T, 0, len(input))
	for _, v := range input {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			unique = append(unique, v)
		}
	}
	return unique
}
