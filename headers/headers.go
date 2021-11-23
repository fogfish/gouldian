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

package headers

import (
	"fmt"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/optics"
	"strings"
)

/*

List of supported HTTP header constants, use them instead of explicit definition
*/
const (
	// Common HTTP headers
	// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
	CacheControl     = µ.Header("Cache-Control")
	Connection       = µ.Header("Connection")
	ContentEncoding  = µ.Header("Content-Encoding")
	ContentLength    = µ.Header("Content-Length")
	ContentType      = µ.Header("Content-Type")
	Date             = µ.Header("Date")
	TransferEncoding = µ.Header("Transfer-Encoding")

	//
	// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Request_fields
	Accept             = µ.Header("Accept")
	AcceptCharset      = µ.Header("Accept-Charset")
	AcceptEncoding     = µ.Header("Accept-Encoding")
	AcceptLanguage     = µ.Header("Accept-Language")
	Authorization      = Authorize("Authorization")
	Cookie             = µ.Header("Cookie")
	Host               = µ.Header("Host")
	IfMatch            = µ.Header("If-Match")
	IfModifiedSince    = µ.Header("If-Modified-Since")
	IfNoneMatch        = µ.Header("If-None-Match")
	IfRange            = µ.Header("If-Range")
	IfUnmodifiedSince  = µ.Header("If-Unmodified-Since")
	Origin             = µ.Header("Origin")
	ProxyAuthorization = µ.Header("Proxy-Authorization")
	Range              = µ.Header("Range")
	UserAgent          = µ.Header("User-Agent")

	//
	// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Response_fields
	ContentLanguage = µ.Header("Content-Language")
	ETag            = µ.Header("ETag")
	Expires         = µ.Header("Expires")
	LastModified    = µ.Header("Last-Modified")
	Link            = µ.Header("Link")
	Location        = µ.Header("Location")
	Server          = µ.Header("Server")
	SetCookie       = µ.Header("Set-Cookie")
)

//
const (
	ApplicationJSON = "application/json"
	ApplicationForm = "application/x-www-form-urlencoded"
	TextPlain       = "text/plain"
	TextHTML        = "text/html"
)

/*

Authorize is "synonym" to Header type. It defines a few Endpoints that simplify
validation of credentials/tokens supplied within the request

  e := µ.GET( µ.Authorization.With(func(string, string) error { ... }) )
  e(mock.Input(mock.Header("Authorization", "Basic foo"))) == nil
  e(mock.Input(mock.Header("Authorization", "Basic bar"))) != nil

*/
type Authorize µ.Header

// To implements matcher for Content type (see Header.To)
func (h Authorize) To(lens optics.Lens) µ.Endpoint {
	return µ.Header(h).To(lens)
}

// Maybe implements matcher for Content type (see Header.Maybe)
func (h Authorize) Maybe(lens optics.Lens) µ.Endpoint {
	return µ.Header(h).Maybe(lens)
}

// With validates content of HTTP Authorization header
func (h Authorize) With(f func(string, string) error) µ.Endpoint {
	return func(ctx *µ.Context) error {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			return µ.Status.Unauthorized(
				µ.WithIssue(
					fmt.Errorf("Unauthorized %s", ctx.Request.URL.Path),
				),
			)
		}

		cred := strings.Split(auth, " ")
		if len(cred) != 2 {
			return µ.Status.Unauthorized(
				µ.WithIssue(
					fmt.Errorf("Unauthorized %v", ctx.Request.URL.Path),
				),
			)
		}

		if err := f(cred[0], cred[1]); err != nil {
			return µ.Status.Unauthorized(µ.WithIssue(err))
		}

		return nil
	}
}
