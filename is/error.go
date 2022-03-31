package is

import (
	"errors"

	"github.com/halimath/assertthat-go/matcher"
)

func Error(target error) matcher.Matcher[error] {
	return matcher.Func[error](func(t matcher.T, got error) {
		if got == nil {
			t.Errorf("expected an error with target %v but got nil", target)
		} else if !errors.Is(got, target) {
			t.Errorf("expected an error with target %v but got %v", target, got)
		}
	})
}

func NoError() matcher.Matcher[error] {
	return matcher.Func[error](func(t matcher.T, got error) {
		if got != nil {
			t.Errorf("expected no error but got %v", got)
		}
	})
}
