package expect

import "fmt"

type tbMock struct {
	errors []string
}

func (*tbMock) Cleanup(func())                    {}
func (*tbMock) Fail()                             {}
func (*tbMock) FailNow()                          {}
func (*tbMock) Failed() bool                      { return false }
func (*tbMock) Fatal(args ...any)                 {}
func (*tbMock) Fatalf(format string, args ...any) {}
func (*tbMock) Helper()                           {}
func (*tbMock) Log(args ...any)                   {}
func (*tbMock) Logf(format string, args ...any)   {}
func (*tbMock) Name() string                      { return "mock" }
func (*tbMock) Setenv(key, value string)          {}
func (*tbMock) Skip(args ...any)                  {}
func (*tbMock) SkipNow()                          {}
func (*tbMock) Skipf(format string, args ...any)  {}
func (*tbMock) Skipped() bool                     { return false }
func (*tbMock) TempDir() string                   { return "tmp" }

func (t *tbMock) Errorf(msg string, args ...any) {
	t.errors = append(t.errors, fmt.Sprintf(msg, args...))
}

func (t *tbMock) Error(args ...any) {
	t.errors = append(t.errors, fmt.Sprint(args...))
}
