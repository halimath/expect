package is

import (
	"strings"

	"github.com/halimath/expect"
)

// StringOfLen expects got to have byte length want.
func StringOfLen(got string, want int) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		gotLen := len(got)
		if gotLen != want {
			t.Errorf("expected %q to have len %d but got %d", got, want, gotLen)
		}
	})
}

// StringContaining expects got to be a string containing want as a substring.
func StringContaining(got, want string) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		if !strings.Contains(got, want) {
			t.Errorf("expected %q to contain %q", got, want)
		}
	})
}

// StringWithPrefix expects got to be a string having prefix want.
func StringWithPrefix(got, want string) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		if !strings.HasPrefix(got, want) {
			t.Errorf("expected %q to have prefix %q", got, want)
		}
	})
}

// StringWithSuffix expects got to be a string having suffix want.
func StringWithSuffix(got, want string) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		if !strings.HasSuffix(got, want) {
			t.Errorf("expected %q to have suffix %q", got, want)
		}
	})
}
