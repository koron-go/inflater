# koron-go/inflater

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/inflater)](https://pkg.go.dev/github.com/koron-go/inflater)
[![Actions/Go](https://github.com/koron-go/inflater/workflows/Go/badge.svg)](https://github.com/koron-go/inflater/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/inflater)](https://goreportcard.com/report/github.com/koron-go/inflater)

## Operations

* V -> S (inflate)
* V -> V
* S -> S (V -> S)
* S -> S (V -> V: map)
* S -> S (V -> discard: filter)
* S + S -> S (join)
* S * S -> S
