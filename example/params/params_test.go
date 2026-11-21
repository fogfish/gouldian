package main

import (
	"net/url"
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/example/internal/exampletest"
)

func TestExample(t *testing.T) {
	exampletest.Run(t, []µ.Routable{qString(), qInt()}, []exampletest.Case{
		{Name: "string", Path: "/echo", Query: url.Values{"q": {"value"}}, Contains: "query string (value)"},
		{Name: "int", Path: "/echo", Query: url.Values{"v": {"12345"}}, Contains: "query number (12345)"},
		{Name: "malformed int", Path: "/echo", Query: url.Values{"v": {"abc"}}, Status: 400, Contains: "v"},
		{Name: "no params", Path: "/echo", Status: 404},
	})
}
