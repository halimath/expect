package expect

import (
	"reflect"
	"testing"
)

func TestStringContaining(t *testing.T) {
	var tm tbMock

	IsStringContaining("oba").Match(&tm, "foobar")
	IsStringContaining("oba").Match(&tm, "spameggs")

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{
			"expected 'spameggs' to contain 'oba'",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestStringWithPrefix(t *testing.T) {
	var tm tbMock

	IsStringWithPrefix("foo").Match(&tm, "foobar")
	IsStringWithPrefix("foo").Match(&tm, "spameggs")

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{
			"expected 'spameggs' to have prefix 'foo'",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestStringWithSuffix(t *testing.T) {
	var tm tbMock

	IsStringWithSuffix("bar").Match(&tm, "foobar")
	IsStringWithSuffix("bar").Match(&tm, "spameggs")

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{
			"expected 'spameggs' to have suffix 'bar'",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
