package expect

import "strings"

// SliceContaining expects got to be a string containing want as a substring.
func IsStringContaining(want string) Matcher[string] {
	return MatcherFunc[string](func(t TB, got string) {
		t.Helper()

		if !strings.Contains(got, want) {
			t.Errorf("expected '%s' to contain '%s'", got, want)
		}
	})
}

// IsStringWithPrefix expects got to be a string having prefix want.
func IsStringWithPrefix(want string) Matcher[string] {
	return MatcherFunc[string](func(t TB, got string) {
		t.Helper()

		if !strings.HasPrefix(got, want) {
			t.Errorf("expected '%s' to have prefix '%s'", got, want)
		}
	})
}

// IsStringWithSuffix expects got to be a string having suffix want.
func IsStringWithSuffix(want string) Matcher[string] {
	return MatcherFunc[string](func(t TB, got string) {
		t.Helper()

		if !strings.HasSuffix(got, want) {
			t.Errorf("expected '%s' to have suffix '%s'", got, want)
		}
	})
}
