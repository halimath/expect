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

func TestLen_string(t *testing.T) {
	got := "hello"
	expect.That(t, got).
		Is(expect.Len(5))
}

func TestMap(t *testing.T) {
	got := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	expect.That(t, got).
		Is(expect.Len(2)).
		Is(expect.MapContaining("foo", 1))
}

func TestSlice(t *testing.T) {
	got := []int{1, 2, 3}
	expect.That(t, got).
		Is(expect.Len(3)).
		Is(expect.SliceContainingInOrder(1, 3))
}

type Mod interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func Even[M Mod]() expect.Matcher {
	return expect.MatcherFunc(func(ctx expect.Context, got any) {
		g, ok := got.(M)
		if !ok {
			ctx.Failf("expected <%v> to be of type <%T>", got, g)
			return
		}

		if g%2 != 0 {
			ctx.Failf("expected <%v> to be even", got)
		}
	})
}

func DivisableBy[M Mod](d M) expect.Matcher {
	return expect.MatcherFunc(func(ctx expect.Context, got any) {
		g, ok := got.(M)
		if !ok {
			ctx.Failf("expected <%v> to be of type <%T>", got, g)
			return
		}

		if g%d != 0 {
			ctx.Failf("expected <%v> to be divisable by <%v>", got, d)
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
