package expect

import "reflect"

// IsNil matches a value to be nil.
func IsNil() Matcher[any] {
	return MatcherFunc[any](func(t TB, got any) {
		t.Helper()

		if got != nil {
			v := reflect.ValueOf(got)
			if v.Kind() == reflect.Pointer {
				v = v.Elem()
			}

			t.Errorf("expected <%v> to be nil", v)
		}
	})
}

// IsNotNil expects got to be non nil.
func IsNotNil() Matcher[any] {
	return MatcherFunc[any](func(t TB, got any) {
		t.Helper()

		if got == nil {
			t.Errorf("expected value to be not nil")
		}
	})
}
