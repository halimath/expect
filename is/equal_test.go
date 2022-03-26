package is

import "testing"

func TestEqual(t *testing.T) {
	Equal("foo").Match(t, "foo")
}

func TestDeepEqual(t *testing.T) {
	DeepEqual("foo").Match(t, "foo")
}
