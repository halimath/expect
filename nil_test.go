package expect

import (
	"reflect"
	"testing"
)

func TestNil(t *testing.T) {
	var tm contextMock

	s := "foo"
	Nil().Match(&tm, &s)
	Nil().Match(&tm, nil)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"expected <foo> to be nil"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestNotNil(t *testing.T) {
	var tm contextMock

	s := "foo"
	NotNil().Match(&tm, &s)
	NotNil().Match(&tm, nil)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"expected value to be not nil"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
