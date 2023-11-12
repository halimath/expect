# expect

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]

A library for writing test expectations using golang 1.18 generics to provide a fluent, readable and type-safe 
expectations.

# Installation

This module uses golang modules and can be installed with

```shell
go get github.com/halimath/expect@main
```

## A note on upgrades

This module previously used the module path `get github.com/halimath/expect-go`, preferred dot imports and
used a different API. If you are upgrading please keep in mind to use the new import path in addition to all
the api changes.

# Usage in tests

`expect` provides two packages

```go
import "github.com/halimath/expect"
```

imports the core framework and

```go
import "github.com/halimath/expect/is"
```

imports the bundled expectations.

The following example demonstates the basic use:

```go
expect.That(t, 
	is.DeepEqualTo(got, MyStruc{
    	Foo: "bar",
    	Spam: "eggs",
	}),
)
```

Throughout the examples as well as the source code, we use the term _got_ to represent the value _gotten_ from
some operation under test (the _actual_ value). We use the term _want_ to describe the _wanted_ (or 
_expected_) value to compare to.

Expectations are written using `expect.That` passing in either a `testing.T`, `testing.B` or `testing.F`
followed by a variable number of `expect.Expectation` values (passing in zero results in a no-op). Here is a
somewhat more complex expectation chain:

```go
got := []int{2, 3, 5, 7, 11, 13, 17, 19}

expect.That(t,
	is.SliceOfLen(got, 8),
	is.SliceContainingInOrder(got, 5, 7, 13),
)
```

`expect.That` will execute all expectations in order and report all errors - it behaves just like regular
calls to `t.Error`. All failed expecations will be reported at the very end of the test.

If you want to fail a test immediately - i.e. calling `t.FailNow()` or `t.Fatal(...)`, use the `FailNow` 
decorator to wrap the expectations. You can combine them with regular onces.

```go
got, err := doSomething()

expect.That(t
	expect.FailNow(is.NoError(err)),
	is.DeepEqualTo(got, MyStruc{
    	Foo: "bar",
    	Spam: "eggs",
	}),
)
```

## Standard expectations

The following table shows the predefined expectations provided by `expect`.

Expectation | Type constraints | Description
-- | -- | --
`is.EqualTo` | `comparable` | Compares given and wanted for equality using the go `==` operator.
`is.DeepEqualTo` | `any` | Compares given and wanted for deep equality using reflection.
`is.NoError` | `error` | Expects the given error value to be `nil`.
`is.Error` | `error` | Expects that the given error to be a non-`nil` error that is of the given target error by using `errors.Is` 
`is.MapOfLen` | `map` | Expects the given value to be a map containing the given number of entries
`is.MapContaining` | `map` | Expects the given value to be a map containing a given key, value pair
`is.SliceOfLen` | `slice` | Expects the given value to be a slice containing the given number of values
`is.SliceContaining` | `slice` | Expects the given value to be a slice containing a given set of values in any order
`is.SliceContainingInOrder` | `slice` | Expects the given value to be a slice containing a given list of values in given order
`is.StringOfLen` | `string` | Expects the given value to be a string containing the given number of bytes (not neccessarily runes)
`is.StringContaining` | `string` | Expects the given value to be a string containing a given substring
`is.StringHavingPrefix` | `string` | Expects the given value to be a string having a given prefix
`is.StringHavingSuffix` | `string` | Expects the given value to be a string having a given suffix

### Deep equality

The `is.DeepEqualTo` expectation is special as compared to the other ones. It uses a recursive algorithm to 
compare the given values deeply traversing nested structures using reflection. It handles all primitive types,
interfaces, maps, slices, arrays and structs. It reports all differences found so test failures are easy to
track down.

The equality checking algorithm can be customized on a per-expectation-invocation level using any of the
following options. All options must be given to the `is.DeepEqualTo` call:

```go
expect.That(t,
	is.DeepEqualTo(map[string]int{}, map[string]int(nil), NilMapsAreEmpty(false)),
)
```

#### Floatint point precision

Passing the `FloatPrecision` option allows you to customize the floating point precision when comparing both
`float32` and `float64`. The default value is 10 decimal digits.

#### Nil slices and maps

By default `nil` slices are considered equal to empty ones as well as `nil` maps are considered equal to empty
ones. You can customize this by passing `NilSlicesAreEmpty(false)` or `NilMapsAreEmpty(false)`.

#### Struct fields

Struct fields can be excluded from the comparison using any of the following methods.

Passing `ExcludeUnexportedStructFields(true)` excludes unexported struct fields (those with a name starting
with a lower case letter) from the comparison. The default is not to exclude them.

Using `ExludeTypes` you can exclude all fields with a type given in the list. `ExcludeTypes` is a slice of
`reflect.Type` so you can pass in any number of types.

`ExcludeFields` allows you to specify path expressions (given as strings) that match a path to a field. The
syntax resembles the format used to report differences (so you can simply copy them from the initial test
failure). In addition, you can use a wildcard `*` to match any field or index value.

The following code sample demonstrates the usage:

```go
type nested struct {
	nestedField string
}

type root struct {
	stringField string
	sliceField  []nested
	mapField    map[string]string
}

first := root{
	stringField: "a",
	sliceField: []nested{
		{nestedField: "b"},
	},
	mapField: map[string]string{
		"foo":  "bar",
		"spam": "eggs",
	},
}

second := root{
	stringField: "a",
	sliceField: []nested{
		{nestedField: "c"},
	},
	mapField: map[string]string{
		"foo":  "bar",
		"spam": "spam and eggs",
	},
}

is.DeepEqualTo(first, second, ExcludeFields{
	".sliceField[*].nestedField",
	".mapField[spam]",
})
```

## Defining you own expectation

Defining you own expectation is very simple: Implement a type that implements the `expect.Expecation` 
interface which contains a single method: `Expect`. The method receives a `expect.TB` value which is a 
striped-down version of `testing.TB` (`testing.TB` contains an unexported method and thus cannot be mocked
in external tests).

Perform the matching steps and invoke any method on `expect.TB` to log a message and/or fail the test.
Matchers are encouraged to use `Error`, `Errorf` or `Fail` and leave `Fatal`, `Fatalf` an `FailNow` to the
`expect.FailNow` decorator. This ensures your're expecations are as flexible as the standard ones.
Nevertheless, if your expectation are always meant to fail the test now, its totatly safe to invoke `FailNow`
and get exactly that behavior.

As most expecations can be implemented by a closure function, `expect` provides the `expect.ExpectFunc`
convenience type. Almost all built-in matchers are implemented using `ExpectFunc`.

The following example shows how to implement an expectation for asserting that a given number is even. The
example uses generics to handle all kinds of integral numbers and uses a constraint interface from the
[golang.org/x/exp/constraints](https://pkg.go.dev/golang.org/x/exp/constraints) module.

```go
func IsEven[T constraints.Integer](got T) expect.Expectation {
	return expect.ExpectFunc(func(t expect.TB) {
		if got%2 != 0 {
			t.Errorf("expected <%v> to be even", got)
		}
	})
}

func TestSomething(t *testing.T) {
	var i int = 22
	expect.That(t, IsEven(i))
}
```

# License

Copyright 2022, 2023 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[ci-img-url]: https://github.com/halimath/expect/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/expect
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/expect
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/expect
[release-img-url]: https://img.shields.io/github/v/release/halimath/expect-go.svg
[release-url]: https://github.com/halimath/expect/releases