package expect

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNoError(t *testing.T) {
	var tm tbMock

	NoError().Match(&tm, nil)
	NoError().Match(&tm, errors.New("failed"))

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{"expected <{failed}> to be nil"},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestError(t *testing.T) {
	var tm tbMock

	err := errors.New("failed")

	IsError(err).Match(&tm, nil)
	IsError(err).Match(&tm, errors.New("other"))
	IsError(err).Match(&tm, err)
	IsError(err).Match(&tm, fmt.Errorf("wrapped %w", err))

	if !reflect.DeepEqual(tm, tbMock{
		errors: []string{
			"expected an error with target failed but got nil",
			"expected an error with target failed but got other",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
