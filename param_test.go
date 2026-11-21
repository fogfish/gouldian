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
	"testing"

	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/mock"
	"github.com/fogfish/it/v2"
)

func TestParams(t *testing.T) {
	type foobar struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}
	type myT struct {
		Val foobar `content:"form"`
	}
	var val myT

	lens := µ.Optics1[myT, foobar]()
	foo := µ.Params(lens)

	req := mock.Input(mock.URL("/?foo=bar&bar=foo"))
	err := foo(req)

	it.Then(t).Should(
		it.Nil(err),
		it.Nil(µ.FromContext(req, &val)),
		it.Equal(val.Val.Foo, "bar"),
		it.Equal(val.Val.Bar, "foo"),
	)
}

func TestParamIs(t *testing.T) {
	foo := µ.Param("foo", "bar")
	success := mock.Input(mock.URL("/?foo=bar"))
	failure := mock.Input(mock.URL("/?bar=foo"))

	it.Then(t).
		Should(it.Equal(foo(success), nil)).
		ShouldNot(it.Equal(foo(failure), nil))
}

func TestParamAny(t *testing.T) {
	foo := µ.ParamAny("foo")
	success1 := mock.Input(mock.URL("/?foo"))
	success2 := mock.Input(mock.URL("/?foo=bar"))
	success3 := mock.Input(mock.URL("/?foo=baz"))
	failure := mock.Input(mock.URL("/?bar=foo"))

	it.Then(t).
		Should(it.Equal(foo(success1), nil)).
		Should(it.Equal(foo(success2), nil)).
		Should(it.Equal(foo(success3), nil)).
		ShouldNot(it.Equal(foo(failure), nil))
}

func TestParamString(t *testing.T) {
	type myT struct{ Val string }

	val := µ.Optics1[myT, string]()
	foo := µ.Param("foo", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, "bar"))
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, "1"))
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
	})
}

func TestParamMaybeString(t *testing.T) {
	type myT struct{ Val string }

	val := µ.Optics1[myT, string]()
	foo := µ.ParamMaybe("foo", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, "bar"))
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, ""))
	})
}

func TestParamInt(t *testing.T) {
	type myT struct{ Val int }

	val := µ.Optics1[myT, int]()
	foo := µ.Param("foo", val)

	t.Run("string", func(t *testing.T) {
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 1))
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
	})
}

func TestParamMaybeInt(t *testing.T) {
	type myT struct{ Val int }

	val := µ.Optics1[myT, int]()
	foo := µ.ParamMaybe("foo", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Then(t).
			Should(it.True(isFailure(foo(req)))).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 0))
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 1))
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 0))
	})
}

func TestParamFloat(t *testing.T) {
	type myT struct{ Val float64 }

	val := µ.Optics1[myT, float64]()
	foo := µ.Param("foo", val)

	t.Run("integer", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 1.0))
	})

	t.Run("double", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1.1"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 1.1))
	})

	t.Run("string", func(t *testing.T) {
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
	})
}

func TestParamMaybeFloat(t *testing.T) {
	type myT struct{ Val float64 }

	val := µ.Optics1[myT, float64]()
	foo := µ.ParamMaybe("foo", val)

	t.Run("double", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1.1"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 1.1))
	})

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Then(t).
			Should(it.True(isFailure(foo(req)))).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 0.0))
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, 0.0))
	})
}

type MyT struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func TestParamJSON(t *testing.T) {
	type MyS struct {
		A string `json:"a"`
		B int    `json:"b"`
	}

	type myT struct{ Val MyS }

	val := µ.Optics1[myT, MyS]()
	foo := µ.ParamJSON("foo", val)

	t.Run("json", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=%7B%22a%22%3A%22abc%22%2C%22b%22%3A10%7D"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, MyS{A: "abc", B: 10}))
	})

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Then(t).
			Should(it.True(isFailure(foo(req)))).
			Should(it.Equal(µ.FromContext(req, &val), nil))
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
	})

	t.Run("badformat", func(t *testing.T) {
		req := mock.Input(mock.URL("/?foo=%7"))

		it.Then(t).
			ShouldNot(it.Equal(foo(req), nil))
	})
}

func TestParamMaybeJSON(t *testing.T) {
	type MyS struct {
		A string `json:"a"`
		B int    `json:"b"`
	}

	type myT struct{ Val MyS }

	val := µ.Optics1[myT, MyS]()
	foo := µ.ParamMaybeJSON("foo", val)

	t.Run("json", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=%7B%22a%22%3A%22abc%22%2C%22b%22%3A10%7D"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, MyS{A: "abc", B: 10}))
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, MyS{}))

	})

	t.Run("badformat", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=%7"))

		it.Then(t).
			Should(it.Equal(foo(req), nil)).
			Should(it.Equal(µ.FromContext(req, &val), nil)).
			Should(it.Equal(val.Val, MyS{}))
	})
}

// // func TestParamOr(t *testing.T) {
// // 	foo := µ.GET(µ.Param(
// // 		param.Or(param.Any("foo"), param.Any("bar")),
// // 	))
// // 	success1 := mock.Input(mock.URL("/?foo"))
// // 	success2 := mock.Input(mock.URL("/?bar"))
// // 	failure := mock.Input(mock.URL("/?baz"))

// // 	it.Ok(t).
// // 		If(foo(success1)).Should().Equal(nil).
// // 		If(foo(success2)).Should().Equal(nil).
// // 		If(foo(failure)).ShouldNot().Equal(nil)
// // }
