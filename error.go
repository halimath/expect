package expect

import (
	"errors"
)

// Error matches go to contain target in its chain. The check is performed
// using errors.Is.
func Error(target error) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		if got == nil {
			ctx.Failf("expected an error with target %v but got nil", target)
			return
		}

		err, ok := got.(error)
		if !ok {
			ctx.Failf("expected an error but got %T", got)
			return
		}

		if !errors.Is(err, target) {
			ctx.Failf("expected an error with target %v but got %v", target, got)
		}
	})
}

// NoError is an alias for Nil testing an error value to be nil.
var NoError = Nil
