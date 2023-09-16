// Package expect provides functions and types to write typesafe, fluent expectations for use in unit tests.
// Package expect uses generics to offer a typesafe and easy to use framework which can be extended to support
// custom assertions/expectations by implementing the Matcher interface.
package expect

// Matcher defines the interface implemented to execute a single expectation.
type Matcher[G any] interface {
	// Match is called with a Context interface and the given value got. Implementations should perform
	// matching logic and call one of the fail methods from Context to mark the test as failed.
	Match(t TB, got G)
}

// MatcherFunc is a convenience type that allows matchers to be implemented using a single function.
type MatcherFunc[G any] func(TB, G)

func (f MatcherFunc[G]) Match(t TB, got G) {
	t.Helper()
	f(t, got)
}

// TB is basically a copy from testing.TB. It is used here to allow other implementations (testing.TB
// contains an unexported private method) to be used (i.e. mocks while testing matchers). All methods of this
// interface work exactly the same as their counterparts from testing.TB type.
type TB interface {
	Cleanup(func())
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	Setenv(key, value string)
	Skip(args ...any)
	SkipNow()
	Skipf(format string, args ...any)
	Skipped() bool
	TempDir() string
}

// ExpectThat starts a new expectation chain using got as the value to expect things from. It uses t to interact
// with the running test. Any failing matches originating from a call to ExpectThat cause the test to fail at
// the end of the execution; it works like caling t.Error(...).
func ExpectThat[G any](t TB, got G, matchers ...Matcher[G]) {
	t.Helper()

	for _, m := range matchers {
		m.Match(t, got)
	}
}

// That is an alias for ExpectThat intended to be used when the package is not dot imported.
func That[T any](t TB, got T) { ExpectThat(t, got) }

// EnsureThat starts a new expectation chain using got as the value to ensure things from. It uses t to interact
// with the running test. Any failing matches originating from a call to EnsureThat cause the test to fail
// immediately; it works like caling t.FailNow(...).
func EnsureThat[G any](t TB, got G, matchers ...Matcher[G]) {
	t.Helper()

	ExpectThat(&ensureTB{t}, got, matchers...)
}

type ensureTB struct {
	TB
}

func (t *ensureTB) Error(args ...any)                 { t.Helper(); t.Fatal(args...) }
func (t *ensureTB) Errorf(format string, args ...any) { t.Helper(); t.Fatalf(format, args...) }
func (t *ensureTB) Fail()                             { t.Helper(); t.FailNow() }
