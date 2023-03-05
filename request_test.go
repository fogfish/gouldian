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

func TestMethodNoMatch(t *testing.T) {
	spec := []struct {
		Verb func(µ.Routable, ...µ.Endpoint) µ.Routable
		Mock mock.Mock
	}{
		{µ.GET, mock.Method("OTHER")},
		{µ.PUT, mock.Method("OTHER")},
		{µ.POST, mock.Method("OTHER")},
		{µ.DELETE, mock.Method("OTHER")},
		{µ.PATCH, mock.Method("OTHER")},
	}

	for _, tt := range spec {
		foo := mock.Endpoint(
			tt.Verb(
				µ.URI(µ.Path("foo")),
			),
		)

		req := mock.Input(tt.Mock, mock.URL("/foo"))
		err := foo(req)
		it.Then(t).ShouldNot(
			it.Nil(err),
		)
	}
}

func TestBodyJson(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	spec := []struct {
		Mock   *µ.Context
		Expect foobar
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/json"),
				mock.JSON(foobar{"foo1", 10}),
			),
			foobar{"foo1", 10},
		},
		{
			mock.Input(
				mock.Header("Content-Type", "application/json"),
				mock.Text(`{"foo":"foo1","bar":10}`),
			),
			foobar{"foo1", 10},
		},
	}

	type request struct {
		FooBar foobar `content:"json"`
	}
	var lens = µ.Optics1[request, foobar]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).Should(
			it.Nil(err),
			it.Nil(µ.FromContext(tt.Mock, &req)),
			it.Equiv(req.FooBar, tt.Expect),
		)
	}
}

func TestBodyJsonNoMatch(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	spec := []struct {
		Mock *µ.Context
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/json"),
				mock.Text(`{"foo:"foo1,"bar":10}`),
			),
		},
	}

	type request struct {
		FooBar foobar `content:"json"`
	}
	var lens = µ.Optics1[request, foobar]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).
			Should(it.Nil(err)).
			ShouldNot(it.Nil(µ.FromContext(tt.Mock, &req)))
	}
}

func TestBodyForm(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	spec := []struct {
		Mock   *µ.Context
		Expect foobar
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/x-www-form-urlencoded"),
				mock.Text("foo=foo1&bar=10"),
			),
			foobar{"foo1", 10},
		},
	}

	type request struct {
		FooBar foobar `content:"form"`
	}
	var lens = µ.Optics1[request, foobar]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).Should(
			it.Nil(err),
			it.Nil(µ.FromContext(tt.Mock, &req)),
			it.Equiv(req.FooBar, tt.Expect),
		)
	}
}

func TestBodyFormNoMatch(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	spec := []struct {
		Mock *µ.Context
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "application/x-www-form-urlencoded"),
				mock.Text("foobar"),
			),
		},
	}

	type request struct {
		FooBar foobar `content:"form"`
	}
	var lens = µ.Optics1[request, foobar]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).
			Should(it.Nil(err)).
			ShouldNot(it.Nil(µ.FromContext(tt.Mock, &req)))
	}
}

func TestBodyText(t *testing.T) {
	spec := []struct {
		Mock   *µ.Context
		Expect string
	}{
		{
			mock.Input(
				mock.Header("Content-Type", "text/plain"),
				mock.Text("foobar"),
			),
			"foobar",
		},
	}

	type request struct {
		FooBar string
	}
	var lens = µ.Optics1[request, string]()

	for _, tt := range spec {
		var req request
		foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
		err := foo(tt.Mock)

		it.Then(t).Should(
			it.Nil(err),
			it.Nil(µ.FromContext(tt.Mock, &req)),
			it.Equiv(req.FooBar, tt.Expect),
		)
	}
}
