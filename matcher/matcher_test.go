package matcher

import (
	"testing"
)

func TestFunc(t *testing.T) {
	m := Func[string](func(t testing.TB, got string) {
	})

	m.Match(t, "got")
}
