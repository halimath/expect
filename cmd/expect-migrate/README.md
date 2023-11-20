# expect-migrate

`expect-migrate` is a small command line tool to upgrade your go test source files
using `github.com/halimath/expect-go@v0.1.0` to `github.com/halimath/expect@latest`. The tool rewrites
go source files that in a lot of cases directly compile and run. In some cases minor work by the developer
is needed.

# Installation

```shell
go install github.com/halimath/expect/cmd/expect-migrate@main
```

# Usage

Given an input file `some_test.go`

```go
package some

import (
	"testing"
	. "github.com/halimath/expect-go"
)

func TestSomething(t *testing.T) {
	var err error
	s := produceValue()

	EnsureThat(t, err).Is(NoError())
	ExpectThat(t, s).Is(Equal("foo"))
}
```

applying

```shell
expect-migrate some_test.go
```

will rewrite the file to 

```go
package some

import (
	"testing"
	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestSomething(t *testing.T) {
	var err error
	s := produceValue()
	expect.That(t, expect.FailNow(is.NoError(err)))
	expect.That(t, is.EqualTo(s, "foo"))
}
```
