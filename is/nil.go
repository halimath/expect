package is

import (
	"github.com/halimath/assertthat-go/matcher"
)

func Nil[T any]() matcher.Matcher[*T] {
	return matcher.Func[*T](func(t matcher.T, got *T) {
		if got != nil {
			t.Errorf("expected %v to be nil", *got)
		}
	})
}

func NotNil[T any]() matcher.Matcher[*T] {
	return matcher.Func[*T](func(t matcher.T, got *T) {
		if got == nil {
			t.Errorf("expected value to be not nil")
		}
	})
}
