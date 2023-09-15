package expect_test

import (
	"testing"

	. "github.com/halimath/expect-go"
)

func TestEqual_string(t *testing.T) {
	got := "something"

	ExpectThat(t, got, IsEqualTo("something"))
}

func TestEqual_ensure(t *testing.T) {
	got := "something"

	EnsureThat(t, got, IsEqualTo("something"))
}

func TestDeepEqual_strings(t *testing.T) {
	got := []string{"foo", "bar"}

	ExpectThat(t, got, IsDeepEqualTo([]string{"foo", "bar"}))
}

func TestLen_string(t *testing.T) {
	got := "hello"
	ExpectThat(t, got, HasLenOf[string](5))
}

func TestMap(t *testing.T) {
	got := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	ExpectThat(t, got, IsMapContaining("foo", 1), HasLenOf[map[string]int](2))
}

func TestSlice(t *testing.T) {
	got := []int{1, 2, 3}
	ExpectThat(t, got, HasLenOf[[]int](3), IsSliceContainingInOrder(1, 3))
}

type Mod interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func IsEven[G Mod]() Matcher[G] {
	return MatcherFunc[G](func(t TB, got G) {
		if got%2 != 0 {
			t.Errorf("expected <%v> to be even", got)
		}
	})
}

func IsDivisableBy[G Mod](d G) Matcher[G] {
	return MatcherFunc[G](func(t TB, got G) {
		if got%d != 0 {
			t.Errorf("expected <%v> to be divisable by <%v>", got, d)
		}
	})
}

func TestCustomMatcher(t *testing.T) {
	var i int = 22
	ExpectThat(t, i, IsEven[int]())
}

func TestCustomMatcher2(t *testing.T) {
	var i int = 22
	ExpectThat(t, i, IsDivisableBy(2))
}
