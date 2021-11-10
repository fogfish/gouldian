//
//   Copyright 2019 Dmitry Kolesnikov, All Rights Reserved
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

/*

Package header defines primitives to match Headers of HTTP requests.

	import "github.com/fogfish/gouldian/header"

	endpoint := µ.GET(
		µ.Header(
			header.Is("Content-Type", "application/json"),
			...
		)
	)
	Json := mock.Header("Content-Type", "application/json")
	endpoint(mock.Input(Json)) == nil

*/
package gouldian

import (
	"strconv"
	"strings"

	"github.com/fogfish/gouldian/optics"
)

/*

Header matches presence of header in the request or match its entire content.

	header.ContentType.JSON,
	header.ContentEncoding.Is(...),

*/
type Header string

// TODO: Path, Param

/*

List of supported HTTP header constants
https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
*/
const (
	Accept             = Header("Accept")
	AcceptCharset      = Header("Accept-Charset")
	AcceptEncoding     = Header("Accept-Encoding")
	AcceptLanguage     = Header("Accept-Language")
	Authorization      = Header("Authorization")
	CacheControl       = Header("Cache-Control")
	Connection         = Header("Connection")
	ContentEncoding    = Header("Content-Encoding")
	ContentLength      = Header("Content-Length")
	ContentType        = Content("Content-Type")
	Cookie             = Header("Cookie")
	Date               = Header("Date")
	Host               = Header("Host")
	IfMatch            = Header("If-Match")
	IfModifiedSince    = Header("If-Modified-Since")
	IfNoneMatch        = Header("If-None-Match")
	IfRange            = Header("If-Range")
	IfUnmodifiedSince  = Header("If-Unmodified-Since")
	Origin             = Header("Origin")
	ProxyAuthorization = Header("Proxy-Authorization")
	Range              = Header("Range")
	TransferEncoding   = Header("Transfer-Encoding")
	UserAgent          = Header("User-Agent")
)

/*

Is matches a header to defined literal value
  e := µ.GET( µ.Header(header.Is("Content-Type", "application/json")) )
  e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil
  e(mock.Input(mock.Header("Content-Type", "text/plain"))) != nil
*/
func (header Header) Is(val string) Endpoint {
	if val == "*" {
		return header.Any()
	}

	return func(req Input) error {
		opt, exists := req.Headers().Get(string(header))
		if exists && strings.HasPrefix(opt, val) {
			return nil
		}
		return NoMatch{}
	}
}

/*

Any is a wildcard matcher of header. It fails if header is not defined.
  e := µ.GET( µ.Header(header.Any("Content-Type")) )
  e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil
  e(mock.Input(mock.Header("Content-Type", "text/plain"))) == nil
  e(mock.Input()) != nil
*/
func (header Header) Any() Endpoint {
	return func(req Input) error {
		_, exists := req.Headers().Get(string(header))
		if exists {
			return nil
		}
		return NoMatch{}
	}
}

/*

String matches a header value to closed variable of string type.
It fails if header is not defined.
  var value string
  e := µ.GET( µ.Header(header.String("Content-Type", &value)) )
  e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil && value == "application/json"
  e(mock.Input()) != nil
*/
func (header Header) String(lens optics.Lens) Endpoint {
	return func(req Input) error {
		val, exists := req.Headers().Get(string(header))
		if !exists {
			return NoMatch{}
		}

		req.Context().Put(lens, val)
		return nil
	}
}

/*

MaybeString matches a header value to closed variable of string type.
It does not fail if header is not defined.
  var value string
  e := µ.GET( µ.Header(header.String("foo", &value)) )
  e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil && value == "application/json"
  e(mock.Input()) == nil
*/
func (header Header) MaybeString(lens optics.Lens) Endpoint {
	return func(req Input) error {
		val, exists := req.Headers().Get(string(header))
		if !exists {
			return nil
		}

		req.Context().Put(lens, val)
		return nil
	}
}

/*

Int matches a header value to closed variable of int type.
It fails if header is not defined.
  var value int
  e := µ.GET( µ.Header(header.Int("Content-Length", &value)) )
  e(mock.Input(mock.Header("Content-Length", "1024"))) == nil && value == 1024
  e(mock.Input()) != nil
*/
func (header Header) Int(lens optics.Lens) Endpoint {
	return func(req Input) error {
		val, exists := req.Headers().Get(string(header))
		if !exists {
			return NoMatch{}
		}

		ivl, err := strconv.Atoi(val)
		if err != nil {
			return NoMatch{}
		}

		req.Context().Put(lens, ivl)
		return nil
	}
}

/*

MaybeInt matches a header value to closed variable of int type.
It does not fail if header is not defined.
  var value int
  e := µ.GET( µ.Header(header.MaybeInt("Content-Length", &value)) )
  e(mock.Input(mock.Header("Content-Length", "1024"))) == nil && value == 1024
  e(mock.Input()) == nil
*/
func (header Header) MaybeInt(lens optics.Lens) Endpoint {
	return func(req Input) error {
		val, exists := req.Headers().Get(string(header))
		if !exists {
			return nil
		}

		ivl, err := strconv.Atoi(val)
		if err != nil {
			return nil
		}

		req.Context().Put(lens, ivl)
		return nil

	}
}

// Content defines headers for content negotiation
type Content Header

// JSON is a syntax sugar to header.ContentType.Is("application/json")
func (h Content) JSON() Endpoint {
	return Header(h).Is("application/json")
}

// Form is a syntax sugar to header.ContentType.Is("application/x-www-form-urlencoded")
func (h Content) Form() Endpoint {
	return Header(h).Is("application/x-www-form-urlencoded")
}

// Text is a syntax sugar to header.ContentType.Is("application/x-www-form-urlencoded")
func (h Content) Text() Endpoint {
	return Header(h).Is("text/plain")
}

// HTML matches Header `???: text/html`
func (h Content) HTML() Endpoint {
	return Header(h).Is("text/html")
}

// Is matches value of HTTP header, Use wildcard string ("*") to match any header value
func (h Content) Is(value string) Endpoint {
	return Header(h).Is(value)
}

// Any matches a header value `???: *`
func (h Content) Any() Endpoint {
	return Header(h).Any()
}

// String matches a header value to closed variable of string type.
func (h Content) String(lens optics.Lens) Endpoint {
	return Header(h).String(lens)
}

// MaybeString matches a header value to closed variable of string type.
func (h Content) MaybeString(lens optics.Lens) Endpoint {
	return Header(h).MaybeString(lens)
}

// Int matches a header value to closed variable of string type.
func (h Content) Int(lens optics.Lens) Endpoint {
	return Header(h).Int(lens)
}

// MaybeInt matches a header value to closed variable of string type.
func (h Content) MaybeInt(lens optics.Lens) Endpoint {
	return Header(h).MaybeInt(lens)
}

// Or is a co-product of header match arrows
//   e := µ.GET(
//     µ.Header(
//       header.Or(
//         header.Is("Content-Type", "application/json"),
//         header.Is("Content-Type", "text/plain"),
//       )
//     )
//   )
//   e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil
//   e(mock.Input(mock.Header("Content-Type", "text/plain"))) == nil
//   e(mock.Input(mock.Header("Content-Type", "text/html"))) != nil
/*
func Or(arrows ...µ.ArrowHeader) µ.ArrowHeader {
	return func(ctx µ.Context, headers µ.Headers) error {
		for _, f := range arrows {
			if err := f(ctx, headers); !errors.Is(err, µ.NoMatch{}) {
				return err
			}
		}
		return µ.NoMatch{}
	}
}
*/

// TODO
// Authorize validates content of HTTP Authorization header
/*
func Authorize(authType string, f func(string) error) µ.Endpoint {
	return func(req *µ.Input) error {
		auth, exists := req.Header("Authorization")
		if !exists {
			return µ.Unauthorized(fmt.Errorf("Unauthorized %v", req.APIGatewayProxyRequest.Path))
		}

		cred := strings.Split(auth, " ")
		if len(cred) != 2 {
			return µ.Unauthorized(fmt.Errorf("Unauthorized %v", req.APIGatewayProxyRequest.Path))
		}

		if strings.ToLower(cred[0]) != strings.ToLower(authType) {
			return µ.Unauthorized(fmt.Errorf("Unauthorized %v", req.APIGatewayProxyRequest.Path))
		}

		if err := f(cred[1]); err != nil {
			return µ.Unauthorized(err)
		}

		return nil
	}
}
*/
