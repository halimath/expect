package is

import (
	"reflect"
	"testing"

	"github.com/halimath/expect/internal/testhelper"
)

func TestStringOfLen(t *testing.T) {
	var tb testhelper.TB

	StringOfLen("foobar", 2).Expect(&tb)
	StringOfLen("foobar", 6).Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"expected \"foobar\" to have len 2 but got 6",
		},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}

func TestStringContaining(t *testing.T) {
	var tb testhelper.TB

	StringContaining("foobar", "oba").Expect(&tb)
	StringContaining("spameggs", "oba").Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"expected \"spameggs\" to contain \"oba\"",
		},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}

func TestStringWithPrefix(t *testing.T) {
	var tb testhelper.TB

	StringWithPrefix("foobar", "foo").Expect(&tb)
	StringWithPrefix("spameggs", "foo").Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"expected \"spameggs\" to have prefix \"foo\"",
		},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}

func TestStringWithSuffix(t *testing.T) {
	var tb testhelper.TB

	StringWithSuffix("foobar", "bar").Expect(&tb)
	StringWithSuffix("spameggs", "bar").Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"expected \"spameggs\" to have suffix \"bar\"",
		},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}

func TestEqualToStringByLines(t *testing.T) {
	var tb testhelper.TB

	EqualToStringByLines("foo\nbar", "foo\nbar").Expect(&tb)
	EqualToStringByLines("foo\nbar", "foobar").Expect(&tb)
	EqualToStringByLines("foo\nbar", "foo").Expect(&tb)
	EqualToStringByLines("foo", "foo\nbar").Expect(&tb)
	EqualToStringByLines("foo\nbar", "foo\nspam").Expect(&tb)
	EqualToStringByLines("foo\n  bar", "foo\n\tbar", DedentLines).Expect(&tb)

	if !reflect.DeepEqual(tb, testhelper.TB{
		ErrFlag: true,

		Logs: []string{
			"expected string to have 1 lines but got 2",
			"at line 0: wanted\n\"foobar\"\nbut got\n\"foo\"",
			"expected string to have 1 lines but got 2",
			"line 1: wanted no line but got\n\"bar\"",
			"expected string to have 2 lines but got 1",
			"line 1: wanted\n\"bar\"\nbut got no line",
			"at line 1: wanted\n\"spam\"\nbut got\n\"bar\""},
	}) {
		t.Errorf("not expected: %#v", tb)
	}
}
