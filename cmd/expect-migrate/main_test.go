package main

import (
	"os"
	"strings"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

// func TestDebug(t *testing.T) {
// 	in := `package test
// 	func TestMemfs_OpenFile(t *testing.T) {
// 		ExpectThat(t, err).Is(NoError())
// 		ExpectThat(t, err).Is(NoError())
// 		ExpectThat(t, err).Is(NoError())

// 		var s string
// 	}
// 	`

// 	// 	in := `
// 	// package test

// 	// import (
// 	// 	"testing"
// 	// 	. "github.com/halimath/expect-go"
// 	// )

// 	// func TestSomething(t *testing.T) {
// 	// 	var err error
// 	// 	s := produceValue()
// 	// 	EnsureThat(t, err).Is(NoError())

// 	// }
// 	// `

// 	var sb strings.Builder
// 	err := migrate("test.go", strings.NewReader(in), &sb)
// 	expect.That(t, is.NoError(err))
// 	t.Log(sb.String())
// }

func TestMigrate_rewriteOldCode(t *testing.T) {
	type testCase struct {
		label, src, want string
	}

	tests := []testCase{
		{
			label: "should rewrite old test code",
			src: `
package test

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
				`,
			want: strings.TrimSpace(`
package test

import (
	"testing"
	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestSomething(t *testing.T) {
	var err error
	s := produceValue()
	expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(s, "foo"))

}
				`),
		},
		{
			label: "should rewrite old test code with inline got value",
			src: `
package test

import (
	"testing"
	. "github.com/halimath/expect-go"
)

func TestSomething(t *testing.T) {
	EnsureThat(t, l).Is(Equal(len("hello, world")))
	EnsureThat(t, f.Close()).Is(NoError())
}
						`,
			want: strings.TrimSpace(`
package test

import (
	"testing"
	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestSomething(t *testing.T) {
	expect.That(t, expect.FailNow(is.EqualTo(l, len("hello, world"))), expect.FailNow(is.NoError(f.Close())))

}
						`),
		},
		{
			label: "should not rewrite new test code",
			src: `
package test

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
							`,

			want: strings.TrimSpace(`
package test

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
						`),
		},
		// -------------------------
		{
			label: "should rewrite code from Run functions with multiple expectations",
			src: `
package test

import (
	"testing"
	. "github.com/halimath/expect-go"
	. "github.com/halimath/fixture"
)

func TestMemfs_OpenFile(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("success", func(t *testing.T, f *memfsFixture) {
			file, err := f.fs.OpenFile("open_file", fsx.O_RDWR|fsx.O_CREATE, 0644)
			EnsureThat(t, err).Is(NoError())

			l, err := file.Write([]byte("hello, world"))
			EnsureThat(t, err).Is(NoError())
			EnsureThat(t, l).Is(Equal(len("hello, world")))

			EnsureThat(t, file.Close()).Is(NoError())

			got, err := fs.ReadFile(f.fs, "open_file")
			EnsureThat(t, err).Is(NoError())
			ExpectThat(t, string(got)).Is(Equal("hello, world"))
		})
}
					`,
			want: `package test

import (
	"testing"
	"github.com/halimath/expect"
	. "github.com/halimath/fixture"
	"github.com/halimath/expect/is"
)

func TestMemfs_OpenFile(t *testing.T) {
	With(t, new(memfsFixture)).
		Run("success", func(t *testing.T, f *memfsFixture) {
			file, err := f.fs.OpenFile("open_file", fsx.O_RDWR|fsx.O_CREATE, 0644)
			expect.That(t, expect.FailNow(is.NoError(err)))

			l, err := file.Write([]byte("hello, world"))
			expect.That(t, expect.FailNow(is.NoError(err)), expect.FailNow(is.EqualTo(l, len("hello, world"))), expect.FailNow(is.NoError(file.Close())))

			got, err := fs.ReadFile(f.fs, "open_file")
			expect.That(t, expect.FailNow(is.NoError(err)), is.EqualTo(string(got), "hello, world"))

		})
}`,
		},
	}

	for _, test := range tests {
		var sb strings.Builder
		err := migrate("unit_test.go", strings.NewReader(test.src), &sb)
		got := strings.TrimSpace(sb.String())

		expect.WithMessage(t, test.label).That(
			is.NoError(err),
			is.EqualToStringByLines(got, test.want),
		)
	}
}

func TestMigrate_realWorldTestcase(t *testing.T) {
	defer func() {
		os.Rename("./testdata/in_test.go~", "./testdata/in_test.go")
	}()

	err := migrateFile("./testdata/in_test.go")
	expect.That(t, expect.FailNow(is.NoError(err)))

	want, err := os.ReadFile("./testdata/want_test.go")
	expect.That(t, expect.FailNow(is.NoError(err)))

	got, err := os.ReadFile("./testdata/in_test.go")
	expect.That(t, expect.FailNow(is.NoError(err)))

	expect.That(t, is.EqualToStringByLines(string(got), string(want)))
}
