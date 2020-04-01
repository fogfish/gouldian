//
//   Copyright 2019 Dmitry Kolesnikov, All Rights Reserved
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package gouldian_test

import (
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/header"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/param"
	"github.com/fogfish/gouldian/path"
	"github.com/fogfish/it"
)

func TestVerbAny(t *testing.T) {
	endpoint := µ.ANY()

	success1 := mock.Input(mock.Method("GET"))
	success2 := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success1)).Should().Equal(nil).
		If(endpoint(success2)).Should().Equal(nil)
}

func TestVerbDelete(t *testing.T) {
	endpoint := µ.DELETE()
	success := mock.Input(mock.Method("DELETE"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbGet(t *testing.T) {
	endpoint := µ.GET()
	success := mock.Input(mock.Method("GET"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbPatch(t *testing.T) {
	endpoint := µ.PATCH()
	success := mock.Input(mock.Method("PATCH"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)

}

func TestVerbPost(t *testing.T) {
	endpoint := µ.POST()
	success := mock.Input(mock.Method("POST"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)

}

func TestVerbPut(t *testing.T) {
	endpoint := µ.PUT()
	success := mock.Input(mock.Method("PUT"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)

}

func TestPath(t *testing.T) {
	foo := µ.GET(µ.Path(path.Is("foo")))
	bar := µ.GET(µ.Path(path.Is("bar")))
	foobar := µ.GET(µ.Path(path.Is("foo"), path.Is("bar")))

	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestPathRoot(t *testing.T) {
	root := µ.GET(µ.Path())

	success := mock.Input(mock.URL("/"))
	failure := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(root(success)).Should().Equal(nil).
		If(root(failure)).ShouldNot().Equal(nil)
}

func TestParam(t *testing.T) {
	foo := µ.GET(µ.Param(param.Is("foo", "bar")))
	bar := µ.GET(µ.Param(param.Is("bar", "foo")))
	foobar := µ.GET(µ.Param(param.Is("foo", "bar"), param.Is("bar", "foo")))

	req := mock.Input(mock.URL("/?foo=bar"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestHeader(t *testing.T) {
	foo := µ.GET(µ.Header(header.Is("foo", "bar")))
	bar := µ.GET(µ.Header(header.Is("bar", "foo")))
	foobar := µ.GET(µ.Header(header.Is("foo", "bar"), header.Is("bar", "foo")))

	req := mock.Input(mock.Header("foo", "bar"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestJWT(t *testing.T) {
	var token µ.AccessToken
	foo := µ.GET(µ.JWT(&token))
	req := mock.Input(
		mock.Auth(
			map[string]interface{}{
				"sub":   "00000000-0000-0000-0000-000000000000",
				"scope": "a b",
			},
		),
	)

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(token.Sub).Should().Equal("00000000-0000-0000-0000-000000000000").
		If(token.Scope).Should().Equal("a b")
}

type foobar struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestBodyJSON(t *testing.T) {
	var value foobar
	foo := µ.GET(µ.Body(&value))
	success := mock.Input(mock.JSON(foobar{"foo", 10}))
	failure1 := mock.Input(
		mock.Header("Content-Type", "application/json"),
		mock.Text("foobar"),
	)
	failure2 := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(foobar{"foo", 10}).
		If(foo(failure1)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestBodyForm(t *testing.T) {
	var value foobar
	foo := µ.GET(µ.Body(&value))
	success := mock.Input(
		mock.Header("Content-Type", "application/x-www-form-urlencoded"),
		mock.Text("foo=foo&bar=10"),
	)
	failure1 := mock.Input(
		mock.Header("Content-Type", "application/x-www-form-urlencoded"),
		mock.Text("foobar"),
	)
	failure2 := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(foobar{"foo", 10}).
		If(foo(failure1)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestText(t *testing.T) {
	var value string
	foo := µ.GET(µ.Text(&value))
	success := mock.Input(mock.Text("foobar"))
	failure := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal("foobar").
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestFMapSuccess(t *testing.T) {
	foo := µ.GET(
		µ.Path(path.Is("foo")),
		µ.FMap(
			func() error { return µ.Ok().Text("bar") },
		),
	)
	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Assert(
		func(be interface{}) bool {
			return be.(error).Error() == "bar"
		},
	)
}

func TestFMap2Success(t *testing.T) {
	foo := µ.GET(
		µ.Path(path.Is("foo")),
		µ.FMap(func() error { return µ.Ok().Text("bar") }),
	)
	bar := µ.GET(
		µ.Path(path.Is("bar")),
		µ.FMap(func() error { return µ.Ok().Text("foo") }),
	)
	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(µ.Or(foo, bar)(req)).Should().Assert(
		func(be interface{}) bool {
			return be.(error).Error() == "bar"
		},
	)
}

func TestFMapFailure(t *testing.T) {
	foo := µ.GET(
		µ.Path(path.Is("foo")),
		µ.FMap(
			func() error { return µ.Unauthorized("") },
		),
	)
	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Assert(
		func(be interface{}) bool {
			return be.(error).Error() == "401: Unauthorized"
		},
	)
}
