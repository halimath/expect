package is

import (
	"strings"
	"unicode"

	"github.com/halimath/expect"
)

// StringOfLen expects got to have byte length want.
func StringOfLen(got string, want int) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

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

// Dedent is intended to be used as a transformer passed to [EqualToStringByLines].
// It removes any prefix whitespace from s thus dedenting each line. This is
// especially usefull if the expected value for a test is written in code as an
// indented multiline raw string literal but the actual lines are not indented.
func DedentLines(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

// EqualToStringByLines compares got and want line by line and reports different
// lines one at a time. This makes it easiert to understand failed expectations
// when comparing large strings.
//
// transformers are applied to all lines, both those obtained from got and want.
// transformers are applied in order (iteratively) and the final transformation
// result is used for comparison.
func EqualToStringByLines(got, want string, transformers ...func(string) string) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		gotLines := strings.Split(got, "\n")
		wantLines := strings.Split(want, "\n")

		lenghtsDiffer := len(gotLines) != len(wantLines)

		if lenghtsDiffer {
			t.Errorf("expected string to have %d lines but got %d", len(wantLines), len(gotLines))
		}

		limit := min(len(gotLines), len(wantLines))

		for i := 0; i < limit; i++ {
			gotLine := gotLines[i]
			wantLine := wantLines[i]

			for _, transformer := range transformers {
				gotLine = transformer(gotLine)
				wantLine = transformer(wantLine)
			}

			if gotLine != wantLine {
				t.Errorf("at line %d: wanted\n%q\nbut got\n%q", i, wantLines[i], gotLines[i])
				if lenghtsDiffer {
					return
				}
			}
		}

		if len(gotLines) > limit {
			for i, line := range gotLines[limit:] {
				t.Errorf("line %d: wanted no line but got\n%q", i+limit, line)
			}
		}

		if len(wantLines) > limit {
			for i, line := range wantLines[limit:] {
				t.Errorf("line %d: wanted\n%q\nbut got no line", i+limit, line)
			}
		}

	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
