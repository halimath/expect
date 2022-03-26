package asserthat_test

import (
	"testing"

	"github.com/halimath/assertthat-go/assert"
	"github.com/halimath/assertthat-go/is"
	"github.com/halimath/assertthat-go/matcher"
)

func TestEqual(t *testing.T) {
	got := "something"

	assert.That(t, got, is.Equal("something"))
}

func TestDeepEqual(t *testing.T) {
	got := []string{"foo", "bar"}

	assert.That(t, got, is.DeepEqual([]string{"foo", "bar"}))
}

type Mod interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func IsEven[M Mod]() matcher.Matcher[M] {
	return matcher.Func[M](func(t testing.TB, got M) {
		if got%2 != 0 {
			t.Errorf("expected %v to be even", got)
		}
	})
}

func IsDivisableBy[M Mod](d M) matcher.Matcher[M] {
	return matcher.Func[M](func(t testing.TB, got M) {
		if got%d != 0 {
			t.Errorf("expected %d to be divisable by %d", got, d)
		}
	})
}

func TestCustomMatcher(t *testing.T) {
	var i int = 22
	assert.That(t, i, IsEven[int]())
}

func TestCustomMatcher2(t *testing.T) {
	var i int = 22
	assert.That(t, i, IsDivisableBy(2))
}
