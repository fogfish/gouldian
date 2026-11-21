<p align="center">
  <img src="./doc/gouldian-v2.svg" height="120" />
  <h3 align="center">µ-Gouldian</h3>
  <p align="center"><strong>A typed, declarative HTTP router for Go that runs the same code on net/http and AWS Lambda. 87 ns per request on a 203-route API, zero allocations, within 1.4× of the fastest router, while binding parameters into structs.</strong></p>
  <p align="center">You don't pay for correctness: 404/405/406/415, RFC 9457 problems and safe defaults are built in.</p>

  <p align="center">
    <!-- Version -->
    <a href="https://github.com/fogfish/gouldian/releases">
      <img src="https://img.shields.io/github/v/tag/fogfish/gouldian?label=version" />
    </a>
    <!-- Documentation -->
    <a href="https://pkg.go.dev/github.com/fogfish/gouldian/v3">
      <img src="https://pkg.go.dev/badge/github.com/gouldian/dynamo" />
    </a>
    <!-- Build Status  -->
    <a href="https://github.com/fogfish/gouldian/actions/">
      <img src="https://github.com/fogfish/gouldian/workflows/build/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/fogfish/gouldian">
      <img src="https://img.shields.io/github/last-commit/fogfish/gouldian.svg" />
    </a>
    <!-- Coverage -->
    <a href="https://coveralls.io/github/fogfish/gouldian?branch=main">
      <img src="https://coveralls.io/repos/github/fogfish/gouldian/badge.svg?branch=main" />
    </a>
  </p>
</p>

--- 

The library is a thin layer of purely functional abstractions to build HTTP services. In the contrast with other HTTP routers, the library resolves a challenge of building simple and declarative api implementations in the absence of pattern matching at Golang. The library also support opaque migration of HTTP service between traditional, containers and serverless environments.

[User Guide](./doc/user-guide.md) |
[Hello World](./example/helloworld/helloworld.go) |
[Other Examples](./example/) |
[Benchmark](#benchmark)


## Inspiration

Microservices have become a design style to evolve system architecture in parallel, implement stable and consistent interfaces within distributed system. An expressive language is required to design the manifold of network interfaces. A pure functional languages fits very well to express communication behavior due they rich techniques to hide the networking complexity. [Finch](https://github.com/finagle/finch) is the best library in Scala for microservice development. Gouldian is heavily inspired by Finch. 

The library solves few practical problems of HTTP service development in Golang:
* The library support opaque migration of HTTP service between traditional, containers and serverless environments. The api implementation remains source compatible regardless the execution environment; The library supports integration with
  * Golang standard HTTP server: [`http.Request`](https://pkg.go.dev/net/http)
  * AWS API Gateway REST API (payload 1.0): [`APIGatewayProxyRequest`](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go)
  * AWS API Gateway HTTP API and Lambda Function URL (payload 2.0): `APIGatewayV2HTTPRequest`
* The library enforces a type safe, pattern-based approach for api definition: handlers receive a typed struct and never touch the raw request.
* A request has three outcomes, not two: *not my route* (`404`/`405`/`406`/`415`), *my route, bad input* (`400` RFC 9457 problem naming the invalid value) and *success*.
* Fast, zero allocation routing 


## Installing

The library requires **Go 1.27** or later. v3 is a major version, see [MIGRATION-v3.md](MIGRATION-v3.md) if you come from v2.

1. Use `go get` to retrieve the library and add it as dependency to your application.

```bash
go get -u github.com/fogfish/gouldian/v3
```

2. Import it in your code

```go
import (
  µ "github.com/fogfish/gouldian/v3"
)
```

## Quick Example

Here is minimal "Hello World!" example that matches any HTTP requests
to `/hello` endpoint. You can run this example locally see the [instructions](example). 

```go
package main

import (
  "net/http"

  µ "github.com/fogfish/gouldian/v3"
  ø "github.com/fogfish/gouldian/v3/output"
  "github.com/fogfish/gouldian/v3/server/httpd"
)

func main() {
  http.ListenAndServe(":8080",
    httpd.Serve(hello()),
  )
}

func hello() µ.Routable {
  return µ.GET(
    µ.URI(µ.Path("hello")),
    func(ctx *µ.Context) error {
      return ø.Status.OK(
        ø.ContentType.TextPlain,
        ø.Send("Hello World!"),
      )
    },
  )
}
```

The same routes run in AWS Lambda:

```go
api := µ.NewServer(µ.WithCORS(&µ.CORS{Origins: []string{"https://app.example"}})).Mount(hello())

lambda.Start(apigatewayv2.Handler(api)) // HTTP API, Function URL
lambda.Start(apigateway.Handler(api))   // REST API
http.ListenAndServe(":8080", httpd.Handler(api))
```

## Benchmark

The library uses [go-http-routing-benchmark](https://github.com/julienschmidt/go-http-routing-benchmark) methodology for benchmarking, using structure of GitHub API as primary benchmark. See [bench](bench) for the v2 baseline, the v3 results and the comparison script (`python3 bench/compare.py`).


## Next steps

* Study [User Guide](doc/user-guide.md).

* Check build-in collection of endpoints to deal with HTTP request: [path](path.go), [query param](param.go), [http header](header.go), [body and other](request.go) 

* Endpoint always returns some `Output` that defines HTTP response. There are three cases of output: HTTP Success, HTTP Failure and general error. See [Output](output.go) type.

* See [example](example) folder for other advanced use-case. 

* Learn about microservice deployment with AWS CDK, in case of serverless development


## How To Contribute

The library is [Apache 2.0](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


The build and testing process requires [Go](https://golang.org) version 1.27 or later.

**Build** and **run** in your development console.

```bash
git clone https://github.com/fogfish/gouldian
cd gouldian
go test -race ./...
go test -run=^$ -bench=. -cpu 1
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/gouldian/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 

## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/gouldian.svg?style=for-the-badge)](LICENSE)
