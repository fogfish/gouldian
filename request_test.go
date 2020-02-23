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
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/path"
	"github.com/fogfish/it"
)

func TestVerbDelete(t *testing.T) {
	endpoint := µ.DELETE()
	success := mock.Input(mock.Method("DELETE"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint.IsMatch(success)).Should().Equal(true).
		If(endpoint.IsMatch(failure)).Should().Equal(false)
}

func TestVerbGet(t *testing.T) {
	endpoint := µ.GET()
	success := mock.Input(mock.Method("GET"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint.IsMatch(success)).Should().Equal(true).
		If(endpoint.IsMatch(failure)).Should().Equal(false)

}

func TestVerbPatch(t *testing.T) {
	endpoint := µ.PATCH()
	success := mock.Input(mock.Method("PATCH"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint.IsMatch(success)).Should().Equal(true).
		If(endpoint.IsMatch(failure)).Should().Equal(false)

}

func TestVerbPost(t *testing.T) {
	endpoint := µ.POST()
	success := mock.Input(mock.Method("POST"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint.IsMatch(success)).Should().Equal(true).
		If(endpoint.IsMatch(failure)).Should().Equal(false)

}

func TestVerbPut(t *testing.T) {
	endpoint := µ.PUT()
	success := mock.Input(mock.Method("PUT"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint.IsMatch(success)).Should().Equal(true).
		If(endpoint.IsMatch(failure)).Should().Equal(false)

}

func TestPath(t *testing.T) {
	foo := µ.GET(µ.Path(path.Is("foo")))
	bar := µ.GET(µ.Path(path.Is("bar")))
	foobar := µ.GET(µ.Path(path.Is("foo"), path.Is("bar")))

	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo.IsMatch(req)).Should().Equal(true).
		If(bar.IsMatch(req)).Should().Equal(false).
		If(foobar.IsMatch(req)).Should().Equal(false)
}

/*
func TestString(t *testing.T) {
	req := gouldian.Mock("/foo/bar")
	foo := ""

	it.Ok(t).
		If(gouldian.Get().Path("foo").String(&foo).IsMatch(req)).
		Should().Equal(true).
		//
		If(foo).Should().Equal("bar").
		//
		If(gouldian.Get().Path("foo").Path("bar").String(&foo).IsMatch(req)).
		Should().Equal(false)
}

func TestInt(t *testing.T) {
	req := gouldian.Mock("/foo/10")
	inv := gouldian.Mock("/foo/bar")
	foo := 0

	it.Ok(t).
		If(gouldian.Get().Path("foo").Int(&foo).IsMatch(req)).
		Should().Equal(true).
		//
		If(foo).Should().Equal(10).
		//
		If(gouldian.Get().Path("foo").Path("bar").Int(&foo).IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").Int(&foo).IsMatch(inv)).
		Should().Equal(false)
}

func TestParam(t *testing.T) {
	req := gouldian.Mock("/foo?bar=foo")

	it.Ok(t).
		If(gouldian.Get().Path("foo").Param("bar", "foo").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").Param("bar", "bar").IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").Param("foo", "").IsMatch(req)).
		Should().Equal(false)
}

func TestHasParam(t *testing.T) {
	req := gouldian.Mock("/foo?bar")
	foo := gouldian.Mock("/foo?bar=foo")

	it.Ok(t).
		If(gouldian.Get().Path("foo").HasParam("bar").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").HasParam("bar").IsMatch(foo)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").HasParam("foo").IsMatch(req)).
		Should().Equal(false)
}

func TestQString(t *testing.T) {
	req := gouldian.Mock("/foo?bar=foo")
	bar := ""

	it.Ok(t).
		If(gouldian.Get().Path("foo").QString("bar", &bar).IsMatch(req)).
		Should().Equal(true).
		If(bar).Should().Equal("foo").
		//
		If(gouldian.Get().Path("foo").QString("foo", &bar).IsMatch(req)).
		Should().Equal(true).
		If(bar).Should().Equal("")
}

func TestQInt(t *testing.T) {
	req := gouldian.Mock("/foo?bar=10")
	inv := gouldian.Mock("/foo?bar=foo")
	bar := 0

	it.Ok(t).
		If(gouldian.Get().Path("foo").QInt("bar", &bar).IsMatch(req)).
		Should().Equal(true).
		If(bar).Should().Equal(10).
		//
		If(gouldian.Get().Path("foo").QInt("foo", &bar).IsMatch(req)).
		Should().Equal(true).
		If(bar).Should().Equal(0).
		//
		If(gouldian.Get().Path("foo").QInt("bar", &bar).IsMatch(inv)).
		Should().Equal(false)
}

func TestHead(t *testing.T) {
	req := gouldian.Mock("/foo").
		With("Content-Type", "application/json")

	it.Ok(t).
		If(gouldian.Get().Path("foo").Head("Content-Type", "application/json").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").Head("Content-Type", "text/plain").IsMatch(req)).
		Should().Equal(false).
		//
		If(gouldian.Get().Path("foo").Head("Accept", "application/json").IsMatch(req)).
		Should().Equal(false)
}

func TestHasHead(t *testing.T) {
	req := gouldian.Mock("/foo").
		With("Content-Type", "application/json")

	it.Ok(t).
		If(gouldian.Get().Path("foo").HasHead("Content-Type").IsMatch(req)).
		Should().Equal(true).
		//
		If(gouldian.Get().Path("foo").HasHead("Accept").IsMatch(req)).
		Should().Equal(false)
}

func TestHeadString(t *testing.T) {
	req := gouldian.Mock("/foo").
		With("Content-Type", "application/json")
	content := ""

	it.Ok(t).
		If(gouldian.Get().Path("foo").HString("Content-Type", &content).IsMatch(req)).
		Should().Equal(true).
		If(content).Should().Equal("application/json").
		//
		If(gouldian.Get().Path("foo").HString("Accept", &content).IsMatch(req)).
		Should().Equal(true).
		If(content).Should().Equal("")
}

type foobar struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestJson(t *testing.T) {
	req := gouldian.Mock("/foo").
		WithJSON(foobar{"foo", 10})
	inv := gouldian.Mock("/foo").
		WithText("foobar")
	val := foobar{}

	it.Ok(t).
		If(gouldian.Get().Path("foo").JSON(&val).IsMatch(req)).
		Should().Equal(true).
		//
		If(val).Should().Equal(foobar{"foo", 10}).
		//
		If(gouldian.Get().Path("foo").JSON(&val).IsMatch(inv)).
		Should().Equal(false)
}

func TestThenSuccess(t *testing.T) {
	req := gouldian.Mock("/foo")
	handle := func() error { return gouldian.Ok().Text("bar") }

	it.Ok(t).
		If(gouldian.Get().Path("foo").FMap(handle)(req)).Should().
		Assert(
			func(be interface{}) bool {
				var rsp *gouldian.Output
				return errors.As(be.(error), &rsp)
			},
		)
}

func TestThenFailure(t *testing.T) {
	req := gouldian.Mock("/foo")
	handle := func() error { return gouldian.Unauthorized("") }

	it.Ok(t).
		If(gouldian.Get().Path("foo").FMap(handle)(req)).Should().
		Assert(
			func(be interface{}) bool {
				var rsp *gouldian.Output
				return !errors.As(be.(error), &rsp)
			},
		)
}

func TestAccessToken(t *testing.T) {
	req := gouldian.Mock("/foo").
		WithAuthorizer(map[string]interface{}{
			"sub":   "00000000-0000-0000-0000-000000000000",
			"scope": "a b",
		})
	val := gouldian.AccessToken{}

	it.Ok(t).
		If(gouldian.Get().Path("foo").AccessToken(&val).IsMatch(req)).
		Should().Equal(true).
		If(val.Sub).Should().Equal("00000000-0000-0000-0000-000000000000").
		If(val.Scope).Should().Equal("a b")
}
*/
