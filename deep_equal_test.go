package expect

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestDeepEqual(t *testing.T) {
	var tm contextMock

	DeepEqual("foo").Match(&tm, "foo")
	DeepEqual("foo").Match(&tm, "bar")

	want := contextMock{
		failures: []string{"values are not deeply equal:\n  want: foo\n   got: bar"},
	}

	if !reflect.DeepEqual(tm, want) {
		t.Errorf("unexpected failures:\n%s\n%s", want, tm)
	}
}

func TestDeepEquals_standardCases(t *testing.T) {
	type test struct {
		want, got any
		diff      diff
	}

	someString := "some value"
	someSlice := []int{0, 1}
	someMap := make(map[string]int)

	type someStruct struct {
		A string
		b int
	}

	tests := []test{
		// Type checking
		{"", 0, diff{{"", "string", "int"}}},

		// nil handling
		{nil, nil, nil},
		{"", nil, diff{{"", "", "<nil>"}}},
		{nil, "", diff{{"", "<nil>", ""}}},

		// strings
		{"foo", "foo", nil},
		{"foo", "bar", diff{{"", "foo", "bar"}}},
		{&someString, &someString, nil},

		// ints
		{0, 0, nil},
		{uint(0), uint(0), nil},

		// floats
		{0.0, 0.0, nil},
		{0.0, 1.0, diff{{"", "0.0000000000", "1.0000000000"}}},

		// bools
		{true, true, nil},
		{false, false, nil},
		{true, false, diff{{"", "true", "false"}}},

		// slices
		{[]int{0, 1}, []int{0, 1}, nil},
		{[]int{0, 1}, someSlice, nil},
		{[]int{}, []int(nil), nil},
		{[]int(nil), []int{}, nil},
		{someSlice, someSlice, nil},
		{[]int{0, 1}, []int{0}, diff{{"[1]", "1", "<missing slice index>"}}},
		{[]int{0, 1}, []int{2, 3}, diff{
			{"[0]", "0", "2"},
			{"[1]", "1", "3"},
		}},
		{[]int{}, []int{1, 2}, diff{
			{"[0]", "<unwanted slice index>", "1"},
			{"[1]", "<unwanted slice index>", "2"},
		}},

		// arrays
		{[2]int{0, 1}, [2]int{0, 1}, nil},
		{[2]int{0, 1}, [1]int{0}, diff{{"", "[2]int", "[1]int"}}},

		// maps
		{map[string]string{}, map[string]string{}, nil},
		{map[string]string{}, map[string]string(nil), nil},
		{map[string]string(nil), map[string]string(nil), nil},
		{map[string]string{"foo": "foo"}, map[string]string{"foo": "foo"}, nil},
		{someMap, someMap, nil},
		{&someMap, &someMap, nil},
		{map[string]string{"foo": "foo"}, map[string]string{"foo": "bar"},
			diff{{
				"[foo]", "foo", "bar",
			}}},
		{map[string]string{"foo": "foo"}, map[string]string{},
			diff{{
				"[foo]", "foo", "<missing map key>",
			}}},
		{map[string]string{"foo": "foo"}, map[string]string{"bar": "bar"},
			diff{
				{"[foo]", "foo", "<missing map key>"},
				{"[bar]", "<missing map key>", "bar"},
			}},

		// structs
		{someStruct{}, someStruct{}, nil},
		{someStruct{"a", 1}, someStruct{"b", 2}, diff{
			{".A", "a", "b"},
			{".b", "1", "2"},
		}},

		// special values
		{fmt.Errorf("foo"), errors.New("foo"), nil},
	}

	for _, test := range tests {
		got := deepEquals(test.want, test.got)

		if !reflect.DeepEqual(test.diff, got) {
			t.Errorf("%#v == %#v\nwant: %#v\n got: %#v", test.want, test.got, test.diff, got)
		}
	}
}

func TestDeepEquals_nilSlicesAreNotEmpty(t *testing.T) {
	got := deepEquals([]int{}, []int(nil), NilSlicesAreEmpty(false))
	want := diff{{"", "[]", "<nil slice>"}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %#v but got %#v", want, got)
	}
}

func TestDeepEquals_nilMapsAreNotEmpty(t *testing.T) {
	got := deepEquals(map[string]int{}, map[string]int(nil), NilMapsAreEmpty(false))
	want := diff{{"", "map[]", "<nil map>"}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %#v but got %#v", want, got)
	}
}

func TestDeepEquals_excludeUnexportedFields(t *testing.T) {
	type s struct {
		f string
	}

	got := deepEquals(s{"a"}, s{"b"}, ExcludeUnexportedStructFields(true))

	if got != nil {
		t.Errorf("expected no diff but got %#v", got)
	}
}

func TestDeepEquals_excludeFieldsOfType(t *testing.T) {
	type s struct {
		f string
		t int
	}

	got := deepEquals(s{"a", 0}, s{"a", 1}, ExcludeTypes{reflect.TypeOf(int(0))})

	if got != nil {
		t.Errorf("expected no diff but got %#v", got)
	}
}
