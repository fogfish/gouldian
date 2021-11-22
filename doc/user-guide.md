# User Guide

- [Overview](#overview)
- [Composition](#composition)
- [Life-cycle](#life-cycle)
- [Primitive Endpoint types](#primitive-endpoint-types)
- [Primitive Endpoints](#primitive-endpoints)
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
* Golang standard HTTP server: [`http.Request`](https://pkg.go.dev/net/http)
* AWS API Gateway: [`APIGatewayProxyRequest`](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go).

`Output` is a sum type that captures a result of `Endpoint` evaluation. Technically, it indicates if:
* the endpoint do not match the request
* the endpoint matches the request
* the endpoint successfully transforms the request to HTTP output
* the endpoint has failed to transform the request

Golang is missing generics and type variance. Therefore, the Output is always an `error` value. The library could possible implement own interface but due to opaque error handling requirement, the interface behind `error` type is used. The library supplies [primitives](../output.go) to declare output using HTTP status codes notation (e.g. `Ok`, `Created`, `BadRequest`, `Unauthorized`, etc).

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

Entire HTTP service is built (see `Serve` combinator) as a co-product Endpoint that defined entire *"api algebra"* of application. Each incoming HTTP request passed to this *"algebra"* - root endpoint for the further evaluation:

```go
service := httpd.Serve(
  µ.GET(µ.Path("a"), /* ... */),
  µ.GET(µ.Path("b"), /* ... */),
  /* ... */
)
```

It is important to understand the life-cycle behavior for development of a [High-Order Endpoints](#high-order-endpoints) and writing a [Unit Testing](#unit-testing) in your application.

1. The library envelops each incoming request to `Input` type and applies it to the endpoint `service(input)`.
2. The resulting value of `error` (aka `Output`) type is matched against
* `NoMatch` causes abort of current *product* `Endpoint`. The request is passed to succeeding *co-product* `Endpoint`.
* `nil` continues evaluation of *product* `Endpoint` to succeeding item.
* `error` aborts the evaluation of the endpoint. The output error value is send to the caller


## Primitive Endpoint types

Each Endpoint is acting either as *pattern matching* or *value extractor*. Pattern matching compares a defined literal (constant) value with a corresponding term at HTTP request. It fails if term is not "equal" to specified value with `NoMatch` response.

```go
// For example, the endpoint uses pattern matching, it is only "matches"
// HTTP request containing URL /foo/bar?baz=zar
µ.GET(
  µ.Path("foo", "bar"),
  µ.Param("baz").To("zar"),
  // ...
)
```

Extractors matches corresponding terms and lift its values to the context so that api implementation can use the value to parametrize the api logic. The primitive extractors support `string`, `int` and `float64` data types. The extractor fails with `NoMatch` if term value cannot be converted to requested type.


```go
// For example, the endpoint uses extractors, it "matches" the HTTP request 
// containing URL /foo/{bar}?baz={zar}
µ.GET(
  µ.Path("foo", bar),
  µ.Param("baz").To(zar),
  // ...
)
```

Extractor uses lenses to inject decoded terms of HTTP request into the application type. Lenses are essential and core feature of the library that ensures type safety of HTTP services.

```go
/*

So far, the library has used a simplified definition of Endpoint

  Endpoint: Input ⟼ Output

from the type-safe perspective of api specification, each endpoint is implemented by function

  F[A, B]: A ⟼ B 

Therefore, endpoint needs to transform Input to A, apply function F and output type B.
*/
type A struct {
  Bar, Zar string
} 

/*

The optics package of this library implement decomposition of product type (struct) into lenses
*/
var bar, zar = optics.ForProduct2(A{})

µ.GET(
  // these lenses are passed to extractors 
  µ.Path("foo", bar),
  µ.Param("baz").To(zar),
  µ.FMap(func(ctx µ.Context) error {
    // the context contains matched value and can be decoded to value of type A
    var a A
    if err := ctx.Get(&a); err != nil {
      return µ.Status.BadRequest(µ.WithIssue(err))
    }
    // ...
  }),
),
```


## Primitive Endpoints

The library delivers set of built-in endpoints to deal with HTTP request processing.

**Match HTTP Verb/Method**

`func Method(verb string) µ.Endpoint` builds the `Endpoint` that matches HTTP Verb. You supplies either a valid HTTP Verb or wildcard (`µ.Any`) to match anything.

```go
e := µ.Join(µ.Method("GET"), /* ... */)
e(mock.Input())

// The library implements a syntax sugar for mostly used HTTP Verbs
e := µ.GET(/* ... */)
e(mock.Input())
```

**Match Path**

`func Path(segments ...interface{}) µ.Endpoint` builds the `Endpoint` that matches URL path from HTTP request. The endpoint considers the path as an ordered sequence of segments, it takes a sequence of either pattern matchers (literals) or extractors. 

```go
// sequence of pattern matchers (literals)
e := µ.Path("foo", "bar")
e(mock.Input(mock.URL("/foo/bar")))
```

Often, implementation of **root** `Endpoint` is required, use `µ.Path` with empty definition.

```go
e := µ.Path()
e(mock.Input(mock.URL("/")))
```

Skip `µ.Path` definition to match any path of the request.

**Extract Path**

The library uses lenses to lift matched path segments into the context. The extractor fails with `NoMatch` if segment value cannot be converted to requested type.

```go
type A struct {
  Bar string
}

var bar := optics.ForProduct1(A{})
e := µ.Path("foo", bar)
e(mock.Input(mock.URL("/foo/bar")))
```

**Params**

The library defines a type `Param` that builds the `Endpoint` to match URL query string from HTTP request. The type defines a functions `Is` and `Any` matches query params; `To` and `Maybe` to extracts values.

```go
e := µ.Param("foo").Is("bar")
e(mock.Input(mock.URL("/?foo=bar")))

e := µ.Param("foo").To(bar)
e(mock.Input(mock.URL("/?foo=bar")))
```

**Headers**

The library defines a type `Header` that builds the `Endpoint` to match URL query string from HTTP request. The type defines a functions `Is` and `Any` matches query params; `To` and `Maybe` to extracts values. See the package `headers` that defines HTTP header constants


```go

e := µ.Header("Content-Type").Is("application/json")
e(mock.Input(mock.Header("Content-Type", "application/json")))

e := µ.Header("Content-Length").To(length)
e(mock.Input(mock.Header("Content-Length", "1024")))
```

**Body**

The library defines `Body` function build `Endpoint` to extract value. The value is then decoded by lens into the application type. The following example decodes body into struct. 

```go
// application type that captures application payload
type User struct {
  Username string `json:"username"` 
}

// type of the request
type A struct {
  User User
}
var user := optics.ForProduct1(A{})

e := µ.Body(user)
e(mock.Input(mock.Text("{\"username\":\"Joe Doe\"}")))
```


**Authentication with AWS Cognito**

The library defines a types `µ.JWT`, `µ.Access` that builds the `Endpoint` to extract JWT access token. The type defines a functions `Is` matches fields of JWT token; `To` and `Maybe` to extracts values. The serverless factory of support automatic decoding of JWT access token, which is provided by AWS Cognito service.

```go
/*

Endpoint matches if HTTP request contains JWT created by AWS Cognito for user
*/ 
type A struct{ User string }

user := optics.Lenses1(MyT{})
e := µ.GET( µ.Access(µ.JWT.Username).To(user) )

/*

Endpoint matches if HTTP request contains JWT created by AWS Cognito for trusted client
*/ 
type A struct{ Client string }

client := optics.Lenses1(MyT{})
e := µ.GET( µ.Access(µ.JWT.ClientID).To(client) )
```


## High-order Endpoints

Usage of combinators is an essential part to declare API from primitive endpoints. The library defines `and-then` product and `or-else` coproduct combinators. They have been discussed earlier in this guide. Use combinators to implement high-order endpoints.

**Product endpoint**

Use the product combinator to declare *conjunctive conditions*. 

```go
// High Order Product Endpoint
//  /search?q=:text
func search(text optics.Lens) µ.Endpoint {
  return µ.Join(
    µ.Path("search"),
    µ.Param("q").To(text))
  )
}

// Use HoC
var text = optics.FromProduct1(A{})
µ.GET( search(text) )
```

**Coproduct endpoint**

A co-product represents either-or endpoint evaluation.

```go
// High Order CoProduct Endpoint
//  /search?q=:text
//  /search/:text
func search(text optics.Lens) µ.Endpoint {
  return µ.Or(
    µ.Path("search", text),
    µ.Join(
      µ.Path("search"),
      µ.Param("q").To(text)),
    ),
  )
}

// Use HoC
var text = optics.FromProduct1(A{})
µ.GET( search(text) )
```

## Mapping Endpoints

A business logic is defined as Endpoint mapper with help of closure functions `Context ⟼ Output`. The library provides `func FMap(f func(µ.Context) error) µ.Endpoint` function. It lifts a transformer into Endpoint so that it is composable with other Endpoints.

```go
µ.GET(
  µ.Path("foo"),
  µ.FMap(func(µ.Context) error { µ.Status.OK() }),
)
```

## Outputs

Every returned value from the mapper/transformer is `Output`, which is implemented as `error` value. The library supplies [primitives](../output.go) to declare output of HTTP response. Endpoint *maps* the request either to successful HTTP status code or failure. The failures are RFC 7807: Problem Details for HTTP APIs.

The library provides factory functions named after HTTP status codes. Use them to declare your intent

```go
µ.GET(
  µ.Path(path.Is("foo")),
  µ.FMap(
    func(µ.Context) error {
      return µ.Status.Ok(
        µ.WithJSON(User{"Joe Doe"}),
      ),
    },
  ),
)
```

## Unit testing

Gouildian support unit testing of API without a needs to spawn actual HTTP server. Each `Endpoint` is a function, mock HTTP Input and validate its result.

```go
endpoint := µ.GET(/* ... */)

request := mock.Input(mock.URL("/foo"))

switch v := endpoint(request).(type) {
  case *µ.Output:
    v.Body == "{\"username\":\"Joe Doe\"}"
  default:
    // error
}
```

The library also supports testing using standard test server

```go
import "net/http/httptest"

httptest.NewServer(
  httpd.Serve(
    µ.GET(/* ... */),
    µ.GET(/* ... */),
    /* ... */
  ),
)
```