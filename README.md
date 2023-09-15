# expect-go

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]


A library for writing test expectations using golang 1.18 generics to provide a fluent, readable and type-safe 
expectations.

# Installation

This module uses golang modules and can be installed with

```shell
go get github.com/halimath/expect-go@main
```

# Usage in tests

`expect-go` is intended to be "dot imported" into your tests. This makes all the exported functions and 
matchers directly available leading to very fluent and readable expectations:

```go
import . "github.com/halimath/expect-go"
```

All following example assume the dot import.

The following example demonstates the basic use:

```go
ExpectThat(t, got,
	IsDeepEqualTo(MyStruc{
    	Foo: "bar",
    	Spam: "eggs",
	}))
```

Expectations are written using `ExpectThat` passing in either a `testing.T` or a `testing.B` followed by the
actual value and zero or more matchers (passing in zero results in a no-op). Here is a somewhat more complex
expectation:

```go
actual := []int{2, 3, 5, 7, 11, 13, 17, 19}

ExpectThat(t, actual,
	HasLenOf[[]int](8),
	IsSliceContainingInOrder(5, 7, 13),
)
```

`ExpectThat` will execute all matchers in order and report all errors. Any matcher that produces a negative
result causes the test to fail at the end. 

If you want to fail a test immediately, use `EnsureThat` which has the same signature:

```go
EnsureThat(t, got,
	IsDeepEqualTo(MyStruc{
    	Foo: "bar",
    	Spam: "eggs",
	}))
```

## Standard matchers

The following table shows the predefined matchers provided by `expect-go`.

Matcher | Type constraints | Description
-- | -- | --
`IsNil` | `any` | Expects a pointer to be `nil`
`IsNotNil` | `any` | Expects a pointer to be non `nil`
`IsEqualTo` | `comparable` | Compares given and wanted for equality using the go `==` operator.
`IsDeepEqualTo` | `any` | Compares given and wanted for deep equality using reflection.
`IsNoError` | `error` | Expects the given error value to be `nil`.
`IsError` | `error` | Expects that the given error to be a non-`nil` error that is of the given target error by using `errors.Is` 
`HasLenOf` | `string`, `array`, `slice`, `map`, `channel` | Expects the length of the given value to equal a given length
`IsMapContaining` | `map` | Expects the given value to be a map containing a given key, value pair
`IsSliceContaining` | `slice` | Expects the given value to be a slice containing a given set of values in any order
`IsSliceContainingInOrder` | `slice` | Expects the given value to be a slice containing a given list of values in given order
`IsStringContaining` | `string` | Expects the given value to be a string containing a given substring
`IsStringHavingPrefix` | `string` | Expects the given value to be a string having a given prefix
`IsStringHavingSuffix` | `string` | Expects the given value to be a string having a given suffix

### Deep equality

The `IsDeepEqualTo` matcher is special as compared to the other ones. It uses a recursive algorithm to compare
the given values deeply traversing nested structures. It handles all primitive types, interfaces, maps, slices,
arrays and structs. It reports all differences found so test failures are easy to track down.

The equality checking algorithm can be customized on a per-matcher-invocation level using any of the following
options. All options must be given to the `IsDeepEqualTo` matcher:

```go
ExpectThat(t, map[string]int{},
	IsDeepEqualTo(map[string]int(nil), NilMapsAreEmpty(false)))
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

got := deepEquals(first, second, ExcludeFields{
	".sliceField[*].nestedField",
	".mapField[spam]",
})
```



## Defining you own matcher

Defining you own matcher is very simple: Implement a type that implements the `Matcher` interface which
contains a single method: `Match`. The method receives a `TB` and the actual value (for internal testability)
this module defines a `TB` interface which is a striped-down version of `testing.TB`.

Perform the matching steps and invoke any method on `TB` to log a message and/or fail the test. Matchers are
encouraged to use `Error` and `Errorf` as this will work with both `ExpectThat` and `EnsureThat` (internally,
`EnsureThat` just uses a delegate wrapper for `TB`, which delegate `Error` to `Fail` ...).

As most matchers can be implemented by a closure function, `expect-go` provides the `MatcherFunc` convenience
type. Almost all built-in matchers are implemented using `MatcherFunc`.

The following example shows how to implement a matcher for asserting that a given number is even. The example
uses generics to handle all kinds of integral numbers.

```go
type Mod interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func IsEven[G Mod]() Matcher[G] {
	return MatcherFunc[G](func(t TB, got G) {
		if got%2 != 0 {
			t.Errorf("expected <%v> to be even", got)
		}
	})
}

func TestCustomMatcher(t *testing.T) {
	var i int = 22
	ExpectThat(t, i, IsEven[int]())
}
```

This example creates a type constraint interface assembling a union of all the number types a modulo operation
is useful for. It then defines a generic factory function `IsEven` to create a custom matcher for a given
integral type implemented as a closure using the `MatcherFunc` type. Note how the naming `IsEven` supports the
fluent style of the expectations.

Note that we need to specify the generic type argument when using the matcher. This is due to the fact, that 
`IsEven` is not accepting any kind of argument. Hopefully, a later version of the go compiler will be able to 
interfer the type argument based on the context it is used in. 

We can rewrite this matcher to be a little bit more versatile, we get the following:

```go
func IsDivisableBy[G Mod](d G) Matcher[G] {
	return MatcherFunc[G](func(t TB, got G) {
		if got%d != 0 {
			t.Errorf("expected <%v> to be divisable by <%v>", got, d)
		}
	})
}

func TestCustomMatcher2(t *testing.T) {
	var i int = 22
	ExpectThat(t, i, IsDivisableBy(2))
}
```

As you can see here, there is no need to specify any generic arguments.

# License

Copyright 2022, 2023 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[ci-img-url]: https://github.com/halimath/expect-go/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/expect-go
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/expect-go
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/expect-go
[release-img-url]: https://img.shields.io/github/v/release/halimath/expect-go.svg
[release-url]: https://github.com/halimath/expect-go/releases