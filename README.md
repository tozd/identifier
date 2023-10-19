# Readable global identifiers

[![pkg.go.dev](https://pkg.go.dev/badge/gitlab.com/tozd/identifier)](https://pkg.go.dev/gitlab.com/tozd/identifier)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/tozd/identifier)](https://goreportcard.com/report/gitlab.com/tozd/identifier)
[![pipeline status](https://gitlab.com/tozd/identifier/badges/main/pipeline.svg?ignore_skipped=true)](https://gitlab.com/tozd/identifier/-/pipelines)
[![coverage report](https://gitlab.com/tozd/identifier/badges/main/coverage.svg)](https://gitlab.com/tozd/identifier/-/graphs/main/charts)

A Go package providing functions to generate and parse readable global identifiers.

Features:

- Identifiers have 128 bits of entropy, making them suitable as global identifiers.
- By default identifiers are random, but you can convert existing
  [UUIDs](https://en.wikipedia.org/wiki/Universally_unique_identifier).
- They are encoded into readable base 58 strings always of 22 characters in length.

## Installation

This is a Go package. You can add it to your project using `go get`:

```sh
go get gitlab.com/tozd/identifier
```

It requires Go 1.17 or newer.

## Usage

See full package documentation on [pkg.go.dev](https://pkg.go.dev/gitlab.com/tozd/identifier#section-documentation).

## GitHub mirror

There is also a [read-only GitHub mirror available](https://github.com/tozd/identifier),
if you need to fork the project there.
