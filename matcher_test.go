package expect

import (
	"testing"
)

func TestFunc(t *testing.T) {
	m := MatcherFunc[string](func(ctx Context, got string) {
	})

	ctx := &contextMock{}
	m.Match(ctx, "got")
}
