package expect

import (
	"reflect"
	"testing"
)

func TestNil(t *testing.T) {
	var tm contextMock

	s := "foo"
	Nil[string]().Match(&tm, &s)
	Nil[any]().Match(&tm, nil)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"expected foo to be nil"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestNotNil(t *testing.T) {
	var tm contextMock

	s := "foo"
	NotNil[string]().Match(&tm, &s)
	NotNil[any]().Match(&tm, nil)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"expected value to be not nil"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
