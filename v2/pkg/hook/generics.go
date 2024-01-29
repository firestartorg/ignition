package hook

func Gather[T any](h Collection) []T {
	var result []T
	for _, hook := range h.Hooks() {
		if hook, ok := hook.(T); ok {
			result = append(result, hook)
		}
	}
	return result
}
