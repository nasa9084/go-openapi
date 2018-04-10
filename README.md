OpenAPI Specification object model

## SYNOPSIS

``` go
package main

import (
    "fmt"

    "github.com/nasa9084/go-openapi"
)

func main() {
    doc, _ := openapi.Load("path/to/spec")
    fmt.Print(doc.Version)
}
```
