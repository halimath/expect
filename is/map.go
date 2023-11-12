package is

import "github.com/halimath/expect"

// MapContaining expects got to contain key with value val.
func MapContaining[T ~map[K]V, K, V comparable](got T, key K, val V) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
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

// MapOfLen expects got to contain want elements.
func MapOfLen[T ~map[K]V, K comparable, V any](got T, want int) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		gotLen := len(got)
		if gotLen != want {
			t.Errorf("expected %v to have len %d but got %d", got, want, gotLen)
		}
	})
}
