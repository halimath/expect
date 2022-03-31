package assert

import (
	"testing"

	"github.com/halimath/assertthat-go/matcher"
)

func TestThat_noFailure(t *testing.T) {
	var m1Called, m2Called int
	m1 := matcher.Func[interface{}](func(matcher.T, interface{}) {
		m1Called++
	})
	m2 := matcher.Func[interface{}](func(matcher.T, interface{}) {
		m2Called++
	})

	That[interface{}](t, "foo", m1, m2)
}
