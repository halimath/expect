package expect

import (
	"fmt"
	"strings"
	"testing"
)

func TestThat_noFailure(t *testing.T) {
	var m1Called, m2Called int
	m1 := MatcherFunc(func(Context, interface{}) {
		m1Called++
	})
	m2 := MatcherFunc(func(Context, interface{}) {
		m2Called++
	})

	ExpectThat(t, "foo").Is(m1).Has(m2)
}

type tbmock struct {
	buf strings.Builder
}

func (m *tbmock) Log(args ...any) {
	fmt.Fprintln(&m.buf, args...)
}

func (m *tbmock) Logf(format string, args ...any) {
	fmt.Fprintf(&m.buf, format, args...)
}

func (m *tbmock) Fail()    {}
func (m *tbmock) FailNow() {}

func TestContext(t *testing.T) {
	m := &tbmock{}
	ctx := context{
		t:    m,
		fail: func() {},
	}

	ctx.Fail("test")

	got := m.buf.String()
	want := "expect_test.go:43: test\n"

	if got != want {
		t.Errorf("expected '%s' but got '%s'", want, got)
	}
}
