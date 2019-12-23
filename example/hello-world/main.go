package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fogfish/gouldian"
)

func main() {
	lambda.Start(gouldian.Serve(hello()))
}

func hello() gouldian.Endpoint {
	return gouldian.Get().Path("hello").FMap(
		func() error { return gouldian.Ok().Text("Hello World!") },
	)
}
