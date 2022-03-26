package is

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/halimath/assertthat-go/matcher"
)

// DeepEqual asserts that given and wanted value are deeply equal. It uses the fantastic
// github.com/go-test/deep library which internally uses reflection to compare for equality.
func DeepEqual[G any](want G) matcher.Matcher[G] {
	return matcher.Func[G](func(t testing.TB, got G) {
		if diff := deep.Equal(want, got); diff != nil {
			t.Error(diff)
		}
	})
}

// Equal asserts that given and wanted are equal in terms of the go equality operator. Thus, it works only on
// types that satisfy comparable.
func Equal[G comparable](want G) matcher.Matcher[G] {
	return matcher.Func[G](func(t testing.TB, got G) {
		if want != got {
			t.Errorf("expected %v to equal %v", want, got)
		}
	})
}
