package main

import (
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/example/internal/exampletest"
)

func TestExample(t *testing.T) {
	exampletest.Run(t, []µ.Routable{echo()}, []exampletest.Case{
		{Name: "echo", Path: "/echo/hello", Header: map[string]string{"Accept": "text/plain"}, Contains: "hello"},
		{Name: "not acceptable", Path: "/echo/hello", Header: map[string]string{"Accept": "image/png"}, Status: 406},
		{Name: "not found", Path: "/other", Status: 404},
	})
}
