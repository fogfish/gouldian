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
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/it"
)

func TestVerbAny(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.ANY(µ.Path()),
	)

	success1 := mock.Input(mock.Method("GET"))
	success2 := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success1)).Should().Equal(nil).
		If(endpoint(success2)).Should().Equal(nil)
}

func TestVerbDelete(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.DELETE(µ.Path()),
	)
	success := mock.Input(mock.Method("DELETE"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbGet(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.GET(µ.Path()),
	)
	success := mock.Input(mock.Method("GET"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbPatch(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.PATCH(µ.Path()),
	)
	success := mock.Input(mock.Method("PATCH"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbPost(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.POST(µ.Path()),
	)
	success := mock.Input(mock.Method("POST"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestVerbPut(t *testing.T) {
	endpoint := mock.Endpoint(
		µ.PUT(µ.Path()),
	)
	success := mock.Input(mock.Method("PUT"))
	failure := mock.Input(mock.Method("OTHER"))

	it.Ok(t).
		If(endpoint(success)).Should().Equal(nil).
		If(endpoint(failure)).ShouldNot().Equal(nil)
}

func TestPath(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(µ.Path("foo")),
	)
	bar := mock.Endpoint(
		µ.GET(µ.Path("bar")),
	)
	foobar := mock.Endpoint(µ.GET(µ.Path("foo", "bar")))

	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestPathRoot(t *testing.T) {
	root := mock.Endpoint(
		µ.GET(µ.Path()),
	)

	success := mock.Input(mock.URL("/"))
	failure := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(root(success)).Should().Equal(nil).
		If(root(failure)).ShouldNot().Equal(nil)
}

// //
// /*
// type MyType []string

// func (id *MyType) Pattern() µ.Endpoint {
// 	return func(req *µ.Input) error {
// 		var (
// 			a string
// 			b string
// 		)

// 		f := path.String(&a).Then(path.String(&b))
// 		switch err := f(segments).(type) {
// 		case µ.Match:
// 			*id = []string{a, b}
// 			return err
// 		default:
// 			return err
// 		}
// 	}
// }

// func TestPathTypeSafePattern(t *testing.T) {
// 	var id MyType

// 	foo := µ.GET(µ.Path(path.Is("foo"), id.Pattern()))
// 	success := mock.Input(mock.URL("/foo/a/b"))
// 	failure1 := mock.Input(mock.URL("/foo/a"))
// 	failure2 := mock.Input(mock.URL("/foo/a/b/c"))

// 	it.Ok(t).
// 		If(foo(success)).Should().Equal(nil).
// 		If(id[0]).Should().Equal("a").
// 		If(id[1]).Should().Equal("b").
// 		If(foo(failure1)).ShouldNot().Equal(nil).
// 		If(foo(failure2)).ShouldNot().Equal(nil)
// }
// */

func TestParam(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.Path(),
			µ.Param("foo").Is("bar"),
		),
	)
	bar := mock.Endpoint(
		µ.GET(
			µ.Path(),
			µ.Param("bar").Is("foo"),
		),
	)
	foobar := mock.Endpoint(
		µ.GET(
			µ.Path(),
			µ.Param("foo").Is("bar"),
			µ.Param("bar").Is("foo"),
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
			µ.Path(),
			µ.Header("foo").Is("bar"),
		),
	)

	bar := mock.Endpoint(
		µ.GET(
			µ.Path(),
			µ.Header("bar").Is("foo"),
		),
	)
	foobar := mock.Endpoint(
		µ.GET(
			µ.Path(),
			µ.Header("foo").Is("bar"),
			µ.Header("bar").Is("foo"),
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
	var lens = optics.ForProduct1(request{})

	var value request
	foo := mock.Endpoint(µ.GET(µ.Path(), µ.Body(lens)))
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
		If(success1.Get(&value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal(foobar{"foo1", 10}).
		//
		If(foo(success2)).Should().Equal(nil).
		If(success2.Get(&value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal(foobar{"foo2", 10}).
		//
		If(foo(failure1)).Should().Equal(nil).
		If(failure1.Get(&value)).ShouldNot().Equal(nil).
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
	var lens = optics.ForProduct1(request{})

	var value request
	foo := mock.Endpoint(µ.GET(µ.Path(), µ.Body(lens)))

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
		If(success1.Get(&value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal(foobar{"foo1", 10}).
		//
		If(foo(success2)).Should().Equal(nil).
		If(success2.Get(&value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal(foobar{"foo2", 10}).
		//
		If(foo(failure1)).Should().Equal(nil).
		If(failure1.Get(&value)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestText(t *testing.T) {
	type request struct {
		FooBar string
	}
	var lens = optics.ForProduct1(request{})

	var value request
	foo := mock.Endpoint(µ.GET(µ.Path(), µ.Body(lens)))
	success := mock.Input(mock.Text("foobar"))
	failure := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(success.Get(&value)).Should().Equal(nil).
		If(value.FooBar).Should().Equal("foobar").
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestContextFree(t *testing.T) {
	foo := mock.Endpoint(µ.GET(µ.Path("foo")))
	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil)

	req.Free()

	it.Ok(t).
		If(foo(req)).ShouldNot().Equal(nil)
}

func TestFMapSuccess(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.Path("foo"),
			µ.FMap(func(*µ.Context) error {
				return µ.Status.OK(µ.WithText("bar"))
			}),
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
	foo := mock.Endpoint(
		µ.GET(
			µ.Path("foo"),
			µ.FMap(func(*µ.Context) error {
				return µ.Status.
					OK(µ.WithText("bar"))
			}),
		),
	)
	bar := mock.Endpoint(
		µ.GET(
			µ.Path("bar"),
			func(*µ.Context) error {
				return µ.Status.
					OK(µ.WithText("foo"))
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

func TestFMapFailure(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.Path("foo"),
			µ.FMap(func(*µ.Context) error {
				return µ.Status.
					Unauthorized(µ.WithIssue(fmt.Errorf("")))
			}),
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
	lens := optics.ForProduct1(request{})

	endpoint := func() µ.Routable {
		return µ.GET(
			µ.Path(),
			µ.Body(lens),
			func(ctx *µ.Context) error {
				var req request
				if err := ctx.Get(&req); err != nil {
					return err
				}

				seq := []Pair{}
				for key, val := range req.Item.Seq {
					if val.Key == 0 {
						seq = append(seq, Pair{Key: key + 1, Val: val.Val})
					}
				}
				req.Item = Item{Seq: seq}
				return µ.Status.OK(µ.WithJSON(req.Item))
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

func TestAccessIs(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.Path(),
			µ.Access(µ.JWT.Sub).Is("sub"),
		),
	)
	success := mock.Input(mock.JWT(µ.JWT{"sub": "sub"}))
	failure1 := mock.Input(mock.JWT(µ.JWT{"sub": "foo"}))
	failure2 := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure1)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestAccessTo(t *testing.T) {
	type MyT struct{ Sub string }
	sub := optics.ForProduct1(MyT{})

	foo := mock.Endpoint(
		µ.GET(
			µ.Path(),
			µ.Access(µ.JWT.Sub).To(sub),
		),
	)

	t.Run("some", func(t *testing.T) {
		var val MyT
		req := mock.Input(mock.JWT(µ.JWT{"sub": "sub"}))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Get(&val)).Should().Equal(nil).
			If(val.Sub).Should().Equal("sub")
	})

	t.Run("none", func(t *testing.T) {
		req := mock.Input()

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestAccessMaybe(t *testing.T) {
	type MyT struct{ Sub string }
	sub := optics.ForProduct1(MyT{})

	foo := mock.Endpoint(
		µ.GET(
			µ.Path(),
			µ.Access(µ.JWT.Sub).Maybe(sub),
		),
	)

	t.Run("some", func(t *testing.T) {
		var val MyT
		req := mock.Input(mock.JWT(µ.JWT{"sub": "sub"}))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Get(&val)).Should().Equal(nil).
			If(val.Sub).Should().Equal("sub")
	})

	t.Run("empty", func(t *testing.T) {
		var val MyT
		req := mock.Input(mock.JWT(µ.JWT{}))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Get(&val)).Should().Equal(nil).
			If(val.Sub).Should().Equal("")
	})

	t.Run("none", func(t *testing.T) {
		req := mock.Input()

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

}
