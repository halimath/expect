package expect

import (
	"testing"
)

func TestFunc(t *testing.T) {
	m := MatcherFunc(func(ctx Context, got any) {
	})

	ctx := &contextMock{}
	m.Match(ctx, "got")
}
