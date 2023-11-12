package is

import (
	"reflect"
	"testing"

	"github.com/halimath/expect/internal/testhelper"
)

func TestMapContaining(t *testing.T) {
	var tb testhelper.TB

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapContaining(m, "a", 1).Expect(&tb)
	MapContaining(m, "a", 2).Expect(&tb)
	MapContaining(m, "z", 5).Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"expected <map[a:1 b:2 c:3]> to contain key <a> with value <2> but got <1>",
			"expected <map[a:1 b:2 c:3]> to contain key <z> but that key does not exist",
		},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}

func TestMapOfLen(t *testing.T) {
	var tb testhelper.TB

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapOfLen(m, 1).Expect(&tb)
	MapOfLen(m, 3).Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"expected map[a:1 b:2 c:3] to have len 1 but got 3",
		},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}
