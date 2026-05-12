package snippet

func ConvertSlice[T any, U any](slice []T) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = any(v).(U)
	}
	return result
}
