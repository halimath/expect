// Package expect provides functions and types to write typesafe, fluent expectations for use in unit tests.
// Package expect uses generics to offer a typesafe and easy to use framework which can be extended to support
// custom assertions/expectations by implementing the Matcher interface.
package expect

import (
	"fmt"
	"path/filepath"
	"runtime"
)

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
type Matcher interface {
	// Match is called with a Context interface and the given value got. Implementations should perform
	// matching logic and call one of the fail methods from Context to mark the test as failed.
	Match(ctx Context, got any)
}

// MatcherFunc is a convenience type that allows matchers to be implemented using a single function.
type MatcherFunc func(ctx Context, got any)

func (f MatcherFunc) Match(ctx Context, got any) {
	f(ctx, got)
}

// TB defines the interface for interacting the a running test. It is a simplification of testing.TB defining
// only what is needed to use the expect library.
type TB interface {
	// Log logs a single message produced from args.
	Log(args ...any)
	// Fail marks the test as failed but continues execution.
	Fail()
	// FailNow marks the test as failed and stops execution.
	FailNow()
}

type failFunc func()

// context implements Context. It calls fail to markt the test as failed.
type context struct {
	t    TB
	fail failFunc
}

var _ Context = &context{}

func (ctx *context) Fail(args ...any) {
	loc := sourceLocation()
	msg := fmt.Sprint(args...)
	ctx.t.Log(loc + ": " + msg)
	ctx.fail()
}

func (ctx *context) Failf(format string, args ...any) {
	ctx.Fail(fmt.Sprintf(format, args...))
}

func sourceLocation() string {
	var pc [50]uintptr
	n := runtime.Callers(3, pc[:])
	if n == 0 {
		panic("expect-go: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	f, ok := frames.Next()
	if !ok {
		panic("expect-go: zero callers found")
	}

	return fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
}

type Clause interface {
	clause()
}

type stopImmediately struct{}

func (stopImmediately) clause() {}

func WithStopImmediately() Clause {
	return stopImmediately{}
}

// Chain defines an intermediate type used to chain expectations on a single type. Normally, test code will
// not use this type but instead call the chaining methods to invoce Matchers.
type Chain struct {
	got any
	ctx Context
}

// ExpectThat starts a new expectation chain using got as the value to expect things from. It uses t to interact
// with the running test.
func ExpectThat(t TB, got any, clauses ...Clause) *Chain {
	ctx := &context{
		t:    t,
		fail: t.Fail,
	}

	for _, clause := range clauses {
		if _, ok := clause.(stopImmediately); ok {
			ctx.fail = t.FailNow
		}
	}

	return &Chain{
		got: got,
		ctx: ctx,
	}
}

// That is an alias for ExpectThat intended to be used when the package is not dot imported.
var That = ExpectThat

// Is adds m to e providing a fluent API.
func (e *Chain) Is(m Matcher) *Chain { return e.runMatcher(m) }

// Has adds m to e providing a fluent API.
func (e *Chain) Has(m Matcher) *Chain { return e.runMatcher(m) }

// And adds m to e providing a fluent API.
func (e *Chain) And(m Matcher) *Chain { return e.runMatcher(m) }

// Matches adds m to e providing a fluent API.
func (e *Chain) Matches(m Matcher) *Chain { return e.runMatcher(m) }

func (e *Chain) runMatcher(m Matcher) *Chain {
	m.Match(e.ctx, e.got)
	return e
}
