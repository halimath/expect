package expect

import (
	"reflect"
	"testing"
)

func TestStringContaining(t *testing.T) {
	var tm contextMock

	StringContaining("oba").Match(&tm, "foobar")
	StringContaining("oba").Match(&tm, "spameggs")
	StringContaining("oba").Match(&tm, 17)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{
			"expected 'spameggs' to contain 'oba'",
			"expected value of type string but got int",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestStringWithPrefix(t *testing.T) {
	var tm contextMock

	StringWithPrefix("foo").Match(&tm, "foobar")
	StringWithPrefix("foo").Match(&tm, "spameggs")
	StringWithPrefix("foo").Match(&tm, 17)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{
			"expected 'spameggs' to have prefix 'foo'",
			"expected value of type string but got int",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestStringWithSuffix(t *testing.T) {
	var tm contextMock

	StringWithSuffix("bar").Match(&tm, "foobar")
	StringWithSuffix("bar").Match(&tm, "spameggs")
	StringWithSuffix("bar").Match(&tm, 17)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{
			"expected 'spameggs' to have suffix 'bar'",
			"expected value of type string but got int",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
