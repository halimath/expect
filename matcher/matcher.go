// Package matcher defines the interface and a convenience type for implementing matchers. Both built-in and
// custom matchers are implemented by providing a type that implements Matcher.
package matcher

import "testing"

// Matcher defines the generic interface for matchers.
type Matcher[G any] interface {
	// Match is called with a testing interface and the given value. Implementations should perform matching
	// logic and call either t.Error or t.Fatal to report test failures.
	Match(t testing.TB, got G)
}

// Func is a convenience type that allows matchers to be implemented using a single function.
type Func[G any] func(t testing.TB, got G)

func (f Func[G]) Match(t testing.TB, got G) {
	f(t, got)
}
