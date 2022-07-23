package expect

func Nil[T any]() Matcher[*T] {
	return MatcherFunc[*T](func(ctx Context, got *T) {
		if got != nil {
			ctx.Failf("expected %v to be nil", *got)
		}
	})
}

func NotNil[T any]() Matcher[*T] {
	return MatcherFunc[*T](func(ctx Context, got *T) {
		if got == nil {
			ctx.Failf("expected value to be not nil")
		}
	})
}
