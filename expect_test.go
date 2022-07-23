package expect

import (
	"testing"
)

func TestThat_noFailure(t *testing.T) {
	var m1Called, m2Called int
	m1 := MatcherFunc[interface{}](func(Context, interface{}) {
		m1Called++
	})
	m2 := MatcherFunc[interface{}](func(Context, interface{}) {
		m2Called++
	})

	That[interface{}](t, "foo").Is(m1).Has(m2)
}
