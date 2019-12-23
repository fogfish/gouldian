// Package gouldian is Go combinator library for building HTTP services.
// The library is a thin layer of purely functional abstractions on top
// of AWS Gateway API. It resolves a challenge of building simple and
// declarative api implementations in the absence of pattern matching.
//
//
// Inspiration
//
// The library is heavily inspired by Scala Finch
// https://github.com/finagle/finch. However, gouldian primary target is
// a serverless development with AWS Lambda and AWS Gateway API.
//
//
// Getting started
//
// Here is minimal "Hello World!" example that matches any HTTP requests
// to /hello endpoint.
//
//   import (
//     "github.com/aws/aws-lambda-go/lambda"
//     "github.com/fogfish/gouldian"
//   )
//
//   func main() {
//     lambda.Start( gouldian.Serve(hello()) )
//   }
//
//   func hello() gouldian.Endpoint {
//     return gouldian.Get().Path("hello").Then(
//        func() error { return gouldian.Ok().Text("Hello World!") }
//     )
//   }
//
// See examples folder for advanced use-case.
//
// Next steps
//
// ↣ Study Endpoint type and its composition
//
// ↣ Check build-in collection of endpoints to deal with HTTP request.
// See types: HTTP, APIGateway
//
// ↣ Endpoint always returns some `Output` that defines HTTP response.
// There are three cases of output: HTTP Success, HTTP Failure and general
// error. See Output, Issue types.
package gouldian
