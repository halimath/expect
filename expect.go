// Package expect provides functions and types to write fluent, readable expectations for use in unit tests.
// Package expect provides an easy to use and easy to extend framework to write expectations using common
// comparisons as well as providing custom expcters.
package expect

import (
	"fmt"
)

// Expectation defines an interface for types that perform an expectation.
type Expectation interface {
	Expect(t TB)
}

// ExpectFunc is a convenience type to satisfy Expection with a bare function.
type ExpectFunc func(TB)

func (f ExpectFunc) Expect(t TB) {
	t.Helper()
	f(t)
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

// Expectations implements a context for running expectations.
type Expectations struct {
	t TB
}

// Using creates a new Expectations value using t to interact with the test runner.
func Using(t TB) *Expectations {
	t.Helper()
	return &Expectations{t: t}
}

// That runs all expecters using got as the actual value. You may run That multiple times and with the same
// or different value for got.
func (e *Expectations) That(expectations ...Expectation) *Expectations {
	e.t.Helper()

	for _, expecter := range expectations {
		expecter.Expect(e.t)
	}

	return e
}

// WithMessage creates a new Expectations value that prefixes all messages with the message produced by
// applying args to format.
func (e *Expectations) WithMessage(format string, args ...any) *Expectations {
	e.t.Helper()

	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}

	return Using(&prefixedTB{TB: e.t, prefix: format + ": "})
}

type prefixedTB struct {
	TB
	prefix string
}

func (p *prefixedTB) args(args []any) []any {
	p.TB.Helper()

	a := make([]any, len(args)+1)
	copy(a[1:], args)
	a[0] = p.prefix
	return a
}

func (p *prefixedTB) format(format string, args []any) string {
	p.TB.Helper()

	return fmt.Sprint(p.prefix, fmt.Sprintf(format, args...))
}

func (p *prefixedTB) Error(args ...any) {
	p.TB.Helper()
	p.TB.Error(p.args(args)...)
}

func (p *prefixedTB) Errorf(format string, args ...any) {
	p.TB.Helper()
	p.TB.Errorf(p.format(format, args))
}

func (p *prefixedTB) Fatal(args ...any) {
	p.TB.Helper()
	p.TB.Fatal(p.args(args)...)
}

func (p *prefixedTB) Fatalf(format string, args ...any) {
	p.TB.Helper()
	p.TB.Fatalf(p.format(format, args))
}

func (p *prefixedTB) Log(args ...any) {
	p.TB.Helper()
	p.TB.Log(p.args(args)...)
}

func (p *prefixedTB) Logf(format string, args ...any) {
	p.TB.Helper()
	p.TB.Logf(p.format(format, args))
}

func (p *prefixedTB) Skip(args ...any) {
	p.TB.Helper()
	p.TB.Log(p.args(args)...)
}

func (p *prefixedTB) Skipf(format string, args ...any) {
	p.TB.Helper()
	p.TB.Skipf(p.format(format, args))
}

// That is a convenience function that runs all expecters on got reporting to t.
func That(t TB, expectations ...Expectation) {
	t.Helper()
	Using(t).That(expectations...)
}

// WithMessage is a convenience function to create an Expectations values with a pre-set prefix.
func WithMessage(t TB, format string, args ...any) *Expectations {
	t.Helper()
	return Using(t).WithMessage(format, args...)
}

// Fail is an Expectation that always fails.
var Fail Expectation = ExpectFunc(func(t TB) { t.Error("test failed") })

// FailNow is a decorator for an Expectation that converts calls to t.Error, t.Errorf and t.Fail to
// corresponding calls of t.Fatal, t.Fatalf, t.FailNow thus causing the test to fail immediately.
func FailNow(expectations ...Expectation) Expectation {
	return ExpectFunc(func(t TB) {
		wrapped := &failNowTB{TB: t}
		for _, e := range expectations {
			e.Expect(wrapped)
		}
	})
}

type failNowTB struct {
	TB
}

func (f *failNowTB) Error(args ...any) {
	f.TB.Fatal(args...)
}

func (f *failNowTB) Errorf(format string, args ...any) {
	f.TB.Fatalf(format, args...)
}

func (f *failNowTB) Fail() {
	f.TB.FailNow()
}
