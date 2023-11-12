package expect

import (
	"reflect"
	"testing"

	"github.com/halimath/expect/internal/testhelper"
)

func TestExpectFunc(t *testing.T) {
	f := ExpectFunc(func(t TB) {
		t.Log("got")
	})

	tm := &testhelper.TB{}
	f.Expect(tm)

	if len(tm.Logs) != 1 {
		t.Fatalf("expected 1 log message but got %v", tm.Logs)
	}

	if tm.Logs[0] != "got" {
		t.Errorf("expected 1st log message to be 'got' but got %v", tm.Logs[0])
	}
}

func TestThat_noFailure(t *testing.T) {
	var e1Called, e2Called int
	e1 := ExpectFunc(func(TB) { e1Called++ })
	e2 := ExpectFunc(func(TB) { e2Called++ })

	That(t, e1, e2)
}

func TestPrefixTB(t *testing.T) {
	var tb testhelper.TB

	prefixed := prefixedTB{TB: &tb, prefix: "prefix: "}

	prefixed.Log("Log")
	prefixed.Logf("Logf")
	prefixed.Error("Error")
	prefixed.Errorf("Errorf")
	prefixed.Fatal("Fatal")
	prefixed.Fatalf("Fatalf")
	prefixed.Skip("Skip")
	prefixed.Skipf("Skipf")

	want := testhelper.TB{
		ErrFlag:     true,
		FatalFlag:   true,
		SkippedFlag: true,
		Logs: []string{
			"prefix: Log",
			"prefix: Logf",
			"prefix: Error",
			"prefix: Errorf",
			"prefix: Fatal",
			"prefix: Fatalf",
			"prefix: Skip",
			"prefix: Skipf",
		},
	}

	if diff := reflect.DeepEqual(tb, want); !diff {
		t.Errorf("TB interaction not equal. Wanted %v but got %v", want, tb)
	}
}

func TestWithMessage(t *testing.T) {
	var tb testhelper.TB

	expect := WithMessage(&tb, "pref%s", "ix")
	expect.That(Fail)

	want := testhelper.TB{
		ErrFlag: true,
		Logs: []string{
			"prefix: test failed",
		},
	}

	if diff := reflect.DeepEqual(tb, want); !diff {
		t.Errorf("TB interaction not equal. Wanted %v but got %v", want, tb)
	}
}

func TestFailNow(t *testing.T) {
	var tb testhelper.TB

	That(&tb, FailNow(Fail))

	want := testhelper.TB{
		FatalFlag: true,
		Logs: []string{
			"test failed",
		},
	}

	if diff := reflect.DeepEqual(tb, want); !diff {
		t.Errorf("TB interaction not equal. Wanted %v but got %v", want, tb)
	}
}
