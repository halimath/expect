package expect

import (
	"reflect"
	"testing"
)

func TestLen(t *testing.T) {
	var tm contextMock

	Len(2).Match(&tm, 32)

	Len(2).Match(&tm, "fo")
	Len(2).Match(&tm, "foo")

	Len(2).Match(&tm, []rune{'f', 'o'})
	Len(2).Match(&tm, []rune{'f', 'o', 'o'})

	Len(2).Match(&tm, map[string]int{"f": 1, "o": 2})
	Len(2).Match(&tm, map[string]int{"f": 1, "o": 2, "b": 3})

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{
			"unable to determine len of <32>",
			"exected <foo> to have len 2 but got 3",
			"exected <[102 111 111]> to have len 2 but got 3",
			"exected <map[b:3 f:1 o:2]> to have len 2 but got 3",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
