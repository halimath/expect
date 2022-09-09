package expect

func MapContaining[K, V comparable](key K, val V) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		m, ok := got.(map[K]V)
		if !ok {
			ctx.Failf("expected %T but got %T", m, got)
			return
		}

		vg, ok := m[key]
		if !ok {
			ctx.Failf("expected <%v> to contain key <%v> but that key does not exist", m, key)
			return
		}

		if vg != val {
			ctx.Failf("expected <%v> to contain key <%v> with value <%v> but got <%v>", m, key, val, vg)
		}
	})
}
