package is

import (
	"github.com/halimath/expect"
	"github.com/halimath/expect/internal/set"
)

// SliceOfLen create an expect.Expectation that expects len(v) == want.
func SliceOfLen[T ~[]S, S any](v T, want int) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		got := len(v)
		if got != want {
			t.Errorf("expected slice with len %d but got slice with len %d: %v", want, got, v)
		}
	})
}

// SliceContaining expects got to be a slice of element type T contain all values given as wants in any order.
// Duplicates in wants are not considered to be contained multiple times in the given slice.
func SliceContaining[S ~[]T, T comparable](v S, wants ...T) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		if len(wants) == 0 {
			return
		}

		wantsMissing := set.New(wants...)

		for _, g := range v {
			wantsMissing.Remove(g)
			if len(wantsMissing) == 0 {
				return
			}
		}

		t.Errorf("%T does not contain %v", v, wantsMissing.ToSlice())
	})
}

// SliceContainingInOrder expects go to be a slice with element type T containing all values given as wants
// in the same order they are given as wants.
func SliceContainingInOrder[S ~[]T, T comparable](v S, wants ...T) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		if len(wants) == 0 {
			return
		}

		for _, g := range v {
			if g == wants[0] {
				wants = wants[1:]
				if len(wants) == 0 {
					return
				}
			}
		}

		t.Errorf("%T does not contain %v in order", v, wants[0])
	})
}
