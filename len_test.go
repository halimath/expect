package expect

import (
	"reflect"
	"testing"
)

func TestLen(t *testing.T) {
	var tm tbMock

	HasLenOf[int](2).Match(&tm, 32)

	HasLenOf[string](2).Match(&tm, "fo")
	HasLenOf[string](2).Match(&tm, "foo")

	HasLenOf[[]rune](2).Match(&tm, []rune{'f', 'o'})
	HasLenOf[[]rune](2).Match(&tm, []rune{'f', 'o', 'o'})

	HasLenOf[map[string]int](2).Match(&tm, map[string]int{"f": 1, "o": 2})
	HasLenOf[map[string]int](2).Match(&tm, map[string]int{"f": 1, "o": 2, "b": 3})

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{
			"unable to determine len of <32>",
			"exected <foo> to have len 2 but got 3",
			"exected <[102 111 111]> to have len 2 but got 3",
			"exected <map[b:3 f:1 o:2]> to have len 2 but got 3",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
