package expect

import (
	"errors"
)

// IsError matches go to contain target in its chain. The check is performed
// using errors.Is.
func IsError(target error) Matcher[error] {
	return MatcherFunc[error](func(t TB, got error) {
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

// NoError is an alias for Nil testing an error value to be nil.
var NoError = IsNil
