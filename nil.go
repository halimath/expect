package expect

import "reflect"

// Nil matches a value to be nil.
func Nil() Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		if got != nil {
			v := reflect.ValueOf(got)
			if v.Kind() == reflect.Pointer {
				v = v.Elem()
			}

			ctx.Failf("expected <%v> to be nil", v)
		}
	})
}

// NotNil expects got to be non nil.
func NotNil() Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		if got == nil {
			ctx.Failf("expected value to be not nil")
		}
	})
}
