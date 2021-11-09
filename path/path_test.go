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

package path_test

import (
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/gouldian/path"
	"github.com/fogfish/it"
)

func TestPathIs(t *testing.T) {
	foo := µ.GET(µ.Path(path.Is("foo")))
	success := mock.Input(mock.URL("/foo"))
	failure := mock.Input(mock.URL("/bar"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestPathAny(t *testing.T) {
	foo := µ.GET(µ.Path(path.Is("foo"), path.Any()))
	bar := µ.GET(µ.Path(path.Is("foo"), path.Is("*")))
	success1 := mock.Input(mock.URL("/foo/bar"))
	success2 := mock.Input(mock.URL("/foo/foo"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(bar(success1)).Should().Equal(nil).
		If(bar(success2)).Should().Equal(nil)
}

func TestPathString(t *testing.T) {
	type myT struct{ Val string }

	val := optics.Lenses1(myT{})
	foo := µ.GET(µ.Path(path.Is("foo"), path.String(val)))

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/foo/bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal("bar")
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/foo/1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal("1")
	})
}

func TestPathInt(t *testing.T) {
	type myT struct{ Val int }

	val := optics.Lenses1(myT{})
	foo := µ.GET(µ.Path(path.Is("foo"), path.Int(val)))

	t.Run("string", func(t *testing.T) {
		req := mock.Input(mock.URL("/foo/bar"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/foo/1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1)
	})
}

//
// TODO: recover
//
// func TestPathOr(t *testing.T) {
// 	pat := path.Is("foo").Or(path.Is("bar"))
// 	foo := µ.GET(µ.Path(pat))
// 	success1 := mock.Input(mock.URL("/foo"))
// 	success2 := mock.Input(mock.URL("/bar"))
// 	failure := mock.Input(mock.URL("/baz"))

// 	it.Ok(t).
// 		If(foo(success1)).Should().Equal(nil).
// 		If(foo(success2)).Should().Equal(nil).
// 		If(foo(failure)).ShouldNot().Equal(nil)
// }

// func TestPathThen(t *testing.T) {
// 	pat := path.Is("foo").Then(path.Is("bar"))
// 	foo := µ.GET(µ.Path(pat))
// 	success := mock.Input(mock.URL("/foo/bar"))
// 	failure1 := mock.Input(mock.URL("/foo"))
// 	failure2 := mock.Input(mock.URL("/foo/bar/baz"))

// 	it.Ok(t).
// 		If(foo(success)).Should().Equal(nil).
// 		If(foo(failure1)).ShouldNot().Equal(nil).
// 		If(foo(failure2)).ShouldNot().Equal(nil)
// }

func TestPathVariableLen(t *testing.T) {
	foo := µ.GET(µ.Path(path.Is("foo")))
	success := mock.Input(mock.URL("/foo"))
	failure := mock.Input(mock.URL("/foo/bar"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestPathSeq(t *testing.T) {
	type myT struct{ Val []string }

	val := optics.Lenses1(myT{})
	foo := µ.GET(µ.Path(path.Is("foo"), path.Seq(val)))

	t.Run("seq0", func(t *testing.T) {
		req := mock.Input(mock.URL("/foo"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("seq1", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/foo/a"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal([]string{"a"})
	})

	t.Run("seqN", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.URL("/foo/a/b/c"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal([]string{"a", "b", "c"})
	})
}
