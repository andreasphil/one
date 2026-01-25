package util

func Filter[T any](list []T, match func(T) bool) (result []T) {
	for _, i := range list {
		if match(i) {
			result = append(result, i)
		}
	}

	return
}
