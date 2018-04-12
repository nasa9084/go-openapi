OpenAPI Specification object model
===

[![GoDoc](https://godoc.org/github.com/nasa9084/go-openapi?status.svg)](https://godoc.org/github.com/nasa9084/go-openapi)

---

## SYNOPSIS

``` go
package main

import (
    "fmt"

    "github.com/nasa9084/go-openapi"
)

func main() {
    doc, _ := openapi.LoadFile("path/to/spec")
    fmt.Print(doc.Version)
}
```
