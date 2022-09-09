package expect_test

import (
	"testing"

	. "github.com/halimath/expect-go"
)

func TestEqual_string(t *testing.T) {
	got := "something"

	ExpectThat(t, got).
		Is(Equal("something"))
}

func TestEqual_failNow(t *testing.T) {
	got := "something"

	ExpectThat(t, got, WithStopImmediately()).
		Is(Equal("something"))
}

func TestDeepEqual_strings(t *testing.T) {
	got := []string{"foo", "bar"}

	ExpectThat(t, got).
		Is(DeepEqual([]string{"foo", "bar"}))
}

func TestLen_string(t *testing.T) {
	got := "hello"
	ExpectThat(t, got).
		Has(Len(5))
}

func TestMap(t *testing.T) {
	got := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	ExpectThat(t, got).
		Is(MapContaining("foo", 1)).
		Has(Len(2))
}

func TestSlice(t *testing.T) {
	got := []int{1, 2, 3}
	ExpectThat(t, got).
		Has(Len(3)).
		Is(SliceContainingInOrder(1, 3))
}

type Mod interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func Even[M Mod]() Matcher {
	return MatcherFunc(func(ctx Context, got any) {
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

func DivisableBy[M Mod](d M) Matcher {
	return MatcherFunc(func(ctx Context, got any) {
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
	ExpectThat(t, i).Is(Even[int]())
}

func TestCustomMatcher2(t *testing.T) {
	var i int = 22
	ExpectThat(t, i).Is(DivisableBy(2))
}
