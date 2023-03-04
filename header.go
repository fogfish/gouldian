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

func isHeaderExists(ctx *Context, header string) error {
	opt := ctx.Request.Header.Get(string(header))
	if opt == "" {
		return ErrNoMatch
	}
	return nil
}

func isHeaderEqString(ctx *Context, header string, value string) error {
	opt := ctx.Request.Header.Get(string(header))
	if opt == "" || !strings.HasPrefix(opt, value) {
		return ErrNoMatch
	}
	return nil
}

func isHeaderEqInt(ctx *Context, header string, value int) error {
	opt := ctx.Request.Header.Get(string(header))
	if opt == "" {
		return ErrNoMatch
	}

	val, err := strconv.Atoi(opt)
	if err != nil {
		return ErrNoMatch
	}

	if val != value {
		return ErrNoMatch
	}

	return nil
}

func isHeaderEqTime(ctx *Context, header string, value time.Time) error {
	opt := ctx.Request.Header.Get(string(header))
	if opt == "" {
		return ErrNoMatch
	}

	val, err := time.Parse(time.RFC1123, opt)
	if err != nil {
		return ErrNoMatch
	}

	if !val.Equal(value) {
		return ErrNoMatch
	}

	return nil
}

// Internal type
type HeaderOf[T MatchableHeaderValues] string

// Any is a wildcard matcher of header. It fails if header is not defined.
func (h HeaderOf[T]) Any(ctx *Context) error {
	return isHeaderExists(ctx, string(h))
}

// Is matches a header to defined literal value.
func (h HeaderOf[T]) Is(value T) Endpoint {
	switch v := any(value).(type) {
	case string:
		return func(ctx *Context) error {
			return isHeaderEqString(ctx, string(h), v)
		}
	case int:
		return func(ctx *Context) error {
			return isHeaderEqInt(ctx, string(h), v)
		}
	case time.Time:
		return func(ctx *Context) error {
			return isHeaderEqTime(ctx, string(h), v)
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
			out.SetIssue(fmt.Errorf("unauthorized %s", ctx.Request.URL.Path))
			return out
		}

		cred := strings.Split(auth, " ")
		if len(cred) != 2 {
			out := NewOutput(http.StatusUnauthorized)
			out.SetIssue(fmt.Errorf("unauthorized %s", ctx.Request.URL.Path))
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

// Type of HTTP Header, Content-Type enumeration
//
//	const ContentType = HeaderEnumContent("Content-Type")
//	µ.ContentType.JSON
type HeaderEnumContent string

// Matches header to any value
func (h HeaderEnumContent) Any(ctx *Context) error {
	return isHeaderExists(ctx, string(h))
}

// Matches value of HTTP header
func (h HeaderEnumContent) Is(value string) Endpoint {
	return func(ctx *Context) error {
		return isHeaderEqString(ctx, string(h), value)
	}
}

// Matches value of HTTP header
func (h HeaderEnumContent) To(lens Lens) Endpoint {
	return HeaderOf[string](h).To(lens)
}

// ApplicationJSON defines header `???: application/json`
func (h HeaderEnumContent) ApplicationJSON(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "application/json")
}

// JSON defines header `???: application/json`
func (h HeaderEnumContent) JSON(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "application/json")
}

// Form defined Header `???: application/x-www-form-urlencoded`
func (h HeaderEnumContent) Form(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "application/x-www-form-urlencoded")
}

// TextPlain defined Header `???: text/plain`
func (h HeaderEnumContent) TextPlain(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "text/plain")
}

// Text defined Header `???: text/plain`
func (h HeaderEnumContent) Text(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "text/plain")
}

// TextHTML defined Header `???: text/html`
func (h HeaderEnumContent) TextHTML(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "text/html")
}

// HTML defined Header `???: text/html`
func (h HeaderEnumContent) HTML(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "text/html")
}

// Type of HTTP Header, Connection enumeration
//
//	const Connection = HeaderEnumConnection("Connection")
//	µ.Connection.KeepAlive
type HeaderEnumConnection string

// Matches header to any value
func (h HeaderEnumConnection) Any(ctx *Context) error {
	return isHeaderExists(ctx, string(h))
}

// Matches value of HTTP header
func (h HeaderEnumConnection) Is(value string) Endpoint {
	return func(ctx *Context) error {
		return isHeaderEqString(ctx, string(h), value)
	}
}

// Matches value of HTTP header
func (h HeaderEnumConnection) To(lens Lens) Endpoint {
	return HeaderOf[string](h).To(lens)
}

// KeepAlive defines header `???: keep-alive`
func (h HeaderEnumConnection) KeepAlive(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "keep-alive")
}

// Close defines header `???: close`
func (h HeaderEnumConnection) Close(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "close")
}

// Type of HTTP Header, Transfer-Encoding enumeration
//
//	const TransferEncoding = HeaderEnumTransferEncoding("Transfer-Encoding")
//	µ.TransferEncoding.Chunked
type HeaderEnumTransferEncoding string

// Matches header to any value
func (h HeaderEnumTransferEncoding) Any(ctx *Context) error {
	return isHeaderExists(ctx, string(h))
}

// Matches value of HTTP header
func (h HeaderEnumTransferEncoding) Is(value string) Endpoint {
	return func(ctx *Context) error {
		return isHeaderEqString(ctx, string(h), value)
	}
}

// Matches value of HTTP header
func (h HeaderEnumTransferEncoding) To(lens Lens) Endpoint {
	return HeaderOf[string](h).To(lens)
}

// Chunked defines header `Transfer-Encoding: chunked`
func (h HeaderEnumTransferEncoding) Chunked(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "chunked")
}

// Identity defines header `Transfer-Encoding: identity`
func (h HeaderEnumTransferEncoding) Identity(ctx *Context) error {
	return isHeaderEqString(ctx, string(h), "identity")
}

// List of supported HTTP header constants
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
const (
	Accept            = HeaderEnumContent("Accept")
	AcceptCharset     = HeaderOf[string]("Accept-Charset")
	AcceptEncoding    = HeaderOf[string]("Accept-Encoding")
	AcceptLanguage    = HeaderOf[string]("Accept-Language")
	CacheControl      = HeaderOf[string]("Cache-Control")
	Connection        = HeaderEnumConnection("Connection")
	ContentEncoding   = HeaderOf[string]("Content-Encoding")
	ContentLength     = HeaderOf[int]("Content-Length")
	ContentType       = HeaderEnumContent("Content-Type")
	Cookie            = HeaderOf[string]("Cookie")
	Date              = HeaderOf[time.Time]("Date")
	From              = HeaderOf[string]("From")
	Host              = HeaderOf[string]("Host")
	IfMatch           = HeaderOf[string]("If-Match")
	IfModifiedSince   = HeaderOf[time.Time]("If-Modified-Since")
	IfNoneMatch       = HeaderOf[string]("If-None-Match")
	IfRange           = HeaderOf[string]("If-Range")
	IfUnmodifiedSince = HeaderOf[time.Time]("If-Unmodified-Since")
	Origin            = HeaderOf[string]("Origin")
	Range             = HeaderOf[string]("Range")
	Referer           = HeaderOf[string]("Referer")
	TransferEncoding  = HeaderEnumTransferEncoding("Transfer-Encoding")
	UserAgent         = HeaderOf[string]("User-Agent")
	Upgrade           = HeaderOf[string]("Upgrade")
)
