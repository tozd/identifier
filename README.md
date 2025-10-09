# Readable global identifiers

[![pkg.go.dev](https://pkg.go.dev/badge/gitlab.com/tozd/identifier)](https://pkg.go.dev/gitlab.com/tozd/identifier)
[![NPM](https://img.shields.io/npm/v/@tozd/identifier.svg)](https://www.npmjs.com/package/@tozd/identifier)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/tozd/identifier)](https://goreportcard.com/report/gitlab.com/tozd/identifier)
[![pipeline status](https://gitlab.com/tozd/identifier/badges/main/pipeline.svg?ignore_skipped=true)](https://gitlab.com/tozd/identifier/-/pipelines)
[![coverage report](https://gitlab.com/tozd/identifier/badges/main/coverage.svg)](https://gitlab.com/tozd/identifier/-/graphs/main/charts)

A package providing functions to generate and parse readable global identifiers.

Features:

- Identifiers have 128 bits of entropy, making them suitable as global identifiers.
- By default identifiers are random, but you can convert existing
  [UUIDs](https://en.wikipedia.org/wiki/Universally_unique_identifier).
- They are encoded into readable base 58 strings always of 22 characters in length.

## Installation

### Go installation

You can add it to your project using `go get`:

```sh
go get gitlab.com/tozd/identifier
```

It requires Go 1.24 or newer.

### TypeScript/JavaScript installation

You can add it to your project using `npm`:

```sh
npm install --save @tozd/identifer
```

It requires node 22 or newer. It works in browsers, too.

## Usage

### Go usage

See full package documentation on [pkg.go.dev](https://pkg.go.dev/gitlab.com/tozd/identifier#section-documentation).

### TypeScript/JavaScript usage

```js
import { Identifier } from "@tozd/identifier"

const id = Identifier.new() // A random identifier.
const s = id.toString()
console.log(s)
Identifier.valid(s) // True.
Identifier.fromString(s) // Is equal to id.
const u = Identifier.fromUUID("c97e2491-dd58-4a4e-b351-d786554e2ae6") // Is equal to Rt7JRSoDY1woPhLidZNvz1.
JSON.stringify({ id }) // Works, id is converted to string.
```

## Related projects

- [Nano ID](https://github.com/ai/nanoid) – a similar project which allows more customization (both choosing
  the alphabet and the size); this project supports only one type of identifiers to make sure everyone using
  it has the same identifiers.

## GitHub mirror

There is also a [read-only GitHub mirror available](https://github.com/tozd/identifier),
if you need to fork the project there.
