package is

import (
	"reflect"
	"testing"
)

func TestEqual(t *testing.T) {
	var tm tMock

	Equal("foo").Match(&tm, "foo")
	Equal("foo").Match(&tm, "bar")

	if !reflect.DeepEqual(tm, tMock{
		errors: []string{"expected foo to equal bar"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestDeepEqual(t *testing.T) {
	var tm tMock

	DeepEqual("foo").Match(&tm, "foo")
	DeepEqual("foo").Match(&tm, "bar")

	if !reflect.DeepEqual(tm, tMock{
		errors: []string{"values are not deeply equal: want\n\"foo\" got\n\"bar\""},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
