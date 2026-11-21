package main

import (
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/example/internal/exampletest"
)

func TestExample(t *testing.T) {
	exampletest.Run(t, []µ.Routable{text(), json(), form()}, []exampletest.Case{
		{Name: "text", Method: "POST", Path: "/echo", Header: map[string]string{"Content-Type": "text/plain"}, Body: "Hello World.", Contains: "Hello World."},
		{Name: "json", Method: "POST", Path: "/echo", Header: map[string]string{"Content-Type": "application/json"}, Body: `{"string": "Hello World.", "number": 1010}`, Contains: "1010"},
		{Name: "form", Method: "POST", Path: "/echo", Header: map[string]string{"Content-Type": "application/x-www-form-urlencoded"}, Body: "string=Hello+World.&number=1010", Contains: "1010"},
		{Name: "bad json", Method: "POST", Path: "/echo", Header: map[string]string{"Content-Type": "application/json"}, Body: "{", Status: 400},
		{Name: "unsupported", Method: "POST", Path: "/echo", Header: map[string]string{"Content-Type": "image/png"}, Body: "x", Status: 415},
		{Name: "empty body", Method: "POST", Path: "/echo", Header: map[string]string{"Content-Type": "text/plain"}, Status: 400, Contains: "body is required"},
	})
}
