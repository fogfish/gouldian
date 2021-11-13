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

	"github.com/fogfish/gouldian/optics"
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
type Header string

/*

Is matches a header to defined literal value.

  e := µ.GET( µ.Header("X-Foo").Is("Bar") )
  e(mock.Input(mock.Header("X-Foo", "Bar"))) == nil
  e(mock.Input(mock.Header("X-Foo", "Baz"))) != nil
*/
func (header Header) Is(val string) Endpoint {
	if val == Any {
		return header.Any
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

  e := µ.GET( µ.Header("X-Foo").Any )
  e(mock.Input(mock.Header("X-Foo", "Bar"))) == nil
  e(mock.Input(mock.Header("X-Foo", "Baz"))) == nil
  e(mock.Input()) != nil
*/
func (header Header) Any(req Input) error {
	_, exists := req.Headers().Get(string(header))
	if exists {
		return nil
	}
	return NoMatch{}
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
func (header Header) To(lens optics.Lens) Endpoint {
	return func(req Input) error {
		if opt, exists := req.Headers().Get(string(header)); exists {
			return req.Context().Put(lens, opt)
		}
		return NoMatch{}
	}
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
func (header Header) Maybe(lens optics.Lens) Endpoint {
	return func(req Input) error {
		if opt, exists := req.Headers().Get(string(header)); exists {
			req.Context().Put(lens, opt)
		}
		return nil
	}
}
