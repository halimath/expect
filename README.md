# assertthat-go

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]


A library for writing test assertions using golang 1.18 generics to provide a fluent, readable and type-safe 
assertions.

# Installation

This module uses golang modules and can be installed with

```shell
go get github.com/halimath/assertthat-go
```

# Usage in tests

This modules provides three packages, two of which are used when writing tests: `assert` and `is`. Both are
provided to form fluent and readable assertions. Here is a basic example:

```go
assert.That(t, got, is.DeepEqual(MyStruc{
    Foo: "bar",
    Spam: "eggs",
}))
```

The only exported symbol from the `assert` package is the `That` function, which acts as a starting point.
It receives your testing interface (it uses `testing.TB` and can be used for regular and benchmark tests),
the given value and a list of matchers. Each matcher is applied in turn. Matchers receive the testing 
interface and can report `Error`s and `Fatal`s.

The `is` package exports a set of predefined matchers that serve a lot of situations. Using two packages to
export matchers and the `That` function has no other reason but to create a readable API. See below for
a list of defined matchers.

## Predefined matchers

The following table shows the predefined matchers.

Matcher | Type constraints | Description
-- | -- | --
`is.Equal` | `comparable` | Compares given and wanted for equality using the go `==` operator.
`is.DeepEqual` | `any` | Compares given and wanted for deep equality using [github.com/go-test/deep](https://github.com/go-test/deep)

## Defining you own matcher

Defining you own matcher is very simple: Implement a type that implements `matcher.Matcher` which contains a
single function: `Match`. The function receives the testing interface and the given value. 

The following example shows how to implement a matcher for asserting that a given number is even. The example
works only for `int`s. Below we show a very that is a little bit more complicated but uses generics to handle
all kinds of integral numbers.

```go

var IsEven matcher.Matcher[int] = matcher.Func[int](func(t testing.TB, got int) {
    if got%2 != 0 {
        t.Errorf("expected %v to be even", got)
    }
})

func TestCustomMatcher(t *testing.T) {
	assert.That(t, 22, IsEven)
}
```

The example uses the `matcher.Func` convenience type to wrap a plain function as a `matcher.Matcher`. As 
stated above this matcher only supports `int`s. To support other integral number types as well, a little bit
of extra effort is needed.

```go
type Mod interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func IsEven[M Mod]() matcher.Matcher[M] {
	return matcher.Func[M](func(t testing.TB, got M) {
		if got%2 != 0 {
			t.Errorf("expected %v to be even", got)
		}
	})
}

func TestCustomMatcher(t *testing.T) {
	assert.That(t, 22, IsEven[int]())
}
```

This example creates a type constraint interface assembling a union of all the number types a modulo operation
is useful for. It then defines a generic factory function to create an `IsEven` matcher for a given integral
type. This matcher can then be used later in the example. Note that we need to specify the generic type
argument when using the matcher. This is due to the fact, that `IsEven` is not accepting any kind of argument. 
Hopefully, a later version of the go compiler will be able to interfer this argument based on the context it
is used in. If we rewrite this matcher to be a little bit more versatile, we get the following:

```go
func IsDivisableBy[M Mod](d M) matcher.Matcher[M] {
	return matcher.Func[M](func(t testing.TB, got M) {
		if got%d != 0 {
			t.Errorf("expected %d to be divisable by %d", got, d)
		}
	})
}

func TestCustomMatcher2(t *testing.T) {
	var i int = 22
	assert.That(t, i, IsDivisableBy(2))
}
```

As you can see here, there is no need to specify any generic arguments.

# License

Copyright 2022 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[ci-img-url]: https://github.com/halimath/assertthat-go/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/assertthat-go
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/assertthat-go
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/assertthat-go
[release-img-url]: https://img.shields.io/github/v/release/halimath/assertthat-go.svg
[release-url]: https://github.com/halimath/assertthat-go/releases