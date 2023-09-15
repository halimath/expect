package expect

import "github.com/halimath/expect-go/internal/set"

// IsSliceContaining expects got to be a slice of element type T contain all values given as wants in any order.
// Duplicates in wants are not considered to be contained multiple times in the given slice.
func IsSliceContaining[T comparable](wants ...T) Matcher[[]T] {
	return MatcherFunc[[]T](func(t TB, got []T) {
		t.Helper()

		if len(wants) == 0 {
			return
		}

		wantsMissing := set.New(wants...)

		for _, g := range got {
			wantsMissing.Remove(g)
			if len(wantsMissing) == 0 {
				return
			}
		}

		t.Errorf("%T does not contain %v", got, wantsMissing.ToSlice())
	})
}

// IsSliceContainingInOrder expects go to be a slice with element type T containing all values given as wants
// in the same order they are given as wants.
func IsSliceContainingInOrder[T comparable](wants ...T) Matcher[[]T] {
	return MatcherFunc[[]T](func(t TB, got []T) {
		t.Helper()

		if len(wants) == 0 {
			return
		}

		for _, g := range got {
			if g == wants[0] {
				wants = wants[1:]
				if len(wants) == 0 {
					return
				}
			}
		}

		t.Errorf("%T does not contain %v in order", got, wants[0])
	})
}
