# koron-go/inflater

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/inflater)](https://pkg.go.dev/github.com/koron-go/inflater)
[![Actions/Go](https://github.com/koron-go/inflater/workflows/Go/badge.svg)](https://github.com/koron-go/inflater/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/inflater)](https://goreportcard.com/report/github.com/koron-go/inflater)

The package `inflater` provides an `Inflater` interface that inflates a value
into multiple values and some operations to combine them `Inflater`s.

## Install and Upgrade

```console
$ go get github.com/koron-go/inflater@latest
```

## Description

The first idea of the `inflater` was very simple: to perform an algorithm that
predicts and expands multiple strings from a single input string using a single
interface combination.  Therefore, the first Inflate function was defined as
`func Inflate(s string) <-chan string`

Combining this with the [rangefunc][rangefunc] added in Go 1.23 and using
generics instead of `string`, we have the following `Inflater` interface.  This
is intended to inflate multiple values based on a `seed`.

```go
type Inflater[V any] interface {
	Inflate(seed V) iter.Seq[V]
}
```

The package `inflater` provides several special `Inflater` types.

* `None` - An Inflater that consumes all input and produces no output, like a black hole.
* `Keep` - An Inflater that outputs the input as is.
* `Map` - An Inflater that converts input with a function and outputs it.
* `Filter` - An Inflater that judges and filters input with a function.

It also provides the ability to combine multiple `Inflater` objects to create another `Inflater` object.

* `Parallel` - Pass the input to multiple `Inflater`s and concatenate their outputs into a single `iter.Seq`.
* `Serial` - Inflate the input with the first `Inflater`, inflate the output of that with the next `Inflater`, and repeat that to get the output iter.Seq

In addition, we provide two `Inflater`s that perform basic operations on strings.

* `Prefix` - An Inflater that outputs multiple values with an each prefix string attached to the input string.
* `Suffix` - An Inflater that outputs multiple values with an each suffix string attached to the input string.

[rangefunc]:https://tip.golang.org/wiki/RangefuncExperiment
