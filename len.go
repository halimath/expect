package expect

import "reflect"

// HasLenOf expects a values length using the builtin len. It supports strings, slices, arrays, channels and maps.
func HasLenOf[T any](want int) Matcher[T] {
	return MatcherFunc[T](func(t TB, got T) {
		t.Helper()

		var l int

		v := reflect.ValueOf(got)

		switch v.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
			l = v.Len()
		default:
			t.Errorf("unable to determine len of <%v>", got)
			return
		}

		if l != want {
			t.Errorf("exected <%v> to have len %d but got %d", got, want, l)
		}
	})
}
