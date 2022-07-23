package expect

import "fmt"

type contextMock struct {
	failures []string
}

var _ Context = &contextMock{}

func (t *contextMock) Failf(msg string, args ...any) {
	t.failures = append(t.failures, fmt.Sprintf(msg, args...))
}

func (t *contextMock) Fail(args ...any) {
	t.failures = append(t.failures, fmt.Sprint(args...))
}
