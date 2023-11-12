package testhelper

import "fmt"

type TB struct {
	ErrFlag     bool
	FatalFlag   bool
	SkippedFlag bool
	Logs        []string
}

func (*TB) Cleanup(func())           {}
func (*TB) Fail()                    {}
func (*TB) FailNow()                 {}
func (t *TB) Failed() bool           { return t.ErrFlag || t.FatalFlag }
func (*TB) Helper()                  {}
func (*TB) Name() string             { return "mock" }
func (*TB) Setenv(key, value string) {}
func (*TB) SkipNow()                 {}
func (t *TB) Skipped() bool          { return t.SkippedFlag }
func (*TB) TempDir() string          { return "tmp" }

func (t *TB) Log(args ...any) {
	t.Logs = append(t.Logs, fmt.Sprint(args...))
}

func (t *TB) Logf(format string, args ...any) {
	t.Log(fmt.Sprintf(format, args...))
}

func (t *TB) Errorf(msg string, args ...any) {
	t.Logf(msg, args...)
	t.ErrFlag = true
}

func (t *TB) Error(args ...any) {
	t.Log(args...)
	t.ErrFlag = true
}

func (t *TB) Fatal(args ...any) {
	t.Log(args...)
	t.FatalFlag = true
}

func (t *TB) Fatalf(format string, args ...any) {
	t.Logf(format, args...)
	t.FatalFlag = true
}

func (t *TB) Skip(args ...any) {
	t.Log(args...)
	t.SkippedFlag = true
}

func (t *TB) Skipf(format string, args ...any) {
	t.Logf(format, args...)
	t.SkippedFlag = true
}
