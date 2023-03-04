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
	"time"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/mock"

	"github.com/fogfish/it"
	itt "github.com/fogfish/it/v2"
)

func TestHeaderProduct(t *testing.T) {
	foo := µ.Endpoints{
		µ.Header("X-Foo", "Bar"),
		µ.Header("X-Bar", "Foo"),
	}.Join

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
	foo := µ.Endpoints{
		µ.Header("X-Foo", "Bar"),
		µ.Header("X-Bar", "Foo"),
	}.Or

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

func TestHeadersMatch(t *testing.T) {
	spec := []struct {
		Endpoint µ.Endpoint
		Header   string
		Value    string
	}{
		{µ.Accept.Any, string(µ.Accept), "*/*"},
		{µ.Accept.Is("*/*"), string(µ.Accept), "*/*"},
		{µ.Accept.ApplicationJSON, string(µ.Accept), "application/json"},
		{µ.Accept.JSON, string(µ.Accept), "application/json"},
		{µ.Accept.Form, string(µ.Accept), "application/x-www-form-urlencoded"},
		{µ.Accept.TextPlain, string(µ.Accept), "text/plain"},
		{µ.Accept.Text, string(µ.Accept), "text/plain"},
		{µ.Accept.TextHTML, string(µ.Accept), "text/html"},
		{µ.Accept.HTML, string(µ.Accept), "text/html"},
		{µ.AcceptCharset.Any, string(µ.AcceptCharset), "utf8"},
		{µ.AcceptCharset.Is("utf8"), string(µ.AcceptCharset), "utf8"},
		{µ.AcceptEncoding.Is("gzip"), string(µ.AcceptEncoding), "gzip"},
		{µ.AcceptLanguage.Is("en"), string(µ.AcceptLanguage), "en"},
		{µ.CacheControl.Is("nocache"), string(µ.CacheControl), "nocache"},
		{µ.Connection.Any, string(µ.Connection), "keep-alive"},
		{µ.Connection.Is("keep-alive"), string(µ.Connection), "keep-alive"},
		{µ.Connection.KeepAlive, string(µ.Connection), "keep-alive"},
		{µ.Connection.Close, string(µ.Connection), "close"},
		{µ.ContentEncoding.Is("foo"), string(µ.ContentEncoding), "foo"},
		{µ.ContentLength.Is(1024), string(µ.ContentLength), "1024"},
		{µ.Header("Content-Length", 1024), string(µ.ContentLength), "1024"},
		{µ.ContentType.JSON, string(µ.ContentType), "application/json"},
		{µ.Cookie.Is("foo"), string(µ.Cookie), "foo"},
		{µ.Date.Is(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(µ.Date), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{µ.Header("Date", time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(µ.Date), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{µ.From.Is("foo"), string(µ.From), "foo"},
		{µ.Host.Is("foo"), string(µ.Host), "foo"},
		{µ.IfMatch.Is("foo"), string(µ.IfMatch), "foo"},
		{µ.IfModifiedSince.Is(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(µ.IfModifiedSince), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{µ.IfNoneMatch.Is("foo"), string(µ.IfNoneMatch), "foo"},
		{µ.IfRange.Is("foo"), string(µ.IfRange), "foo"},
		{µ.IfUnmodifiedSince.Is(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(µ.IfUnmodifiedSince), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{µ.Origin.Is("foo"), string(µ.Origin), "foo"},
		{µ.Range.Is("foo"), string(µ.Range), "foo"},
		{µ.Referer.Is("foo"), string(µ.Referer), "foo"},
		{µ.TransferEncoding.Any, string(µ.TransferEncoding), "chunked"},
		{µ.TransferEncoding.Is("chunked"), string(µ.TransferEncoding), "chunked"},
		{µ.TransferEncoding.Chunked, string(µ.TransferEncoding), "chunked"},
		{µ.TransferEncoding.Identity, string(µ.TransferEncoding), "identity"},
		{µ.UserAgent.Is("foo"), string(µ.UserAgent), "foo"},
		{µ.Upgrade.Is("foo"), string(µ.Upgrade), "foo"},
		{µ.HeaderAny("X-Value"), "X-Value", "bar"},
		{µ.Header("X-Value", "bar"), "X-Value", "bar"},
		{µ.Header("X-Value", "bar"), "x-value", "bar"},
	}

	for _, tt := range spec {
		req := mock.Input(mock.Header(tt.Header, tt.Value))
		err := tt.Endpoint(req)
		itt.Then(t).Should(
			itt.Nil(err),
		)
	}
}

func TestHeadersNoMatch(t *testing.T) {
	spec := []struct {
		Endpoint µ.Endpoint
		Header   string
		Value    string
	}{
		{µ.Accept.Any, string(µ.AcceptEncoding), "*/*"},
		{µ.ContentLength.Is(1024), string(µ.Accept), "10240"},
		{µ.ContentLength.Is(1024), string(µ.ContentLength), "text"},
		{µ.ContentLength.Is(1024), string(µ.ContentLength), "10240"},
		{µ.Date.Is(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(µ.Accept), "Wed, 01 Feb 2023 10:20:30 UTC"},
		{µ.Date.Is(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(µ.Date), "text"},
		{µ.Date.Is(time.Date(2023, 02, 01, 10, 20, 30, 0, time.UTC)), string(µ.Date), "Wed, 11 Feb 2023 10:20:30 UTC"},
		{µ.HeaderAny("X-Value"), "Y-Value", "bar"},
	}

	for _, tt := range spec {
		req := mock.Input(mock.Header(tt.Header, tt.Value))
		err := tt.Endpoint(req)
		itt.Then(t).ShouldNot(
			itt.Nil(err),
		)
	}
}

func TestHeaderString(t *testing.T) {
	type myT struct{ Val string }

	val := µ.Optics1[myT, string]()
	foo := µ.Header("X-Value", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal("bar")
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "1"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
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

	val := µ.Optics1[myT, string]()
	foo := µ.HeaderMaybe("X-Value", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal("bar")
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Foo", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal("")
	})
}

func TestHeaderInt(t *testing.T) {
	type myT struct{ Val int }

	val := µ.Optics1[myT, int]()
	foo := µ.Header("X-Value", val)

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
			If(µ.FromContext(req, &val)).Should().Equal(nil).
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

	val := µ.Optics1[myT, int]()
	foo := µ.HeaderMaybe("X-Value", val)

	t.Run("string", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(0)
	})

	t.Run("number", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Value", "1024"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(1024)
	})

	t.Run("nomatch", func(t *testing.T) {
		var val myT
		req := mock.Input(mock.Header("X-Foo", "bar"))

		it.Ok(t).
			If(foo(req)).Should().Equal(nil).
			If(µ.FromContext(req, &val)).Should().Equal(nil).
			If(val.Val).Should().Equal(0)
	})
}

func TestHeaderAuthorize(t *testing.T) {
	auth := func(scheme, token string) error {
		if token == "foo" {
			return nil
		}
		return errors.New("unauthorized")
	}
	foo := µ.Authorization(auth)

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

	t.Run("noheader", func(t *testing.T) {
		req := mock.Input()

		it.Ok(t).
			If(foo(req)).ShouldNot().Equal(nil)
	})

}

func TestHeaderContentJSON(t *testing.T) {
	for _, header := range []string{
		"application/json",
		"application/x-www-form-urlencoded",
		"text/plain",
		"text/html",
	} {
		foo := µ.Header(headers.ContentType, header)
		success := mock.Input(mock.Header("Content-Type", header))
		failure := mock.Input(mock.Header("Content-Type", "some/value"))

		it.Ok(t).
			If(foo(success)).Should().Equal(nil).
			If(foo(failure)).ShouldNot().Equal(nil)
	}
}
