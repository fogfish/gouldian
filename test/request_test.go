/*

  Copyright 2019 Dmitry Kolesnikov, All Rights Reserved

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package gouldian_test

import (
	"fmt"
	"net/http"
	"testing"

	µ "github.com/fogfish/gouldian"
	ø "github.com/fogfish/gouldian/emitter"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestVerbAny(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.ANY(µ.URI()),
	)

	success1 := mock.Input(mock.Method("GET"))
	success2 := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success1)).Should().Equal(nil).
		If(endpoint(success2)).Should().Equal(nil)
}

func TestVerbDelete(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.DELETE(µ.URI()),
	)
	success := mock.Input(mock.Method("DELETE"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbGet(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.GET(µ.URI()),
	)
	success := mock.Input(mock.Method("GET"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbPatch(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.PATCH(µ.URI()),
	)
	success := mock.Input(mock.Method("PATCH"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbPost(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.POST(µ.URI()),
	)
	success := mock.Input(mock.Method("POST"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbPut(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.PUT(µ.URI()),
	)
	success := mock.Input(mock.Method("PUT"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestPath(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(µ.URI(µ.Path("foo"))),
	)
	bar := mock.Endpoint(
		µ.GET(µ.URI(µ.Path("bar"))),
	)
	foobar := mock.Endpoint(µ.GET(µ.URI(µ.Path("foo"), µ.Path("bar"))))

	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestPathRoot(t *testing.T) {
	root := mock.Endpoint(
		µ.GET(µ.URI()),
	)

	success := mock.Input(mock.URL("/"))
	failure := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(root(success)).Should().Equal(nil).
		If(root(failure)).ShouldNot().Equal(nil)
}

func TestParam(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Param("foo", "bar"),
		),
	)
	bar := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Param("bar", "foo"),
		),
	)
	foobar := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Param("foo", "bar"),
			µ.Param("bar", "foo"),
		),
	)

	req := mock.Input(mock.URL("/?foo=bar"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestHeader(t *testing.T) {
	foo1 := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Header("foo", "bar"),
		),
	)

	bar := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Header("bar", "foo"),
		),
	)
	foobar := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Header("foo", "bar"),
			µ.Header("bar", "foo"),
		),
	)

	req := mock.Input(mock.Header("foo", "bar"))

	it.Ok(t).
		If(foo1(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestBodyJSON(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	type request struct {
		FooBar foobar
	}
	var lens = µ.Optics1[request, foobar]()

	var value request
	foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
	success1 := mock.Input(
		mock.JSON(foobar{"foo1", 10}),
	)
	success2 := mock.Input(
		mock.Header("content-type", "application/json"),
		mock.Text("{\"foo\":\"foo2\",\"bar\":10}"),
	)
	failure1 := mock.Input(
		mock.Header("Content-Type", "application/json"),
		mock.Text("foobar"),
	)
	failure2 := mock.Input()

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(µ.FromContext(success1, &value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal(foobar{"foo1", 10}).
		//
		If(foo(success2)).Should().Equal(nil).
		If(µ.FromContext(success2, &value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal(foobar{"foo2", 10}).
		//
		If(foo(failure1)).Should().Equal(nil).
		If(µ.FromContext(failure1, &value)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestBodyForm(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	type request struct {
		FooBar foobar `content:"form"`
	}
	var lens = µ.Optics1[request, foobar]()

	var value request
	foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))

	success1 := mock.Input(
		mock.Header("Content-Type", "application/x-www-form-urlencoded"),
		mock.Text("foo=foo1&bar=10"),
	)
	success2 := mock.Input(
		mock.Header("content-type", "application/x-www-form-urlencoded"),
		mock.Text("foo=foo2&bar=10"),
	)
	failure1 := mock.Input(
		mock.Header("Content-Type", "application/x-www-form-urlencoded"),
		mock.Text("foobar"),
	)
	failure2 := mock.Input()

	it.Ok(t).
		//
		If(foo(success1)).Should().Equal(nil).
		If(µ.FromContext(success1, &value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal(foobar{"foo1", 10}).
		//
		If(foo(success2)).Should().Equal(nil).
		If(µ.FromContext(success2, &value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal(foobar{"foo2", 10}).
		//
		If(foo(failure1)).Should().Equal(nil).
		If(µ.FromContext(failure1, &value)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestText(t *testing.T) {
	type request struct {
		FooBar string
	}
	var lens = µ.Optics1[request, string]()

	var value request
	foo := mock.Endpoint(µ.GET(µ.URI(), µ.Body(lens)))
	success := mock.Input(mock.Text("foobar"))
	failure := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(µ.FromContext(success, &value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal("foobar").
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestContextFree(t *testing.T) {
	foo := mock.Endpoint(µ.GET(µ.URI(µ.Path("foo"))))
	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil)

	req.Free()

	it.Ok(t).
		If(foo(req)).ShouldNot().Equal(nil)
}

func TestHandlerSuccess(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo")),
			func(*µ.Context) error {
				return ø.Status.OK(ø.Send("bar"))
			},
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

func TestFMapSuccess(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo"), µ.Path(a)),
			µ.FMap(func(ctx *µ.Context, t *T) error {
				return ø.Status.OK(ø.Send(t.A))
			}),
		),
	)
	req := mock.Input(mock.URL("/foo/bar"))

	it.Ok(t).
		If(foo(req)).Should().Assert(
		func(be interface{}) bool {
			return be.(error).Error() == "bar"
		},
	)
}

func TestMapSuccess(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo"), µ.Path(a)),
			µ.Map(func(ctx *µ.Context, t *T) (*T, error) {
				return t, nil
			}),
		),
	)
	req := mock.Input(mock.URL("/foo/bar"))

	it.Ok(t).
		If(foo(req)).Should().Assert(
		func(be interface{}) bool {
			return be.(error).Error() == "{\"A\":\"bar\"}"
		},
	)
}

func TestHandler2Success(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo")),
			func(*µ.Context) error {
				return ø.Status.OK(ø.Send("bar"))
			},
		),
	)
	bar := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("bar")),
			func(*µ.Context) error {
				return ø.Status.OK(ø.Send("foo"))
			},
		),
	)
	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(µ.Endpoints{foo, bar}.Or(req)).Should().Assert(
		func(be interface{}) bool {
			return be.(error).Error() == "bar"
		},
	)
}

func TestHandlerFailure(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo")),
			func(*µ.Context) error {
				return ø.Status.Unauthorized(ø.Error(fmt.Errorf("")))
			},
		),
	)
	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Assert(
		func(be interface{}) bool {
			switch v := be.(type) {
			case *µ.Output:
				return v.Status == http.StatusUnauthorized
			default:
				return false
			}
		},
	)
}

func TestFMapFailure(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo"), µ.Path(a)),
			µ.FMap(func(*µ.Context, *T) error {
				return ø.Status.Unauthorized(ø.Error(fmt.Errorf("")))
			}),
		),
	)
	req := mock.Input(mock.URL("/foo/bar"))

	it.Ok(t).
		If(foo(req)).Should().Assert(
		func(be interface{}) bool {
			switch v := be.(type) {
			case *µ.Output:
				return v.Status == http.StatusUnauthorized
			default:
				return false
			}
		},
	)
}

func TestMapFailure(t *testing.T) {
	type T struct{ A string }
	a := µ.Optics1[T, string]()

	foo := mock.Endpoint(
		µ.GET(
			µ.URI(µ.Path("foo"), µ.Path(a)),
			µ.Map(func(*µ.Context, *T) (*T, error) {
				return nil, ø.Status.Unauthorized(ø.Error(fmt.Errorf("")))
			}),
		),
	)
	req := mock.Input(mock.URL("/foo/bar"))

	it.Ok(t).
		If(foo(req)).Should().Assert(
		func(be interface{}) bool {
			switch v := be.(type) {
			case *µ.Output:
				return v.Status == http.StatusUnauthorized
			default:
				return false
			}
		},
	)
}

func TestBodyLeak(t *testing.T) {
	type Pair struct {
		Key int    `json:"key,omitempty"`
		Val string `json:"val,omitempty"`
	}
	type Item struct {
		Seq []Pair `json:"seq,omitempty"`
	}
	type request struct {
		Item Item
	}
	lens := µ.Optics1[request, Item]()

	endpoint := func() µ.Routable {
		return µ.GET(
			µ.URI(),
			µ.Body(lens),
			func(ctx *µ.Context) error {
				var req request
				if err := µ.FromContext(ctx, &req); err != nil {
					return err
				}

				seq := []Pair{}
				for key, val := range req.Item.Seq {
					if val.Key == 0 {
						seq = append(seq, Pair{Key: key + 1, Val: val.Val})
					}
				}
				req.Item = Item{Seq: seq}
				return ø.Status.OK(ø.Send(req.Item))
			},
		)
	}

	foo := mock.Endpoint(endpoint())
	for val, expect := range map[string]string{
		"{\"seq\":[{\"val\":\"a\"},{\"val\":\"b\"}]}":                 "{\"seq\":[{\"key\":1,\"val\":\"a\"},{\"key\":2,\"val\":\"b\"}]}",
		"{\"seq\":[{\"val\":\"c\"}]}":                                 "{\"seq\":[{\"key\":1,\"val\":\"c\"}]}",
		"{\"seq\":[{\"val\":\"d\"},{\"val\":\"e\"},{\"val\":\"f\"}]}": "{\"seq\":[{\"key\":1,\"val\":\"d\"},{\"key\":2,\"val\":\"e\"},{\"key\":3,\"val\":\"f\"}]}",
	} {
		req := mock.Input(
			mock.Method("GET"),
			mock.Header("Content-Type", "application/json"),
			mock.Text(val),
		)
		out := foo(req)
		it.Ok(t).If(out.Error()).Should().Equal(expect)
	}
}
