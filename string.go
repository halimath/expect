package expect

import "strings"

// SliceContaining expects got to be a string containing want as a substring.
func StringContaining(want string) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		s, ok := got.(string)
		if !ok {
			ctx.Failf("expected value of type string but got %T", got)
			return
		}

		if !strings.Contains(s, want) {
			ctx.Failf("expected '%s' to contain '%s'", got, want)
		}
	})
}

// StringWithPrefix expects got to be a string having prefix want.
func StringWithPrefix(want string) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		s, ok := got.(string)
		if !ok {
			ctx.Failf("expected value of type string but got %T", got)
			return
		}

		if !strings.HasPrefix(s, want) {
			ctx.Failf("expected '%s' to have prefix '%s'", got, want)
		}
	})
}

// StringWithSuffix expects got to be a string having suffix want.
func StringWithSuffix(want string) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		s, ok := got.(string)
		if !ok {
			ctx.Failf("expected value of type string but got %T", got)
			return
		}

		if !strings.HasSuffix(s, want) {
			ctx.Failf("expected '%s' to have suffix '%s'", got, want)
		}
	})
}
