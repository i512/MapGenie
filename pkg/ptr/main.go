package ptr

func Take[V any](v V) *V {
	return &v
}
