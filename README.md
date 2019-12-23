# Gouldian

<img src="./doc/logo.svg" width="320" align="left"/>

The library is Go combinator library for building HTTP services.
The library is a thin layer of purely functional abstractions on top
of AWS Gateway API. It resolves a challenge of building simple and
declarative api implementations in the absence of pattern matching.

[![Documentation](https://godoc.org/github.com/fogfish/gouldian?status.svg)](http://godoc.org/github.com/fogfish/gouldian)
[![Build Status](https://secure.travis-ci.org/fogfish/gouldian.svg?branch=master)](http://travis-ci.org/fogfish/gouldian)
[![Git Hub](https://img.shields.io/github/last-commit/fogfish/gouldian.svg)](http://travis-ci.org/fogfish/gouldian)
[![Coverage Status](https://coveralls.io/repos/github/fogfish/gouldian/badge.svg?branch=master)](https://coveralls.io/github/fogfish/gouldian?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/fogfish/gouldian)](https://goreportcard.com/report/github.com/fogfish/gouldian)



## Inspiration

Microservices have become a design style to evolve system architecture in parallel, implement stable and consistent interfaces. An expressive language is required to design the variety of network interfaces. A pure functional languages fits very well to express communication behavior due they rich techniques to hide the networking complexity. [Finch](https://github.com/finagle/finch) is the best library in Scala for microservice development.

Gouldian is heavily inspired by [Finch](https://github.com/finagle/finch). However, it is primarily designed for serverless application to implement microservices using AWS Lambda and AWS API Gateway. 


## Getting started

The library requires **Go 1.13** or later due to usage of [new error interface](https://blog.golang.org/go1.13-errors).

The latest version of the library is available at `master` branch. All development, including new features and bug fixes, take place on the `master` branch using forking and pull requests as described in contribution guidelines.

Here is minimal "Hello World!" example that matches any HTTP requests
to `/hello` endpoint. You can run this example locally see the [instructions](example/hello-world). 

```go
package main

import (
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/fogfish/gouldian"
)

func main() {
  lambda.Start( gouldian.Serve(hello()) )
}

func hello() gouldian.Endpoint {
  return gouldian.Get().Path("hello").Then(
      func() error { return gouldian.Ok().Text("Hello World!") }
  )
}
```

See [example](example) folder for advanced use-case. The library  [api specification](http://godoc.org/github.com/fogfish/gouldian) is available via Go doc.



## Next steps

* Study [Endpoint](endpoint.go) type and its composition

* Check build-in [collection of endpoints](request.go) to deal with HTTP request. See types: [HTTP](http://godoc.org/github.com/fogfish/gouldian/#HTTP), [APIGateway](http://godoc.org/github.com/fogfish/gouldian/#APIGateway)

* Endpoint always returns some `Output` that defines HTTP response. There are three cases of output: HTTP Success, HTTP Failure and general error. See [Output](http://godoc.org/github.com/fogfish/gouldian/#Output), [Issue](http://godoc.org/github.com/fogfish/gouldian/#Issue) types.

* Learn about microservice deployment with AWS CDK.


## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


The build and testing process requires [Go](https://golang.org) version 1.13 or later.

**Build** and **run** in your development console.

```bash
git clone https://github.com/fogfish/golem
cd golem
go test -cover ./...
```

## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/gouldian.svg?style=for-the-badge)](LICENSE)