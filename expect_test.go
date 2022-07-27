package expect

import (
	"testing"
)

func TestThat_noFailure(t *testing.T) {
	var m1Called, m2Called int
	m1 := MatcherFunc(func(Context, interface{}) {
		m1Called++
	})
	m2 := MatcherFunc(func(Context, interface{}) {
		m2Called++
	})

	That(t, "foo").Is(m1).Has(m2)
}
