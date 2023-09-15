package expect

import (
	"reflect"
	"testing"
)

func TestSliceContaining(t *testing.T) {
	var tm tbMock

	s := []int{1, 2, 3, 4}
	IsSliceContaining(1, 3).Match(&tm, s)
	IsSliceContaining(3, 1).Match(&tm, s)
	IsSliceContaining(1, 5).Match(&tm, s)

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{"[]int does not contain [5]"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestSliceContainingInOrder(t *testing.T) {
	var tm tbMock

	s := []int{1, 2, 3, 4}
	IsSliceContainingInOrder(1, 3).Match(&tm, s)
	IsSliceContainingInOrder(3, 1).Match(&tm, s)
	IsSliceContainingInOrder(1, 5).Match(&tm, s)

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{
			"[]int does not contain 1 in order",
			"[]int does not contain 5 in order",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
