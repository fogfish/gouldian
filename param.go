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
	"net/url"

	"github.com/fogfish/gouldian/optics"
)

/*

Param type defines primitives to match query param in the HTTP requests.

  import "github.com/fogfish/gouldian/header"

  endpoint := µ.GET(
    µ.Param("foo").Is("bar"),
  )

  endpoint(
    mock.Input(
			mock.URL("/?foo=bar")
    )
  ) == nil

*/
type Param string

/*

Is matches a param key to defined literal value

  e := µ.GET( µ.Param("foo").Is("bar") )
  e(mock.Input(mock.URL("/?foo=bar"))) == nil
  e(mock.Input(mock.URL("/?bar=foo"))) != nil
*/
func (key Param) Is(val string) Endpoint {
	if val == Any {
		return key.Any
	}

	return func(req Input) error {
		opt, exists := req.Params().Get(string(key))
		if exists && opt == val {
			return nil
		}
		return NoMatch{}
	}
}

/*

Any is a wildcard matcher of param key. It fails if key is not defined.

  e := µ.GET( µ.Param("foo").Any )
  e(mock.Input(mock.URL("/?foo"))) == nil
  e(mock.Input(mock.URL("/?foo=bar"))) == nil
  e(mock.Input(mock.URL("/?foo=baz"))) == nil
  e(mock.Input()) != nil
*/
func (key Param) Any(req Input) error {
	_, exists := req.Params().Get(string(key))
	if exists {
		return nil
	}
	return NoMatch{}
}

/*

To matches param value to the request context. It uses lens abstraction to
decode value into Golang type. The Endpoint causes no-match if param
value cannot be decoded to the target type. See optics.Lens type for details.

  type myT struct{ Val string }

  x := optics.Lenses1(myT{})
  e := µ.GET( µ.Param("foo").To(x) )
  e(mock.Input(mock.URL("/?foo=bar"))) == nil
*/
func (key Param) To(lens optics.Lens) Endpoint {
	return func(req Input) error {
		opt, exists := req.Params().Get(string(key))
		if exists {
			return req.Context().Put(lens, opt)
		}
		return NoMatch{}
	}
}

/*

Maybe matches param value to the request context. It uses lens abstraction to
decode value into Golang type. The Endpoint does not cause no-match
if header value cannot be decoded to the target type. See optics.Lens type for details.

  type myT struct{ Val string }

  x := optics.Lenses1(myT{})
  e := µ.GET( µ.Param("foo").To(x) )
  e(mock.Input(mock.URL("/?foo=bar"))) == nil

*/
func (key Param) Maybe(lens optics.Lens) Endpoint {
	return func(req Input) error {
		if opt, exists := req.Params().Get(string(key)); exists {
			req.Context().Put(lens, opt)
		}
		return nil
	}
}

/*

JSON matches a param key to closed struct.
It assumes that key holds JSON value as url encoded string
*/
func (key Param) JSON(lens optics.Lens) Endpoint {
	return func(req Input) error {
		opt, exists := req.Params().Get(string(key))
		if !exists {
			return NoMatch{}
		}

		str, err := url.QueryUnescape(opt)
		if err != nil {
			return NoMatch{}
		}

		return req.Context().Put(lens, str)
	}
}

/*

MaybeJSON matches a param key to closed struct.
It assumes that key holds JSON value as url encoded string.
It does not fail if key is not defined.
*/
func (key Param) MaybeJSON(lens optics.Lens) Endpoint {
	return func(req Input) error {
		opt, exists := req.Params().Get(string(key))
		if !exists {
			return nil
		}

		str, err := url.QueryUnescape(opt)
		if err != nil {
			return nil
		}

		req.Context().Put(lens, str)
		return nil
	}
}
