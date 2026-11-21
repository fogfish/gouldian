package main

import (
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/example/internal/exampletest"
)

func TestExample(t *testing.T) {
	exampletest.Run(t, []µ.Routable{get(), post(), put(), patch(), delete()}, []exampletest.Case{
		{Name: "get", Path: "/echo", Contains: "GET echo"},
		{Name: "post", Method: "POST", Path: "/echo", Contains: "POST echo"},
		{Name: "put", Method: "PUT", Path: "/echo", Contains: "PUT echo"},
		{Name: "patch", Method: "PATCH", Path: "/echo", Contains: "PATCH echo"},
		{Name: "delete", Method: "DELETE", Path: "/echo", Contains: "DELETE echo"},
		{Name: "head", Method: "HEAD", Path: "/echo"},
		{Name: "not allowed", Method: "TRACE", Path: "/echo", Status: 405},
	})
}
