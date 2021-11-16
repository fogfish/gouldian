/*

  Copyright 2019 Dmitry Kolesnikov, All Rights Reserved

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

/*

Package gouldian is Go combinator library for building HTTP services.
The library is a thin layer of purely functional abstractions on top
of AWS Gateway API. It resolves a challenge of building simple and
declarative api implementations in the absence of pattern matching.


Inspiration

The library is heavily inspired by Scala Finch
https://github.com/finagle/finch. However, gouldian primary target is
a serverless development with AWS Lambda and AWS API Gateway.


Getting started

Here is minimal "Hello World!" example that matches any HTTP requests
to /hello endpoint.

  package main

  import (
    "github.com/aws/aws-lambda-go/lambda"
    µ "github.com/fogfish/gouldian"
  )

  func main() {
    lambda.Start( µ.Serve(hello()) )
  }

  func hello() µ.Endpoint {
    return µ.GET(µ.Path("hello")).FMap(
      func() error { return µ.Ok().Text("Hello World!") },
    )
  }

See examples folder for advanced use-case.

Next steps

↣ Study Endpoint type and its composition

↣ Check build-in collection of endpoints to deal with HTTP request.
See types: HTTP, APIGateway

↣ Endpoint always returns some `Output` that defines HTTP response.
There are three cases of output: HTTP Success, HTTP Failure and general
error. See Output, Issue types.

*/
package gouldian
