package main

import (
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/example/internal/exampletest"
)

func TestExample(t *testing.T) {
	exampletest.Run(t, []µ.Routable{text(), json()}, []exampletest.Case{
		{Name: "text", Path: "/echo", Header: map[string]string{"Accept": "text/plain"}, Contains: "hello world."},
		{Name: "json", Path: "/echo", Header: map[string]string{"Accept": "application/json"}, Contains: `"text"`},
		{Name: "unacceptable", Path: "/echo", Header: map[string]string{"Accept": "image/png"}, Status: 406},
	})
}
