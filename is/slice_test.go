package is

import (
	"reflect"
	"testing"

	"github.com/halimath/expect/internal/testhelper"
)

func TestSliceOfLen(t *testing.T) {
	var tb testhelper.TB

	s := make([]string, 3)

	SliceOfLen(s, 1).Expect(&tb)
	SliceOfLen(s, 3).Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs:    []string{"expected slice with len 1 but got slice with len 3: [  ]"},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}

func TestSliceContaining(t *testing.T) {
	var tb testhelper.TB

	s := []int{1, 2, 3, 4}
	SliceContaining(s).Expect(&tb)
	SliceContaining(s, 1, 3).Expect(&tb)
	SliceContaining(s, 3, 1).Expect(&tb)
	SliceContaining(s, 1, 5).Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs:    []string{"[]int does not contain [5]"},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}

func TestSliceContainingInOrder(t *testing.T) {
	var tb testhelper.TB

	s := []int{1, 2, 3, 4}
	SliceContainingInOrder(s).Expect(&tb)
	SliceContainingInOrder(s, 1, 3).Expect(&tb)
	SliceContainingInOrder(s, 3, 1).Expect(&tb)
	SliceContainingInOrder(s, 1, 5).Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"[]int does not contain 1 in order",
			"[]int does not contain 5 in order",
		},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}
