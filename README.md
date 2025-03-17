# Opt

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/metafates/opt)

Opt is a Go package for safe abstractions over optional values.

Inspired by the [Option type in Rust] and follows the same ideas and function signatures.

## Features

- Represent explicitly set values. For example: `{"b":2,"a":null}` and `{"b":2}` would be different states for `a` - explicit and implicit `None`.
- All the encoding and decoding functionality: json, gob, sql, text & binary.
- Adapters to construct options from pointers, zero values, and proto messages.
- No reflection.

## Install

```bash
go get github.com/metafates/opt
```

## Usage

See [example_test.go](./example_test.go) for more examples.

```go
package main

import (
	"fmt"
	"time"

	"github.com/metafates/opt"
)

type User struct {
	Birth opt.Opt[time.Time]
}

func (u User) Age() opt.Opt[int] {
	return opt.Map(u.Birth, func(t time.Time) int {
		return time.Now().Year() - t.Year()
	})
}

func isAdult(age int) bool {
	return age >= 18
}

func getUser(id int) opt.Opt[User] {
	// ...
}

func main() {
	if opt.AndThen(getUser(42), User.Age).IsSomeAnd(isAdult) {
		fmt.Println("🍺!")
	}
}
```

[Option type in Rust]: https://doc.rust-lang.org/std/option/enum.Option.html
