package expect

// IsEqualTo asserts that given and wanted are equal in terms of the go equality operator. Thus, it works only on
// types that satisfy comparable.
func IsEqualTo[T comparable](want T) Matcher[T] {
	return MatcherFunc[T](func(t TB, got T) {
		t.Helper()
		if want != got {
			t.Errorf("values are not equal\nwant: %v\ngot:  %v", want, got)
		}
	})
}
