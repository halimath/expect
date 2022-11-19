package expect

// Equal asserts that given and wanted are equal in terms of the go equality operator. Thus, it works only on
// types that satisfy comparable.
func Equal[G comparable](want G) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		ctx.T().Helper()
		if want != got {
			ctx.Failf("values are not equal\nwant: %v\ngot:  %v", want, got)
		}
	})
}
