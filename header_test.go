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
	"errors"
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/it"
)

func TestHeaderIs(t *testing.T) {
	foo := µ.GET(µ.Header("X-Value").Is("some"))
	success := mock.Input(mock.Header("X-Value", "some"))
	failure := mock.Input(mock.Header("X-Value", "none"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestHeaderProduct(t *testing.T) {
	foo := µ.GET(
		µ.Join(
			µ.Header("X-Foo").Is("Bar"),
			µ.Header("X-Bar").Is("Foo"),
		),
	)

	t.Run("success", func(t *testing.T) {
		req := mock.Input(
			mock.Header("X-Foo", "Bar"),
			mock.Header("X-Bar", "Foo"),
		)

		it.Ok(t).
			If(foo(req)).Should().Equal(nil)
	})

	t.Run("incorrect", func(t *testing.T) {
		req := mock.Input(
			mock.Header("X-Foo", "Baz"),
			mock.Header("X-Bar", "Foo"),
		)

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("some", func(t *testing.T) {
		req := mock.Input(
			mock.Header("X-Foo", "Bar"),
		)

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("none", func(t *testing.T) {
		req := mock.Input()

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestHeaderCoProduct(t *testing.T) {
	foo := µ.GET(
		µ.Or(
			µ.Header("X-Foo").Is("Bar"),
			µ.Header("X-Bar").Is("Foo"),
		),
	)

	t.Run("success", func(t *testing.T) {
		req := mock.Input(
			mock.Header("X-Foo", "Bar"),
			mock.Header("X-Bar", "Foo"),
		)

		it.Ok(t).
			If(foo(req)).Should().Equal(nil)
	})

	t.Run("incorrect", func(t *testing.T) {
		req := mock.Input(
			mock.Header("X-Foo", "Baz"),
			mock.Header("X-Bar", "Foo"),
		)

		it.Ok(t).
			If(foo(req)).Should().Equal(nil)
	})

	t.Run("some", func(t *testing.T) {
		req := mock.Input(
			mock.Header("X-Foo", "Bar"),
		)

		it.Ok(t).
			If(foo(req)).Should().Equal(nil)
	})

	t.Run("none", func(t *testing.T) {
		req := mock.Input()

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestHeaderIsLowerCase(t *testing.T) {
	foo := µ.GET(µ.Header("X-Value").Is("bar"))
	success := mock.Input(mock.Header("x-value", "bar"))
	failure := mock.Input(mock.Header("x-value", "baz"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestHeaderAny(t *testing.T) {
	foo := µ.GET(µ.Header("X-Value").Any())
	bar := µ.GET(µ.Header("X-Value").Is("_"))

	success1 := mock.Input(mock.Header("X-Value", "bar"))
	success2 := mock.Input(mock.Header("X-Value", "baz"))
	failure := mock.Input()

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil).
		If(bar(success1)).Should().Equal(nil).
		If(bar(success2)).Should().Equal(nil).
		If(bar(failure)).ShouldNot().Equal(nil)
}

func TestHeaderString(t *testing.T) {
	type myT struct{ Val string }

	val := optics.Lenses1(myT{})
	foo := µ.GET(µ.Header("X-Value").To(val))

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal("bar")
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal("1")
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.Header("X-Foo", "bar"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestHeaderMaybeString(t *testing.T) {
	type myT struct{ Val string }

	val := optics.Lenses1(myT{})
	foo := µ.GET(µ.Header("X-Value").Maybe(val))

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal("bar")
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Foo", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal("")
	})
}

func TestHeaderInt(t *testing.T) {
	type myT struct{ Val int }

	val := optics.Lenses1(myT{})
	foo := µ.GET(µ.Header("X-Value").To(val))

	t.Run("string", func(t *testing.T) {
		req := mock.Input(mock.Header("X-Value", "bar"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "1024"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1024)
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.Header("X-Foo", "bar"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestHeaderMaybeInt(t *testing.T) {
	type myT struct{ Val int }

	val := optics.Lenses1(myT{})
	foo := µ.GET(µ.Header("X-Value").Maybe(val))

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal(0)
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "1024"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1024)
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Foo", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(req.Context().Get(&val)).Should().Equal(nil).
			If(val.Val).Should().Equal(0)
	})
}

func TestAuthorize(t *testing.T) {
	auth := func(scheme, token string) error {
		if token == "foo" {
			return nil
		}
		return errors.New("unauthorized")
	}
	foo := µ.GET(headers.Authorization.With(auth))

	t.Run("bearer", func(t *testing.T) {
		req := mock.Input(mock.Header("Authorization", "Bearer foo"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil)
	})

	t.Run("invalid", func(t *testing.T) {
		req := mock.Input(mock.Header("Authorization", "Digest_foo"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

	t.Run("nomatch", func(t *testing.T) {
		req := mock.Input(mock.Header("Authorization", "Bearer bar"))

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})
}

func TestContentJSON(t *testing.T) {
	for header, endpoint := range map[string]µ.Endpoint{
		"application/json":                  headers.ContentType.JSON,
		"application/x-www-form-urlencoded": headers.ContentType.Form,
		"text/plain":                        headers.ContentType.Text,
		"text/html":                         headers.ContentType.HTML,
	} {
		foo := µ.GET(endpoint)
		success := mock.Input(mock.Header("Content-Type", header))
		failure := mock.Input(mock.Header("Content-Type", "some/value"))

		it.Ok(t).
			If(foo(success)).Should().Equal(nil).
			If(foo(failure)).ShouldNot().Equal(nil)
	}
}
