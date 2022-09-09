package expect

import "reflect"

// Len expects a values length using the builtin len. It supports strings, slices, arrays, channels and maps.
func Len(want int) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()

		var l int

		v := reflect.ValueOf(got)

		switch v.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
			l = v.Len()
		default:
			ctx.Failf("unable to determine len of <%v>", got)
			return
		}

		if l != want {
			ctx.Failf("exected <%v> to have len %d but got %d", got, want, l)
		}
	})
}
