package expect_test

import (
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func Test_Something(t *testing.T) {
	slice := make([]string, 5)
	var err error

	expect.That(t,
		is.NoError(err),
		is.SliceOfLen(slice, 5),
	)
}

func Test_Slices(t *testing.T) {
	slice := []string{"a", "b", "c"}

	expect.That(t,
		is.SliceOfLen(slice, 3),
		is.SliceContaining(slice, "b", "a"),
		is.SliceContainingInOrder(slice, "a", "b"),
	)
}

func TestEqual_string(t *testing.T) {
	got := "something"

	expect.That(t, is.EqualTo(got, "something"))
}

func TestDeepEqual_strings(t *testing.T) {
	got := []string{"foo", "bar"}

	expect.That(t,
		is.DeepEqualTo(got, []string{"foo", "bar"}),
	)
}

func TestLen_string(t *testing.T) {
	got := "hello"
	expect.That(t, is.StringOfLen(got, 5))
}

func TestMap(t *testing.T) {
	got := map[string]int{
		"foo": 1,
		"bar": 2,
	}

	expect.That(t,
		is.MapContaining(got, "foo", 1),
		is.MapOfLen(got, 2),
	)
}

func TestSlice(t *testing.T) {
	got := []int{1, 2, 3}
	expect.That(t,
		is.SliceOfLen(got, 3),
		is.SliceContainingInOrder(got, 1, 3),
	)
}

// The following type definitions can also be used from the golang.org/x/constraints module. We use these
// "copies" here to avoid a build-time dependency just for this demonstrations case.

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

func IsEven[T Integer](got T) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		if got%2 != 0 {
			t.Errorf("expected <%v> to be even", got)
		}
	})
}

func TestCustomMatcher(t *testing.T) {
	var i int = 22
	expect.That(t, IsEven(i))
}
