package conv

func Pointer[T any](v any) (*T, bool) {
	if v == nil {
		return nil, true
	}
	data, ok := v.(*T)
	return data, ok
}
