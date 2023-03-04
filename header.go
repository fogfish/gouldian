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

package gouldian

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReadableHeaderValues interface {
	int | string | time.Time
}

type MatchableHeaderValues interface {
	ReadableHeaderValues | Lens
}

/*
Header combinator defines primitives to match Headers of HTTP requests.

	endpoint := µ.GET(
	  µ.Header("X-Foo", "Bar"),
	)

	endpoint(
	  mock.Input(
	    mock.Header("X-Foo", "Bar")
	  )
	) == nil
*/
func Header[T MatchableHeaderValues](hdr string, val T) Endpoint {
	switch v := any(val).(type) {
	case string:
		return HeaderOf[string](hdr).Is(v)
	case int:
		return HeaderOf[int](hdr).Is(v)
	case time.Time:
		return HeaderOf[time.Time](hdr).Is(v)
	case Lens:
		return HeaderOf[Lens](hdr).To(v)
	default:
		panic("type system failure")
	}
}

/*
HeaderAny is a wildcard matcher of header. It fails if header is not defined.

	e := µ.GET( µ.HeaderAny("X-Foo") )
	e(mock.Input(mock.Header("X-Foo", "Bar"))) == nil
	e(mock.Input(mock.Header("X-Foo", "Baz"))) == nil
	e(mock.Input()) != nil
*/
func HeaderAny(hdr string) Endpoint {
	return HeaderOf[string](hdr).Any
}

/*
HeaderMaybe matches header value to the request context. It uses lens abstraction to
decode HTTP header into Golang type. The Endpoint does not cause no-match
if header value cannot be decoded to the target type. See optics.Lens type for details.

	type myT struct{ Val string }

	x := µ.Optics1[myT, string]()
	e := µ.GET(µ.HeaderMaybe("X-Foo", x))
	e(mock.Input(mock.Header("X-Foo", "Bar"))) == nil
*/
func HeaderMaybe(header string, lens Lens) Endpoint {
	return func(ctx *Context) error {
		if opt := ctx.Request.Header.Get(string(header)); opt != "" {
			ctx.Put(lens, opt)
		}
		return nil
	}
}

// Internal type
type HeaderOf[T MatchableHeaderValues] string

// Any is a wildcard matcher of header. It fails if header is not defined.
func (h HeaderOf[T]) Any(ctx *Context) error {
	opt := ctx.Request.Header.Get(string(h))
	if opt != "" {
		return nil
	}
	return ErrNoMatch
}

// Is matches a header to defined literal value.
func (h HeaderOf[T]) Is(value T) Endpoint {
	switch v := any(value).(type) {
	case string:
		return func(ctx *Context) error {
			opt := ctx.Request.Header.Get(string(h))
			if opt != "" && strings.HasPrefix(opt, v) {
				return nil
			}
			return ErrNoMatch
		}
	case int:
		return func(ctx *Context) error {
			opt := ctx.Request.Header.Get(string(h))
			val, err := strconv.Atoi(opt)
			if err == nil && val == v {
				return nil
			}
			return ErrNoMatch
		}
	case time.Time:
		return func(ctx *Context) error {
			t := v.UTC().Format(time.RFC1123)
			opt := ctx.Request.Header.Get(string(h))
			if opt != "" && t == opt {
				return nil
			}
			return ErrNoMatch
		}
	default:
		panic("invalid type")
	}
}

// To matches header value to the request context. It uses lens abstraction to
// decode HTTP header into Golang type. The Endpoint causes no-match if header
// value cannot be decoded to the target type. See optics.Lens type for details.
func (h HeaderOf[T]) To(lens Lens) Endpoint {
	return func(ctx *Context) error {
		if opt := ctx.Request.Header.Get(string(h)); opt != "" {
			return ctx.Put(lens, opt)
		}
		return ErrNoMatch
	}
}

/*
Authorization defines Endpoints that simplify validation of credentials/tokens
supplied within the request

	e := µ.GET( µ.Authorization(func(string, string) error { ... }) )
	e(mock.Input(mock.Header("Authorization", "Basic foo"))) == nil
	e(mock.Input(mock.Header("Authorization", "Basic bar"))) != nil
*/
func Authorization(f func(string, string) error) Endpoint {
	return func(ctx *Context) error {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			out := NewOutput(http.StatusUnauthorized)
			out.SetIssue(fmt.Errorf("Unauthorized %s", ctx.Request.URL.Path))
			return out
		}

		cred := strings.Split(auth, " ")
		if len(cred) != 2 {
			out := NewOutput(http.StatusUnauthorized)
			out.SetIssue(fmt.Errorf("Unauthorized %s", ctx.Request.URL.Path))
			return out
		}

		if err := f(cred[0], cred[1]); err != nil {
			out := NewOutput(http.StatusUnauthorized)
			out.SetIssue(err)
			return out
		}

		return nil
	}
}

// List of supported HTTP header constants
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
const (
	// Accept            = HeaderEnumContent("Accept")
	// AcceptCharset     = HeaderOf[string]("Accept-Charset")
	// AcceptEncoding    = HeaderOf[string]("Accept-Encoding")
	// AcceptLanguage    = HeaderOf[string]("Accept-Language")
	// Authorization     = HeaderOf[string]("Authorization")
	// CacheControl      = HeaderOf[string]("Cache-Control")
	// Connection        = HeaderEnumConnection("Connection")
	// ContentEncoding   = HeaderOf[string]("Content-Encoding")
	// ContentLength     = HeaderEnumContentLength("Content-Length")
	// ContentType       = HeaderEnumContent("Content-Type")
	// Cookie            = HeaderOf[string]("Cookie")
	// Date              = HeaderOf[time.Time]("Date")
	// From              = HeaderOf[string]("From")
	// Host              = HeaderOf[string]("Host")
	// IfMatch           = HeaderOf[string]("If-Match")
	// IfModifiedSince   = HeaderOf[time.Time]("If-Modified-Since")
	// IfNoneMatch       = HeaderOf[string]("If-None-Match")
	// IfRange           = HeaderOf[string]("If-Range")
	// IfUnmodifiedSince = HeaderOf[time.Time]("If-Unmodified-Since")
	// Origin            = HeaderOf[string]("Origin")
	// Range             = HeaderOf[string]("Range")
	// Referer           = HeaderOf[string]("Referer")
	// TransferEncoding  = HeaderEnumTransferEncoding("Transfer-Encoding")
	// UserAgent         = HeaderOf[string]("User-Agent")
	Upgrade = HeaderOf[string]("Upgrade")
)
