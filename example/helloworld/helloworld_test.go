package main

import (
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/example/internal/exampletest"
)

func TestExample(t *testing.T) {
	exampletest.Run(t, []µ.Routable{hello()}, []exampletest.Case{
		{Name: "hello", Path: "/hello", Contains: "Hello World!"},
		{Name: "method", Method: "POST", Path: "/hello", Status: 405},
	})
}
