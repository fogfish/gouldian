package gouldian_test

import (
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it/v2"
	"testing"
)

func TestMethod(t *testing.T) {
	spec := []struct {
		Verb func(µ.Routable, ...µ.Endpoint) µ.Routable
		Mock mock.Mock
	}{
		{µ.GET, mock.Method("GET")},
		{µ.PUT, mock.Method("PUT")},
		{µ.POST, mock.Method("POST")},
		{µ.DELETE, mock.Method("DELETE")},
		{µ.PATCH, mock.Method("PATCH")},
		{µ.ANY, mock.Method("GET")},
		{µ.ANY, mock.Method("PUT")},
	}

	for _, tt := range spec {
		foo := mock.Endpoint(
			tt.Verb(
				µ.URI(µ.Path("foo")),
			),
		)

		req := mock.Input(tt.Mock, mock.URL("/foo"))
		err := foo(req)
		it.Then(t).Should(
			it.Nil(err),
		)
	}
}
