# User Guide

- [Overview](#overview)
- [Composition](#composition)
- [Life-cycle](#life-cycle)
- [Endpoint types](#endpoint-types)
- [High-Order Endpoints](#high-order-endpoints)
- [Mapping Endpoints](#mapping-endpoints)
- [Outputs](#outputs)
- [Unit Testing](#unit-testing)

## Overview

A `µ.Endpoint` is a key abstraction in the framework. It is a *pure function* that takes HTTP request as Input and return Output (result of request evaluation).

```go
/*

Endpoint: Input ⟼ Output
*/
type Endpoint func(*Input) error
```

`Input` is a convenient wrapper of HTTP request with some Gouldian specific context. This library supports integration with 
* [`APIGatewayProxyRequest`](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go).

`Output` is a sum type that captures a result of `Endpoint` evaluation. Technically, it indicates if:
* it do not match the request
* it matches the request
* it successfully transforms the request to HTTP output
* it fails to transform the request

Golang is missing generics and type variance. Therefore, the Output is always an `error` value. The library supplies [primitives](../output.go) to declare output using HTTP status codes notation (e.g. `Ok`, `Created`, `BadRequest`, `Unauthorized`, etc).

## Composition

`Endpoint A` can be composed with `Endpoint B` into new `Endpoint C`. The library supports two combinators: and-then, or-else.

**and-then**

Use `and-then` to build product Endpoint: `A × B ⟼ C`. The product type matches Input if each composed function successfully matches it. Compose them with `Then` function or variadic alternative `µ.Join`.

```go
var a: µ.Endpoint = /* ... */
var b: µ.Endpoint = /* ... */

//
// You can use `Then` method declared for the Endpoint type
c := a.Then(b)

//
// Alternatively, variadic function `Join` does same for sequence of Endpoints
c := µ.Join(a, b)
```

**or-else**

Use `or-else` to build co-product Endpoint: `A ⨁ B ⟼ C`. The co-product is also known as sum-type that matches the first successful function.

```go
var a: µ.Endpoint = /* ... */
var b: µ.Endpoint = /* ... */

//
// You can use `Or` method declared for the Endpoint type
c := a.Or(b)

//
// Alternatively, variadic `Or` variant does same for sequence of Endpoints
c := µ.Or(a, b)
```

These rules of Endpoint composition allow developers to build any complex HTTP request handling from small re-usable block.


## Life-cycle

The function `gouldian.Serve` builds a co-product Endpoint to define entire *"api algebra"* of application. Each incoming HTTP request passed to this *"algebra"*. It is important to understand the life-cycle behavior for development of a [High-Order Endpoints](#high-order-endpoints) and writing a [Unit Testing](#unit-testing) in your application.

1. The library envelops each incoming request to `Input` type and applies it to the endpoint `api(input)`.
2. The resulting value of `error` (aka `Output`) type is matched against
* `NoMatch` causes abort of current *product* `Endpoint`. The request is passed to succeeding *co-product* `Endpoint`.
* `nil` continues evaluation of *product* `Endpoint` to succeeding item.
* `error` aborts the evaluation of the application. The output error value is output to HTTP client

```
        NoMatch: next co-product          error: return Output
            +------+                            +---------+--->
           /        \                          /         /
[ A × B × C × D ] ⨁ [ A × B × C × D ] ⨁ [ A × B × C × D ]
  \   /                              \
   +-+                                +---> 
 nil: next product                  nil: exits co-product
```

## Endpoint types

Gouldian library delivers set of built-in endpoints to deal with HTTP request processing.

**Match HTTP Verb/Method**

`func Method(verb string) µ.Endpoint` builds the `Endpoint` that matches HTTP Verb. You supplies either a valid HTTP Verb or wildcard to match anything.

```go
e := µ.Join(µ.Method("GET"), /* ... */)
e(mock.Input())

// The library implements a syntax sugar for mostly used HTTP Verbs
e := µ.GET(/* ... */)
e(mock.Input())
```

**Match Path**

`func Path(arrows ...path.Arrow) µ.Endpoint` builds the `Endpoint` that matches URL path from HTTP request. The endpoint considers the path as an ordered sequence of segments, it takes a corresponding product of segment pattern matchers/extractors (they are defined in [`path`](../path/path.go) package).

```go
e := µ.Path(path.Is("foo"), path.Is("bar"))
e(mock.Input(mock.URL("/foo/bar")))
```

Often, implementation of **root** `Endpoint` is required, use `µ.Path` with empty definition.

```go
e := µ.Path()
e(mock.Input(mock.URL("/")))
```

Skip `µ.Path` definition to match all the segments (entire path of URL).

**Extract Path**

The library implements path extractors endpoints. They lift matched path segments to values of corresponding type. The extractor fails with `NoMatch` if segment value cannot be converted to requested type.
* `path.String`
* `path.Int`

```go
var bar string
e := µ.Path(path.Is("foo"), path.String(&bar))
e(mock.Input(mock.URL("/foo/bar")))
```

**Params**

A handing of query string params for HTTP request is consistent with matching/extracting path segments.

`func Param(arrows ...param.Arrow) µ.Endpoint` builds the `Endpoint` that matches URL query string from HTTP request. The endpoint considers a query params as a hashmap, it takes a product of params matchers/extractors (they are defined in [`param`](../param/param.go) package). Functions `param.Is` and `param.Any` matches query params; `param.String`, `param.MaybeString`, `param.Int` and `param.MaybeInt` extracts values.

```go
var text string
e := µ.Param(
  param.Is("foo", "bar"),
  param.String("q", &text),
)
e(mock.Input(mock.URL("/?foo=bar&q=text")))
```

**Headers**

A handing of HTTP headers is consistent with matching/extracting path segments.

`func Header(arrows ...header.Arrow) µ.Endpoint` builds the `Endpoint` that matches HTTP request headers. The endpoint considers headers as a hashmap, it takes a product of header matchers/extractors (they are defined in [`header`](../header/header.go) package). Functions `header.Is` and `header.Any` matches headers; `header.String`, `header.MaybeString`, `header.Int` and `header.MaybeInt` extracts values.

```go
var length int
e := µ.Header(
  header.Is("Content-Type", "application/json"),
  header.Int("Content-Length", &length),
)
e(mock.Input(
  mock.Header("Content-Type", "application/json"),
  mock.Header("Content-Length", "1024"),
))
```

**Bodies**

The library defines `Endpoint` to decode and extract body of HTTP request. It supports `µ.Text` and `µ.JSON`. The JSON endpoint does not match if `json.Unmarshal` returns error.

```go
type User struct {
  Username string `json:"username"` 
}

var user User
e := µ.Body(&user)
e(mock.Input(mock.Text("{\"username\":\"Joe Doe\"}")))
```

**Authentication with AWS Cognito**

The library defines a type `µ.AccessToken` and `func JWT(val *µ.AccessToken) µ.Endpoint` to extract JWT access token, which is provided by AWS Cognito service.


## High-order Endpoints

Usage of combinators is an essential part to declare API from primitive endpoints. The library defines `and-then` product and `or-else` coproduct combinators. They have been discussed earlier in this guide. Use combinators to implement high-order endpoints.

**Product endpoint**

Use the product combinator to declare *conjunctive conditions*. 

```go
// High Order Product Endpoint
//  /search?q=:text
func search(q *string) µ.Endpoint {
  return µ.Join(
    µ.Path(path.Is("search")),
    µ.Param(param.String("q", q))
  )
}

// Use HoC
var q string
µ.GET( search(&q) )
```

**Coproduct endpoint**

A co-product represents either-or endpoint evaluation.

```go
// High Order CoProduct Endpoint
//  /search?q=:text
//  /search/:text
func search(q *string) µ.Endpoint {
  return µ.Or(
    µ.Path(path.Is("search"), path.String(q)),
    µ.Join(
      µ.Path(path.Is("search")),
      µ.Param(param.String("q", q)),
    ),
  )
}

// Use HoC
var q string
µ.GET( search(&q) )
```

## Mapping Endpoints

A business logic is defined as Endpoint mapper with help of closure functions `Ø ⟼ Output`. The library provides `func FMap(f func() error) µ.Endpoint` function. It lifts a transformer into Endpoint so that it is composable with other Endpoints.

```go
µ.GET(
  µ.Path(path.Is("foo")),
  µ.FMap(func() error { µ.Ok() }),
)
```

## Outputs

Every returned value from the mapper/transformer is `Output`, which is implemented as `error` value. The library supplies [primitives](../output.go) to declare output of HTTP response. Endpoint *maps* the request either to successful HTTP status code or failure. The failures are RFC 7807: Problem Details for HTTP APIs.

The library provides factory functions named after HTTP status codes. Use them to declare your intent

```go
µ.GET(
  µ.Path(path.Is("foo")),
  µ.FMap(
    func() error {
      µ.Ok().
        With("Content-Type", "application/json").
        JSON(User{"Joe Doe"})
    }
  ),
)
```

## Unit testing

Gouildian support unit testing of API without a needs to spawn actual HTTP server.
Each `Endpoint` is a function, mock HTTP Input and validate its result.

```go
endpoint := µ.GET(/* ... */)

request := mock.Input(mock.URL("/foo"))

endpoint(request).Body == "{\"username\":\"Joe Doe\"}"
```
