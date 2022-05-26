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

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestParamIs(t *testing.T) {
	foo := µ.Param("foo", "bar")
	success := mock.Input(mock.URL("/?foo=bar"))
	failure := mock.Input(mock.URL("/?bar=foo"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestParamAny(t *testing.T) {
	foo := µ.ParamAny("foo")
	success1 := mock.Input(mock.URL("/?foo"))
	success2 := mock.Input(mock.URL("/?foo=bar"))
	success3 := mock.Input(mock.URL("/?foo=baz"))
	failure := mock.Input(mock.URL("/?bar=foo"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(foo(success3)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestParamString(t *testing.T) {
	type myT struct{ Val string }

	val := µ.Optics1[myT, string]()
	foo := µ.Param("foo", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal("bar")
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal("1")
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestParamMaybeString(t *testing.T) {
	type myT struct{ Val string }

	val := µ.Optics1[myT, string]()
	foo := µ.ParamMaybe("foo", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal("bar")
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal("")
	})
}

func TestParamInt(t *testing.T) {
	type myT struct{ Val int }

	val := µ.Optics1[myT, int]()
	foo := µ.Param("foo", val)

	t.Run("string", func(t *testing.T) {
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1)
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestParamMaybeInt(t *testing.T) {
	type myT struct{ Val int }

	val := µ.Optics1[myT, int]()
	foo := µ.ParamMaybe("foo", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(0)
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1)
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(0)
	})
}

func TestParamFloat(t *testing.T) {
	type myT struct{ Val float64 }

	val := µ.Optics1[myT, float64]()
	foo := µ.Param("foo", val)

	t.Run("integer", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1.0)
	})

	t.Run("double", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1.1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1.1)
	})

	t.Run("string", func(t *testing.T) {
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestParamMaybeFloat(t *testing.T) {
	type myT struct{ Val float64 }

	val := µ.Optics1[myT, float64]()
	foo := µ.ParamMaybe("foo", val)

	t.Run("double", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=1.1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1.1)
	})

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(0.0)
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(0.0)
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

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(MyS{A: "abc", B: 10})
	})

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).ShouldNot().Equal(nil)
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("badformat", func(t *testing.T) {
		req := mock.Input(mock.URL("/?foo=%7"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
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

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(MyS{A: "abc", B: 10})
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?bar=foo"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(MyS{})

	})

	t.Run("badformat", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/?foo=%7"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(MyS{})
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
