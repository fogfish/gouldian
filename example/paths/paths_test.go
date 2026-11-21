package main

import (
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/example/internal/exampletest"
)

func TestExample(t *testing.T) {
	exampletest.Run(t, []µ.Routable{root(), path(), file()}, []exampletest.Case{
		{Name: "root", Path: "/", Contains: "matches root (/)"},
		{Name: "path", Path: "/path/value", Contains: "matches /path/value"},
		{Name: "file", Path: "/file/so/me/va/lue", Contains: "matches /file/so/me/va/lue"},
		{Name: "unknown", Path: "/path", Status: 404},
	})
}
