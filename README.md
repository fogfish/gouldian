<p align="center">
  <img src="./doc/gouldian-v2.svg" height="120" />
  <h3 align="center">µ-Gouldian</h3>
  <p align="center"><strong>Go combinator library for building containerized and serverless HTTP services.</strong></p>

  <p align="center">
    <!-- Version -->
    <a href="https://github.com/fogfish/gouldian/releases">
      <img src="https://img.shields.io/github/v/tag/fogfish/gouldian?label=version" />
    </a>
    <!-- Documentation -->
    <a href="https://pkg.go.dev/github.com/fogfish/gouldian">
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
    <!-- Go Card -->
    <a href="https://goreportcard.com/report/github.com/fogfish/gouldian">
      <img src="https://goreportcard.com/badge/github.com/fogfish/gouldian" />
    </a>
    <!-- Maintainability -->
    <a href="https://codeclimate.com/github/fogfish/gouldian/maintainability">
      <img src="https://api.codeclimate.com/v1/badges/633dc8add2dfc0e7f88e/maintainability" />
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
* The library support opaque migration of HTTP service between traditional, containers and serverless environments. The api implementation remains source compatible regardless the execution environment;  
* The library enforces a type safe, pattern-based approach for api definition.
* Fast, zero allocation routing 


## Installing

The library requires **Go 1.13** or later due to usage of [new error interface](https://blog.golang.org/go1.13-errors).

The latest version of the library is available at `main` branch. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

1. Use `go get` to retrieve the library and add it as dependency to your application.

```bash
go get -u github.com/fogfish/gouldian
```

2. Import it in your code

```go
import (
  µ "github.com/fogfish/gouldian"
)
```

## Quick Example

Here is minimal "Hello World!" example that matches any HTTP requests
to `/hello` endpoint. You can run this example locally see the [instructions](example). 

```go
package main

import (
  µ "github.com/fogfish/gouldian"
  "github.com/fogfish/gouldian/server/httpd"
  "net/http"
)

func main() {
  http.ListenAndServe(":8080",
    httpd.Serve(hello()),
  )
}

func hello() µ.Routable {
  return µ.GET(
    µ.Path("hello"),
    µ.FMap(func(ctx µ.Context) error {
      return µ.Status.OK(µ.WithText("Hello World!"))
    }),
  )
}
```

## Benchmark

The library uses [go-http-routing-benchmark](https://github.com/julienschmidt/go-http-routing-benchmark) methodology for benchmarking, using structure of GitHub API as primary benchmark. The results are obtained on the reference hardware such as AWS m6i.large and a1.large instances.

**m6i.large** 3.5 GHz 3rd generation Intel Xeon Scalable processors:
* It takes **10.9M** routing decisions per second, taking about **110 ns/op** and consuming about **0 allocs/ops**.
* It performs **7.4M** requests/responses for the single endpoint with one parameter, taking about **162 ns/op** and consuming **24 B/op** with **2 allocs/op**.

**a1.large** AWS Graviton Processor with 64-bit Arm Neoverse cores:
* It takes **2M** routing decisions per second, taking about **520 ns/op** and consuming about **0 allocs/ops**.
* It performs **1.5M** requests/responses for the single endpoint with one parameter, taking about **763 ns/op** and consuming **24 B/op** with **2 allocs/op**.


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


The build and testing process requires [Go](https://golang.org) version 1.13 or later.

**Build** and **run** in your development console.

```bash
git clone https://github.com/fogfish/gouldian
cd gouldian
go test
go test -run=^$ -bench=. -cpu 1
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/gouldian/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 

## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/gouldian.svg?style=for-the-badge)](LICENSE)
