package matcher

import (
	"testing"
)

func TestFunc(t *testing.T) {
	m := Func[string](func(t T, got string) {
	})

	m.Match(t, "got")
}
