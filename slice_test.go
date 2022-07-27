package expect

import (
	"reflect"
	"testing"
)

func TestSliceContaining(t *testing.T) {
	var tm contextMock

	s := []int{1, 2, 3, 4}
	SliceContaining(1, 3).Match(&tm, s)
	SliceContaining(3, 1).Match(&tm, s)
	SliceContaining(1, 5).Match(&tm, s)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"[]int does not contain [5]"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestSliceContainingInOrder(t *testing.T) {
	var tm contextMock

	s := []int{1, 2, 3, 4}
	SliceContainingInOrder(1, 3).Match(&tm, s)
	SliceContainingInOrder(3, 1).Match(&tm, s)
	SliceContainingInOrder(1, 5).Match(&tm, s)

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{
			"[]int does not contain 1 in order",
			"[]int does not contain 5 in order",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
