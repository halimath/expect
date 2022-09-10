package expect

import "fmt"

type tbMock struct{}

func (tbMock) Cleanup(func())                    {}
func (tbMock) Error(args ...any)                 {}
func (tbMock) Errorf(format string, args ...any) {}
func (tbMock) Fail()                             {}
func (tbMock) FailNow()                          {}
func (tbMock) Failed() bool                      { return false }
func (tbMock) Fatal(args ...any)                 {}
func (tbMock) Fatalf(format string, args ...any) {}
func (tbMock) Helper()                           {}
func (tbMock) Log(args ...any)                   {}
func (tbMock) Logf(format string, args ...any)   {}
func (tbMock) Name() string                      { return "mock" }
func (tbMock) Setenv(key, value string)          {}
func (tbMock) Skip(args ...any)                  {}
func (tbMock) SkipNow()                          {}
func (tbMock) Skipf(format string, args ...any)  {}
func (tbMock) Skipped() bool                     { return false }
func (tbMock) TempDir() string                   { return "tmp" }

type contextMock struct {
	failures []string
	t        tbMock
}

var _ Context = &contextMock{}

func (t *contextMock) Failf(msg string, args ...any) {
	t.failures = append(t.failures, fmt.Sprintf(msg, args...))
}

func (t *contextMock) Fail(args ...any) {
	t.failures = append(t.failures, fmt.Sprint(args...))
}

func (t *contextMock) T() TB { return t.t }
