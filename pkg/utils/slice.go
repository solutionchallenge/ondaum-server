package utils

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
