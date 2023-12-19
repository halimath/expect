package is

import (
	"errors"

	"github.com/halimath/expect"
)

// Error matches got to contain target in its chain. The check is performed
// using errors.Is.
// A special case happens when target is nil. In this case, Error behaves identical to NoError. This improves
// testing convenience when writing table based tests that test on both error and non error conditions.
func Error(got, target error) expect.Expectation {
	if target == nil {
		return NoError(got)
	}

	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		if got == nil {
			t.Errorf("expected an error with target %v but got nil", target)
			return
		}

		if !errors.Is(got, target) {
			t.Errorf("expected an error with target %v but got %v", target, got)
		}
	})
}

// NoError expects v to be nil.
func NoError(v error) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		if v != nil {
			t.Errorf("expected no error but got %q", v)
		}
	})
}
