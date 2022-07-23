package expect_test

import (
	"testing"

	"github.com/halimath/expect-go"
)

func TestEqual(t *testing.T) {
	got := "something"

	expect.That(t, got).
		Is(expect.Equal("something"))
}

func TestEqual_failNow(t *testing.T) {
	got := "something"

	expect.That(t, got, expect.StopImmediately{}).
		Is(expect.Equal("something"))
}

func TestDeepEqual(t *testing.T) {
	got := []string{"foo", "bar"}

	expect.That(t, got).
		Is(expect.DeepEqual([]string{"foo", "bar"}))
}

type Mod interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func Even[M Mod]() expect.Matcher[M] {
	return expect.MatcherFunc[M](func(ctx expect.Context, got M) {
		if got%2 != 0 {
			ctx.Failf("expected %v to be even", got)
		}
	})
}

func DivisableBy[M Mod](d M) expect.Matcher[M] {
	return expect.MatcherFunc[M](func(ctx expect.Context, got M) {
		if got%d != 0 {
			ctx.Failf("expected %d to be divisable by %d", got, d)
		}
	})
}

func TestCustomMatcher(t *testing.T) {
	var i int = 22
	expect.That(t, i).Is(Even[int]())
}

func TestCustomMatcher2(t *testing.T) {
	var i int = 22
	expect.That(t, i).Is(DivisableBy(2))
}
