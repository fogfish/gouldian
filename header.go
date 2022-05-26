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
	"strings"
)

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
func Header[T Pattern](hdr string, val T) Endpoint {
	switch v := any(val).(type) {
	case string:
		return header(hdr).Is(v)
	case Lens:
		return header(hdr).To(v)
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
	return header(hdr).Any
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
type header string

/*

Is matches a header to defined literal value.
*/
func (header header) Is(val string) Endpoint {
	if val == Any {
		return header.Any
	}

	return func(ctx *Context) error {
		opt := ctx.Request.Header.Get(string(header))
		if opt != "" && strings.HasPrefix(opt, val) {
			return nil
		}
		return ErrNoMatch
	}
}

/*

Any is a wildcard matcher of header. It fails if header is not defined.
*/
func (header header) Any(ctx *Context) error {
	opt := ctx.Request.Header.Get(string(header))
	if opt != "" {
		return nil
	}
	return ErrNoMatch
}

/*

To matches header value to the request context. It uses lens abstraction to
decode HTTP header into Golang type. The Endpoint causes no-match if header
value cannot be decoded to the target type. See optics.Lens type for details.
*/
func (header header) To(lens Lens) Endpoint {
	return func(ctx *Context) error {
		if opt := ctx.Request.Header.Get(string(header)); opt != "" {
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
			return Status.Unauthorized(
				WithIssue(
					fmt.Errorf("Unauthorized %s", ctx.Request.URL.Path),
				),
			)
		}

		cred := strings.Split(auth, " ")
		if len(cred) != 2 {
			return Status.Unauthorized(
				WithIssue(
					fmt.Errorf("Unauthorized %v", ctx.Request.URL.Path),
				),
			)
		}

		if err := f(cred[0], cred[1]); err != nil {
			return Status.Unauthorized(WithIssue(err))
		}

		return nil
	}
}
