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

An `core.Endpoint` is a key abstraction in the framework. It is a *pure function* that takes HTTP request as Input and return Output, result of request evaluation.

```go
/*

Endpoint: Input ⟼ Output
*/
type Endpoint func(*Input) error
```

`Input` is a convenient wrapper of HTTP request with some Gouldian specific context. This library release support only [`APIGatewayProxyRequest`](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go).

`Output` is a sum type that captures a result of `Endpoint` evaluation. Technically, it indicates if it:
* do not match the request
* matches the request
* successfully transforms the request to HTTP output
* fails to apply request transformation

Golang is missing generics and type variance. Therefore, the Output of HTTP request evaluation is always an `error` value. The library supplies [primitives](../output.go) to declare successful output, they resembles HTTP status codes (e.g. `Ok`, `Created`, `BadRequest`, `Unauthorized`, etc).

## Composition

Any `Endpoint A` can be composed with `Endpoint B` into new `Endpoint C`. The library supports two combinators: and-then, or-else.

**and-then**

Use `and-then` to build product Endpoint: `A × B ⟼ C`. The product type matches Input if each composed function successfully matches it.

You can either chain endpoints with `Then` function

```go
var a: core.Endpoint = /* ... */
var b: core.Endpoint = /* ... */

//
// You can use `Then` method declared for the Endpoint type
c := a.Then(b)

//
// Alternatively, variadic function `Join` does same for sequence of Endpoints
c := core.Join(a, b)
```

**or-else**

Use `or-else` to build co-product Endpoint: `A ⨁ B ⟼ C`. The co-product is also known as sum-type that matches the first successful function.

```go
var a: core.Endpoint = /* ... */
var b: core.Endpoint = /* ... */

//
// You can use `Or` method declared for the Endpoint type
c := a.Or(b)

//
// Alternatively, variadic `Or` variant does same for sequence of Endpoints
c := core.Or(a, b)
```

These rules of Endpoint composition allows you to build any complex HTTP request matching login from small re-usable block declared in this library and also defined reusable high-order Endpoint specific for you application. 


## Life-cycle

The function `gouldian.Serve` builds a co-product Endpoint to define entire *"api algebra"* of application. Each incoming HTTP request passed to this Endpoint. It is important to understand the life-cycle behavior for development of a [High-Order Endpoints](#high-order-endpoints) and writing a [Unit Testing](#unit-testing) in your application.

2. The library envelops each incoming request to `Input` type and applies it to the endpoint `api(input)`.
3. The resulting value of `error` (aka `Output`) type is matched against
  a. `NoMatch` causes abort of current *product* `Endpoint`. The request is passed to succeeding *co-product* `Endpoint`.
  b. `nil` continues evaluation of *product* `Endpoint` to succeeding item.
  c. `error` aborts the evaluation of the application. The output error value is output to HTTP client

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

`func Method(verb string) core.Endpoint` builds an `Endpoint` that matches HTTP Verb. You either supplies a valid HTTP Verb or wildcard to match anything.

```go
e := core.Join(µ.Method("GET"), /* ... */)
e(mock.Input())

// The library implements a syntax sugar for mostly used HTTP Verbs
e := µ.GET(/* ... */)
e(mock.Input())
```

**Match Path**

`func Path(arrows ...path.Arrow) core.Endpoint` builds an `Endpoint` that matches arbitrary URL path from HTTP request. The endpoint considers a path as a sequence of segments, it takes a corresponding product of segment pattern matchers/extractors (they are defined in [`path`](../path/path.go) package).

```go
e := µ.Path(path.Is("foo"))
e(mock.Input(mock.URL("/foo")))
```

Often, implementation of **root** `Endpoint` is required, use `µ.Path` with empty definition.

```go
e := µ.Path()
e(mock.Input(mock.URL("/")))
```

Skip `µ.Path` definition to match all the segments, entire path of URL.

**Extract Path**

There are path extractors endpoints that lifts a matched path segment to value of a requested type. The extractor fails with `NoMatch` if segment value cannot be converted to requested type.
* `path.String`
* `path.Int`

```go
var bar string
e := µ.Path(path.Is("foo"), path.String(&bar))
e(mock.Input(mock.URL("/foo/bar")))
```

**Params**

A handing of query string params for HTTP request is consistent with matching/extraction of path segments.

`func Param(arrows ...param.Arrow) core.Endpoint` builds an `Endpoint` that matches URL query string from HTTP request. The endpoint considers a query params as a hashmap, it takes a product of params matchers/extractors (they are defined in [`param`](../param/param.go) package). Functions `param.Is` and `param.Any` matches query params; `param.String`, `param.MaybeString`, `param.Int` and `param.MaybeInt` extracts values.

```go
var text string
e := µ.Param(
  param.Is("foo", "bar"),
  param.String("q", &text),
)
e(mock.Input(mock.URL("/?foo=bar&q=text")))
```

**Headers**

A handing of query string params for HTTP request is consistent with matching/extraction of path segments.

`func Header(arrows ...header.Arrow) core.Endpoint` builds an `Endpoint` that matches HTTP request headers. The endpoint considers headers as a hashmap, it takes a product of header matchers/extractors (they are defined in [`header`](../header/header.go) package). Functions `header.Is` and `header.Any` matches headers; `header.String`, `header.MaybeString`, `header.Int` and `header.MaybeInt` extracts values.

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

The library defines `Endpoint` to decode and extract body of HTTP request. It supports `µ.Text` and `µ.JSON`. The JSON endpoint would not match is `json.Unmarshal` returns error.

```go
type User struct {
  Username string `json:"username"` 
}

var user User
e := µ.Body(&user)
e(mock.Input(mock.Text("{\"username\":\"Joe Doe\"}")))
```

**Authentication with AWS Cognito**

The library defines a type `core.AccessToken` and `func AccessToken(val *core.AccessToken) core.Endpoint` to extract JWT access token as it is provided by AWS Cognito service.


## High-order Endpoints

Usage of combinators is an essential part to specify API from primitive endpoints. The library define `and-then` product and `or-else` coproduct combinators. They have been discussed earlier in this guide. These combinator is an essential tools to declare re-usage high-order endpoints

**Product endpoint**

Use the product combinator to declare *conjunctive conditions*. All variadic functions implemented by the library are product.

```go
// High Order Product Endpoint
//  /search?q=:text
func search(q *string) core.Endpoint {
  return core.Join(
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
func search(q *string) core.Endpoint {
  return core.Or(
    µ.Path(path.Is("search"), path.String(q)),
    core.Join(
      µ.Path(path.Is("search")),
      µ.Param(param.String("q", q)),
    ),
  )
}


// Endpoint coproduct
either := core.Or(µ.Path(/*...*/), µ.Path(/*...*/))

// Path coproduct matches /foo and /bar
either := µ.Path(path.Or(path.Is("foo"), path.Is("bar"))
```

## Mapping Endpoints

A business logic is defined as Endpoint mapper with help of closure functions `Ø ⟼ Output`. The library provides `func FMap(f func() error) core.Endpoint` function. It lifts a transformer into Endpoint so that it is composable with other Endpoints.

```go
µ.GET(
  µ.Path(path.Is("foo")),
  µ.FMap(func() error { µ.Ok() }),
)
```

## Outputs

Every returned value from the mapper/transformer is `Output`, which is implemented as `error` value. The library supplies [primitives](../output.go) to declare output of HTTP response. You *maps* the request either to successful HTTP status code or failure. The failures are RFC 7807: Problem Details for HTTP APIs.

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
