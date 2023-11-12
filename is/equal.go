package is

import "github.com/halimath/expect"

// EqualTo asserts that given and wanted are equal in terms of the go equality operator. Thus, it works only on
// types that satisfy comparable.
func EqualTo[T comparable](got, want T) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		if want != got {
			t.Errorf("values are not equal\nwant: %v\ngot:  %v", want, got)
		}
	})
}
