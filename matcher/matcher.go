// Package matcher defines the interface and a convenience type for implementing matchers. Both built-in and
// custom matchers are implemented by providing a type that implements Matcher.
package matcher

// T defines the interface needed by matchers. It is actually a stripped down version of testing.TB which
// only requires what we need.
type T interface {
	// Errorf reports an an error but does not stop test execution.
	Errorf(msg string, args ...any)

	// Fatalf reports an error and stops execution of the test.
	Fatalf(msg string, args ...any)
}

// Matcher defines the generic interface for matchers.
type Matcher[G any] interface {
	// Match is called with a testing interface and the given value. Implementations should perform matching
	// logic and call either t.Error or t.Fatal to report test failures.
	Match(t T, got G)
}

// Func is a convenience type that allows matchers to be implemented using a single function.
type Func[G any] func(t T, got G)

func (f Func[G]) Match(t T, got G) {
	f(t, got)
}
