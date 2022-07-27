package expect

import (
	"reflect"
	"testing"
)

func TestMapContaining(t *testing.T) {
	var tm contextMock

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapContaining("a", 1).Match(&tm, m)
	MapContaining("a", 2).Match(&tm, m)
	MapContaining("z", 5).Match(&tm, m)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{
			"expected <map[a:1 b:2 c:3]> to contain key <a> with value <2> but got <1>",
			"expected <map[a:1 b:2 c:3]> to contain key <z> but that key does not exist",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
