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

The library is heavily inspired by Scala [Finch](https://github.com/finagle/finch). However, gouldian primary target is a serverless development with AWS Lambda and AWS Gateway API.



## Getting started

Here is minimal "Hello World!" example that matches any HTTP requests
to /hello endpoint.

```go
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

See [example](example) folder for advanced use-case and its [documentation](http://godoc.org/github.com/fogfish/gouldian)



## Next steps

* Study [Endpoint](endpoint.go) type and its composition
* Check build-in [collection of endpoints](request.go) to deal with HTTP request. See types: [HTTP](http://godoc.org/github.com/fogfish/gouldian/#HTTP), [APIGateway](http://godoc.org/github.com/fogfish/gouldian/#APIGateway)
* Endpoint always returns some `Output` that defines HTTP response. There are three cases of output: HTTP Success, HTTP Failure and general error. See [Output](http://godoc.org/github.com/fogfish/gouldian/#Output), [Issue](http://godoc.org/github.com/fogfish/gouldian/#Issue) types.



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