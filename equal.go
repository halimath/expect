package expect

import "reflect"

// DeepEqual asserts that given and wanted value are deeply equal by using reflection to inspect and dive into
// nested structures.
func DeepEqual[T any](want T) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		if !reflect.DeepEqual(want, got) {
			ctx.Failf("values are not deeply equal: want\n%#v got\n%#v", want, got)
		}
	})
}

// Equal asserts that given and wanted are equal in terms of the go equality operator. Thus, it works only on
// types that satisfy comparable.
func Equal[G comparable](want G) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
		if want != got {
			ctx.Failf("values are not equal: want\n%v got\n%v", want, got)
		}
	})
}
