package expect

import (
	"reflect"
	"testing"
)

func TestEqual(t *testing.T) {
	var tm tbMock

	IsEqualTo("foo").Match(&tm, "foo")
	IsEqualTo("foo").Match(&tm, "bar")

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{"values are not equal\nwant: foo\ngot:  bar"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
