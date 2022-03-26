package assert

import (
	"testing"

	"github.com/halimath/assertthat-go/matcher"
)

// That implements the entry function to apply a set of matchers on a given value.
// That accepts the testing interface, to given value and a variable list of matchers
// to apply. That calls each matchers Match method in the order the matchers are given.
// Thus, using t.Fatal in one matcher will prevent the following matchers from running.
func That[G any](t testing.TB, got G, matcher ...matcher.Matcher[G]) {
	for _, m := range matcher {
		m.Match(t, got)
	}
}
