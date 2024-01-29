package ignition

func Coalesce[V any](val *V, defVal V) V {
	if val == nil {
		return defVal
	}
	return *val
}
