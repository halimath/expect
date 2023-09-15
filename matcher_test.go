package expect

import (
	"testing"
)

func TestFunc(t *testing.T) {
	m := MatcherFunc[string](func(TB, string) {
	})

	ctx := &tbMock{}
	m.Match(ctx, "got")
}
