package expect

import (
	"errors"
)

func Error(target error) Matcher[error] {
	return MatcherFunc[error](func(ctx Context, got error) {
		if got == nil {
			ctx.Failf("expected an error with target %v but got nil", target)
		} else if !errors.Is(got, target) {
			ctx.Failf("expected an error with target %v but got %v", target, got)
		}
	})
}

func NoError() Matcher[error] {
	return MatcherFunc[error](func(ctx Context, got error) {
		if got != nil {
			ctx.Failf("expected no error but got %v", got)
		}
	})
}
