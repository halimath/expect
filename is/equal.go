package is

import (
	"reflect"

	"github.com/halimath/assertthat-go/matcher"
)

// DeepEqual asserts that given and wanted value are deeply equal by using reflection to inspect and dive into
// nested structures.
func DeepEqual[G any](want G) matcher.Matcher[G] {
	return matcher.Func[G](func(t matcher.T, got G) {
		if !reflect.DeepEqual(want, got) {
			t.Errorf("values are not deeply equal: want\n%#v got\n%#v", want, got)
		}
	})
}

// Equal asserts that given and wanted are equal in terms of the go equality operator. Thus, it works only on
// types that satisfy comparable.
func Equal[G comparable](want G) matcher.Matcher[G] {
	return matcher.Func[G](func(t matcher.T, got G) {
		if want != got {
			t.Errorf("expected %v to equal %v", want, got)
		}
	})
}
