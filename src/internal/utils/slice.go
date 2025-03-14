package utils

func ConvertMapToSlice[T any](m map[string]T) []T {
	var s []T
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
