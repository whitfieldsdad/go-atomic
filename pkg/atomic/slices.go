package atomic

func concatSlices[T any](slices [][]T) []T {
	var sz int
	for _, s := range slices {
		sz += len(s)
	}
	result := make([]T, sz)
	for i, s := range slices {
		i += copy(result[i:], s)
	}
	return result
}
