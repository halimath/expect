// Package expect provides functions and types to write typesafe, fluent expectations for use in unit tests.
// Package expect uses generics to offer a typesafe and easy to use framework which can be extended to support
// custom assertions/expectations by implementing the Matcher interface.
package expect

// Context defines the context used by Matcher to interact with the running test (i.e. failing it.). Using the
// Context interface Matchers can reject values with a custom message which will fail the test either
// immediately or deferred, depending on the surrounding call.
type Context interface {
	// Fail fails the test producing a message using the provided args.
	Fail(args ...any)

	// Failf fails the test producing a message by formatting args according to format.
	Failf(format string, args ...any)
}

// Matcher defines the interface implemented to execute a single expectation.
type Matcher[T any] interface {
	// Match is called with a Context interface and the given value got. Implementations should perform
	// matching logic and call one of the fail methods from Context to mark the test as failed.
	Match(ctx Context, got T)
}

// MatcherFunc is a convenience type that allows matchers to be implemented using a single function.
type MatcherFunc[T any] func(ctx Context, got T)

func (f MatcherFunc[T]) Match(ctx Context, got T) {
	f(ctx, got)
}

// TB defines the interface for interacting the a running test. It is a simplification of testing.TB defining
// only what is needed to use the expect library.
type TB interface {
	// Log logs a single message produced from args.
	Log(args ...any)
	// Logf logs a single message produced by formatting args according to format.
	Logf(format string, args ...any)
	// Fail marks the test as failed but continues execution.
	Fail()
	// FailNow marks the test as failed and stops execution.
	FailNow()
}

// failContext implements Context with deferred test failure.
type failContext struct {
	t TB
}

var _ Context = &failContext{}

func (ctx *failContext) Fail(args ...any) {
	ctx.t.Log(args...)
	ctx.t.Fail()
}

func (ctx *failContext) Failf(format string, args ...any) {
	ctx.t.Logf(format, args...)
	ctx.t.Fail()
}

// failNowContext implements Context with immediate test failure.
type failNowContext struct {
	t TB
}

var _ Context = &failNowContext{}

func (ctx *failNowContext) Fail(args ...any) {
	ctx.t.Log(args...)
	ctx.t.FailNow()
}

func (ctx *failNowContext) Failf(format string, args ...any) {
	ctx.t.Logf(format, args...)
	ctx.t.FailNow()
}

type Clause interface {
	clause()
}

type StopImmediately struct{}

func (StopImmediately) clause() {}

// Chain defines an intermediate type used to chain expectations on a single type. Normally, test code will
// not use this type but instead call the chaining methods to invoce Matchers.
type Chain[T any] struct {
	got T
	ctx Context
}

// That starts a new expectation chain using got as the value to expect things from. It uses t to interact
// with the running test.
func That[T any](t TB, got T, clauses ...Clause) *Chain[T] {
	var ctx Context = &failContext{t: t}

	for _, clause := range clauses {
		if _, ok := clause.(StopImmediately); ok {
			ctx = &failNowContext{t: t}
		}
	}

	return &Chain[T]{
		got: got,
		ctx: ctx,
	}
}

// Is adds m to e providing a fluent API.
func (e *Chain[T]) Is(m Matcher[T]) *Chain[T] { return e.runMatcher(m) }

// Has adds m to e providing a fluent API.
func (e *Chain[T]) Has(m Matcher[T]) *Chain[T] { return e.runMatcher(m) }

// And adds m to e providing a fluent API.
func (e *Chain[T]) And(m Matcher[T]) *Chain[T] { return e.runMatcher(m) }

// Matches adds m to e providing a fluent API.
func (e *Chain[T]) Matches(m Matcher[T]) *Chain[T] { return e.runMatcher(m) }

func (e *Chain[T]) runMatcher(m Matcher[T]) *Chain[T] {
	m.Match(e.ctx, e.got)
	return e
}
