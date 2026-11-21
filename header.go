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
	"mime"
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
A literal is matched by exact (case-insensitive) value of the header.

	endpoint := µ.GET(
	  µ.Header("X-Foo", "Bar"),
	)

	endpoint(
	  mock.Input(
	    mock.Header("X-Foo", "Bar")
	  )
	) == nil

An absent header or differing literal is NoMatch. A header that is present but
cannot be decoded by the lens is Failure (400).
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
decode HTTP header into Golang type. The Endpoint matches if the header is not
defined. A header that is present but cannot be decoded is Failure.

	type myT struct{ Val string }

	x := µ.Optics1[myT, string]()
	e := µ.GET(µ.HeaderMaybe("X-Foo", x))
	e(mock.Input(mock.Header("X-Foo", "Bar"))) == nil
*/
func HeaderMaybe(header string, lens Lens) Endpoint {
	return func(ctx *Context) error {
		if opt := ctx.Request.Header.Get(header); opt != "" {
			return ctx.put(SourceHeader, header, lens, opt)
		}
		return nil
	}
}

func isHeaderExists(ctx *Context, header string) error {
	opt := ctx.Request.Header.Get(header)
	if opt == "" {
		return NoMatchHeader
	}
	return nil
}

func isHeaderEqString(ctx *Context, header string, value string) error {
	opt := ctx.Request.Header.Get(header)
	if opt == "" || !strings.EqualFold(strings.TrimSpace(opt), value) {
		return NoMatchHeader
	}
	return nil
}

// isHeaderHasToken matches the header with the list of tokens (comma
// separated), e.g. Connection: keep-alive, Upgrade
func isHeaderHasToken(ctx *Context, header string, token string) error {
	for _, opt := range ctx.Request.Header.Values(header) {
		for _, t := range strings.Split(opt, ",") {
			if strings.EqualFold(strings.TrimSpace(t), token) {
				return nil
			}
		}
	}
	return NoMatchHeader
}

func isHeaderEqInt(ctx *Context, header string, value int) error {
	opt := ctx.Request.Header.Get(header)
	if opt == "" {
		return NoMatchHeader
	}

	val, err := strconv.Atoi(opt)
	if err != nil {
		return newFailure(SourceHeader, header, err)
	}

	if val != value {
		return NoMatchHeader
	}

	return nil
}

func isHeaderEqTime(ctx *Context, header string, value time.Time) error {
	opt := ctx.Request.Header.Get(header)
	if opt == "" {
		return NoMatchHeader
	}

	val, err := time.Parse(time.RFC1123, opt)
	if err != nil {
		val, err = http.ParseTime(opt)
	}
	if err != nil {
		return newFailure(SourceHeader, header, err)
	}

	if !val.Equal(value) {
		return NoMatchHeader
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
// decode HTTP header into Golang type. The Endpoint is NoMatch if the header is
// absent, it is Failure if the value cannot be decoded to the target type.
func (h HeaderOf[T]) To(lens Lens) Endpoint {
	return func(ctx *Context) error {
		if opt := ctx.Request.Header.Get(string(h)); opt != "" {
			return ctx.put(SourceHeader, string(h), lens, opt)
		}
		return NoMatchHeader
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
		unauthorized := func(err error) error {
			out := NewOutput(http.StatusUnauthorized)
			out.SetIssue(err)
			return out
		}

		cred := strings.Fields(ctx.Request.Header.Get("Authorization"))
		if len(cred) != 2 {
			return unauthorized(fmt.Errorf("unauthorized %s", ctx.Request.URL.Path))
		}

		if err := f(cred[0], cred[1]); err != nil {
			return unauthorized(err)
		}

		return nil
	}
}

// Type of HTTP Header, Content-Type and Accept enumeration
//
//	const ContentType = HeaderEnumContent("Content-Type")
//	µ.ContentType.JSON
//
// The Content-Type matches by media type, parameters (charset) are ignored,
// the comparison is case-insensitive. A mismatch is NoMatchContentType, the
// server responds with 415 if no other endpoint matches.
//
// The Accept header is negotiated: comma separated list, q-values, \*/\* and
// type/\* are supported. A missing header accepts anything. A mismatch is
// NoMatchAccept, the server responds with 406 if no other endpoint matches.
type HeaderEnumContent string

func (h HeaderEnumContent) isAccept() bool {
	return strings.EqualFold(string(h), "Accept")
}

// Matches header to any value
func (h HeaderEnumContent) Any(ctx *Context) error {
	return isHeaderExists(ctx, string(h))
}

// Is matches the header to one of media types
func (h HeaderEnumContent) Is(values ...string) Endpoint {
	return func(ctx *Context) error { return h.match(ctx, values...) }
}

// Matches value of HTTP header
func (h HeaderEnumContent) To(lens Lens) Endpoint {
	return HeaderOf[string](h).To(lens)
}

func (h HeaderEnumContent) match(ctx *Context, types ...string) error {
	if h.isAccept() {
		return matchAccept(ctx.Request.Header.Values(string(h)), types)
	}

	opt := ctx.Request.Header.Get(string(h))
	mt, _, err := mime.ParseMediaType(opt)
	if err != nil {
		return NoMatchContentType
	}

	for _, t := range types {
		if mediaRange(t, mt) {
			return nil
		}
	}
	return NoMatchContentType
}

// mediaRange checks if media type `mt` is within the range `rng`
func mediaRange(rng, mt string) bool {
	rng = strings.ToLower(strings.TrimSpace(rng))
	mt = strings.ToLower(mt)
	switch {
	case rng == mt || rng == "*/*" || rng == "*":
		return true
	case strings.HasSuffix(rng, "/*"):
		return strings.HasPrefix(mt, rng[:len(rng)-1])
	default:
		return false
	}
}

func matchAccept(values []string, types []string) error {
	var seen bool
	for _, v := range values {
		for _, item := range strings.Split(v, ",") {
			if strings.TrimSpace(item) == "" {
				continue
			}
			seen = true

			rng, params, err := mime.ParseMediaType(item)
			if err != nil {
				continue
			}
			if q, ok := params["q"]; ok {
				if w, err := strconv.ParseFloat(q, 64); err == nil && w <= 0 {
					continue
				}
			}

			for _, t := range types {
				if mediaRange(rng, t) || mediaRange(t, rng) {
					return nil
				}
			}
		}
	}

	if !seen {
		return nil
	}
	return NoMatchAccept
}

func (h HeaderEnumContent) is(ctx *Context, mt string) error { return h.match(ctx, mt) }

// ApplicationJSON defines header `???: application/json`
func (h HeaderEnumContent) ApplicationJSON(ctx *Context) error { return h.is(ctx, "application/json") }

// JSON defines header `???: application/json`
func (h HeaderEnumContent) JSON(ctx *Context) error { return h.is(ctx, "application/json") }

// Form defined Header `???: application/x-www-form-urlencoded`
func (h HeaderEnumContent) Form(ctx *Context) error {
	return h.is(ctx, "application/x-www-form-urlencoded")
}

// TextPlain defined Header `???: text/plain`
func (h HeaderEnumContent) TextPlain(ctx *Context) error { return h.is(ctx, "text/plain") }

// Text defined Header `???: text/plain`
func (h HeaderEnumContent) Text(ctx *Context) error { return h.is(ctx, "text/plain") }

// TextHTML defined Header `???: text/html`
func (h HeaderEnumContent) TextHTML(ctx *Context) error { return h.is(ctx, "text/html") }

// HTML defined Header `???: text/html`
func (h HeaderEnumContent) HTML(ctx *Context) error { return h.is(ctx, "text/html") }

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
		return isHeaderHasToken(ctx, string(h), value)
	}
}

// Matches value of HTTP header
func (h HeaderEnumConnection) To(lens Lens) Endpoint {
	return HeaderOf[string](h).To(lens)
}

// KeepAlive defines header `???: keep-alive`
func (h HeaderEnumConnection) KeepAlive(ctx *Context) error {
	return isHeaderHasToken(ctx, string(h), "keep-alive")
}

// Close defines header `???: close`
func (h HeaderEnumConnection) Close(ctx *Context) error {
	return isHeaderHasToken(ctx, string(h), "close")
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
		return isHeaderHasToken(ctx, string(h), value)
	}
}

// Matches value of HTTP header
func (h HeaderEnumTransferEncoding) To(lens Lens) Endpoint {
	return HeaderOf[string](h).To(lens)
}

// Chunked defines header `Transfer-Encoding: chunked`
func (h HeaderEnumTransferEncoding) Chunked(ctx *Context) error {
	return isHeaderHasToken(ctx, string(h), "chunked")
}

// Identity defines header `Transfer-Encoding: identity`
func (h HeaderEnumTransferEncoding) Identity(ctx *Context) error {
	return isHeaderHasToken(ctx, string(h), "identity")
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
