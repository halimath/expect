package expect

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNoError(t *testing.T) {
	var tm contextMock

	NoError().Match(&tm, nil)
	NoError().Match(&tm, errors.New("failed"))

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{"expected no error but got failed"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestError(t *testing.T) {
	var tm contextMock

	err := errors.New("failed")

	Error(err).Match(&tm, nil)
	Error(err).Match(&tm, errors.New("other"))
	Error(err).Match(&tm, err)
	Error(err).Match(&tm, fmt.Errorf("wrapped %w", err))

	if !reflect.DeepEqual(tm, contextMock{
		failures: []string{
			"expected an error with target failed but got nil",
			"expected an error with target failed but got other",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
