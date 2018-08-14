OpenAPI Specification object model
===

[![GoDoc](https://godoc.org/github.com/nasa9084/go-openapi?status.svg)](https://godoc.org/github.com/nasa9084/go-openapi)
[![Build Status](https://travis-ci.org/nasa9084/go-openapi.svg?branch=master)](https://travis-ci.org/nasa9084/go-openapi)
[![codecov](https://codecov.io/gh/nasa9084/go-openapi/branch/master/graph/badge.svg)](https://codecov.io/gh/nasa9084/go-openapi)

---

## Overview

This is an implementation of [OpenAPI Specification 3.0](https://github.com/OAI/OpenAPI-Specification) object model.

## Synopsis

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

## Status

* [x] Model definition
* [x] Load OpenAPI 3.0 spec file
* [ ] Resolve Reference object
  * [x] Resolve #/component reference
  * [ ] Resolve other file reference
* [ ] Validation
  * [x] Validate spec values
    * [ ] test for validation
      * [x] Document
      * [x] Info
      * [x] Contact
      * [x] License
      * [x] Server
      * [x] ServerVariable
      * [x] Paths
      * [x] PathItem
      * [ ] Operation
      * [ ] Parameter
      * [ ] RequestBody
      * [ ] Responses
      * [ ] Response
      * [ ] Callbacks
      * [ ] Callback
      * [ ] Schema
      * [ ] Example
      * [ ] MediaType
      * [ ] Header
      * [ ] Link
      * [ ] Encoding
      * [x] Discriminator
      * [x] XML
      * [ ] Components
      * [x] SecurityScheme
      * [x] OAuthFlows
      * [x] OAuthFlow
      * [ ] SecurityRequirement
      * [x] Tag
      * [x] ExternalDocumentation
  * [ ] Validate HTTP Request
  * [ ] Validate HTTP Response
