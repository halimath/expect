package expect

import (
	"reflect"
	"testing"
)

func TestNil(t *testing.T) {
	var tm tbMock

	s := "foo"
	IsNil().Match(&tm, &s)
	IsNil().Match(&tm, nil)

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{"expected <foo> to be nil"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestNotNil(t *testing.T) {
	var tm tbMock

	s := "foo"
	IsNotNil().Match(&tm, &s)
	IsNotNil().Match(&tm, nil)

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{"expected value to be not nil"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
