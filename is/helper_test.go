package is

import "fmt"

type tMock struct {
	errors []string
	fatals []string
}

func (t *tMock) Errorf(msg string, args ...any) {
	t.errors = append(t.errors, fmt.Sprintf(msg, args...))
}

func (t *tMock) Fatalf(msg string, args ...any) {
	t.fatals = append(t.fatals, fmt.Sprintf(msg, args...))
}
