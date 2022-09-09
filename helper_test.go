package expect

import "fmt"

type tbMock struct{}

func (tbMock) Log(args ...any) {}
func (tbMock) Fail()           {}
func (tbMock) FailNow()        {}
func (tbMock) Helper()         {}

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
