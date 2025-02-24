# Opt

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/metafates/opt)

Opt is a Go package for safe abstractions over optional values.

Inspired by the [Option type in Rust] and follows the same ideas and function signatures.

## Install

```bash
go get github.com/metafates/opt
```

## Usage

```go
package main

import (
    "fmt"
    "math"
    "github.com/metafates/opt"
)

func divide(numerator, denominator float64) Opt[float64] {
    if denominator == 0 {
        return None[float64]()
    }

    return Some(numerator / denominator)
}

func main() {
    result := divide(2, 3)

    fmt.Println(result.UnwrapOr(math.MaxInt))
}
```

[Option type in Rust]: https://doc.rust-lang.org/std/option/enum.Option.html
