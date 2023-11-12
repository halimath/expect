package is

import (
	"reflect"
	"testing"

	"github.com/halimath/expect/internal/testhelper"
)

func TestEqualTo(t *testing.T) {
	var tb testhelper.TB

	EqualTo("foo", "foo").Expect(&tb)
	EqualTo("bar", "foo").Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs:    []string{"values are not equal\nwant: foo\ngot:  bar"},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}
