package expect

import (
	"testing"
)

func TestThat_noFailure(t *testing.T) {
	var m1Called, m2Called int
	var m1 Matcher[string] = MatcherFunc[string](func(TB, string) {
		m1Called++
	})
	var m2 Matcher[string] = MatcherFunc[string](func(TB, string) {
		m2Called++
	})

	ExpectThat(t, "foo", m1, m2)
}
