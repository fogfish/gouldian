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
	"strings"
)

/*

Header type defines primitives to match Headers of HTTP requests.

  import "github.com/fogfish/gouldian/header"

  endpoint := µ.GET(
    µ.Header("X-Foo").Is("Bar"),
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
		panic("")
	}
}

func HeaderAny(hdr string) Endpoint {
	return header(hdr).Any
}

/*

Maybe matches header value to the request context. It uses lens abstraction to
decode HTTP header into Golang type. The Endpoint does not cause no-match
if header value cannot be decoded to the target type. See optics.Lens type for details.

  type myT struct{ Val string }

  x := optics.Lenses1(myT{})
  e := µ.GET(µ.Header("X-Foo").To(x))
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

type header string

/*

Is matches a header to defined literal value.

  e := µ.GET( µ.Header("X-Foo").Is("Bar") )
  e(mock.Input(mock.Header("X-Foo", "Bar"))) == nil
  e(mock.Input(mock.Header("X-Foo", "Baz"))) != nil
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

  e := µ.GET( µ.Header("X-Foo").Any )
  e(mock.Input(mock.Header("X-Foo", "Bar"))) == nil
  e(mock.Input(mock.Header("X-Foo", "Baz"))) == nil
  e(mock.Input()) != nil
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

  type myT struct{ Val string }

  x := optics.Lenses1(myT{})
  e := µ.GET(µ.Header("X-Foo").To(x))
  e(mock.Input(mock.Header("X-Foo", "Bar"))) == nil
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

Value outputs header value as the result of HTTP response
*/
// func HeaderValue(header string, value string) Result {
// 	return func(out *Output) error {
// 		out.Headers = append(out.Headers,
// 			struct {
// 				Header string
// 				Value  string
// 			}{header, value},
// 		)
// 		return nil
// 	}
// }
