package expect

import (
	"reflect"
	"testing"
)

func TestEqual(t *testing.T) {
	var tm contextMock

	Equal("foo").Match(&tm, "foo")
	Equal("foo").Match(&tm, "bar")

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"values are not equal\nwant: foo\ngot:  bar"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestDeepEqual(t *testing.T) {
	var tm contextMock

	DeepEqual("foo").Match(&tm, "foo")
	DeepEqual("foo").Match(&tm, "bar")

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"\nvalues are not deeply equal\n\nwant: foo\ngot:  bar\n      ▲▲▲\n"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
