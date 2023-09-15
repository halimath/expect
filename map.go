package expect

func IsMapContaining[K, V comparable](key K, val V) Matcher[map[K]V] {
	return MatcherFunc[map[K]V](func(t TB, got map[K]V) {
		t.Helper()

		vg, ok := got[key]
		if !ok {
			t.Errorf("expected <%v> to contain key <%v> but that key does not exist", got, key)
			return
		}

		if vg != val {
			t.Errorf("expected <%v> to contain key <%v> with value <%v> but got <%v>", got, key, val, vg)
		}
	})
}
