package is

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/halimath/expect/internal/testhelper"
)

func TestNoError(t *testing.T) {
	var tm testhelper.TB

	NoError(nil).Expect(&tm)
	NoError(errors.New("failed")).Expect(&tm)

	if !reflect.DeepEqual(tm, testhelper.TB{
		ErrFlag: true,
		Logs:    []string{"expected no error but got \"failed\""},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}

func TestError(t *testing.T) {
	var tm testhelper.TB

	err := errors.New("failed")

	Error(nil, nil).Expect(&tm)
	Error(nil, err).Expect(&tm)
	Error(errors.New("other"), err).Expect(&tm)
	Error(err, err).Expect(&tm)
	Error(fmt.Errorf("wrapped %w", err), err).Expect(&tm)

	if !reflect.DeepEqual(tm, testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"expected an error with target failed but got nil",
			"expected an error with target failed but got other",
		},
	}) {
		t.Errorf("not expected: %#v", tm)
	}
}
