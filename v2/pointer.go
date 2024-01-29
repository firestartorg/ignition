package ignition

func Dereference[V any](val *V) V {
	if val == nil {
		var zero V
		return zero
	}
	return *val
}

func Reference[V any](val V) *V {
	return &val
}
