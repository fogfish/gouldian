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

/*

Package param defines primitives to match URL Query parameters of HTTP requests.

	import "github.com/fogfish/gouldian/param"

	endpoint := µ.GET( µ.Param(param.Is("foo", "bar"), ...) )
	endpoint(mock.Input(mock.URL("/?foo=bar"))) == nil

*/
package gouldian

import (
	"net/url"

	"github.com/fogfish/gouldian/optics"
)

/*

Param matches presence of query param in the request or match its entire content.

	µ.Param("foo").Is("bar")
*/
type Param string

/*

Is matches a param key to defined literal value

  e := µ.Param(param.Is("foo", "bar"))
  e(mock.Input(mock.URL("/?foo=bar"))) == nil
  e(mock.Input(mock.URL("/?bar=foo"))) != nil
*/
func (key Param) Is(val string) Endpoint {
	if val == Any {
		return key.Any()
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

  e := µ.Param(param.Any("foo"))
  e(mock.Input(mock.URL("/?foo"))) == nil
  e(mock.Input(mock.URL("/?foo=bar"))) == nil
  e(mock.Input(mock.URL("/?foo=baz"))) == nil
  e(mock.Input(mock.URL("/?bar=foo"))) != nil
*/
func (key Param) Any() Endpoint {
	return func(req Input) error {
		_, exists := req.Params().Get(string(key))
		if exists {
			return nil
		}
		return NoMatch{}
	}
}

func (key Param) To(lens optics.Lens) Endpoint {
	return func(req Input) error {
		if opt, exists := req.Params().Get(string(key)); exists {
			return req.Context().Put(lens, opt)
		}
		return NoMatch{}
	}
}

func (key Param) Maybe(lens optics.Lens) Endpoint {
	return func(req Input) error {
		if opt, exists := req.Params().Get(string(key)); exists {
			req.Context().Put(lens, opt)
		}
		return nil
	}
}

/*

String matches a param key to closed variable of string type.
It fails if key is not defined.

  e := µ.Param(param.String("foo", FOO))
  e(mock.Input(mock.URL("/?foo=bar"))) == nil && *ctx.String(FOO) == "bar"
  e(mock.Input(mock.URL("/?foo=1"))) == nil && *ctx.String(FOO) == "1"
*/
func (key Param) String(lens optics.Lens) Endpoint {
	return func(req Input) error {
		if opt, exists := req.Params().Get(string(key)); exists {
			return req.Context().Put(lens, opt)
		}
		return NoMatch{}
	}
}

/*

MaybeString matches a param key to closed variable of string type.
It does not fail if key is not defined.

  const FOO µ.Symbol = iota
  e := µ.Param(param.String("foo", FOO))
  e(mock.Input(mock.URL("/?foo=bar"))) == nil && *ctx.String(FOO) == "bar"
  e(mock.Input(mock.URL("/?bar=1"))) == nil && *ctx.String(FOO) == ""
*/
func (key Param) MaybeString(lens optics.Lens) Endpoint {
	return func(req Input) error {
		if opt, exists := req.Params().Get(string(key)); exists {
			req.Context().Put(lens, opt)
		}
		return nil
	}
}

/*

Int matches a param key to closed variable of int type.
It fails if key is not defined.

  const FOO µ.Symbol = iota
  e := µ.Param(param.Int("foo", FOO))
  e(mock.Input(mock.URL("/?foo=1"))) == nil && *ctx.Int(FOO) == 1
  e(mock.Input(mock.URL("/?foo=bar"))) != nil
*/
func (key Param) Int(lens optics.Lens) Endpoint {
	return func(req Input) error {
		opt, exists := req.Params().Get(string(key))
		if !exists {
			return NoMatch{}
		}

		return req.Context().Put(lens, opt)
	}
}

/*

MaybeInt matches a param key to closed variable of int type.
It does not fail if key is not defined.

  const FOO µ.Symbol = iota
  e := µ.GET( µ.Param(param.MaybeInt("foo", &value)) )
  e(mock.Input(mock.URL("/?foo=1"))) == nil && *ctx.Int(FOO) == 1
  e(mock.Input(mock.URL("/?foo=bar"))) == nil && *ctx.Int(FOO) == 0
*/
func (key Param) MaybeInt(lens optics.Lens) Endpoint {
	return func(req Input) error {
		opt, exists := req.Params().Get(string(key))
		if !exists {
			return nil
		}

		req.Context().Put(lens, opt)
		return nil
	}
}

/*

Float matches a param key to closed variable of float64 type.
It fails if key is not defined.

  const FOO µ.Symbol = iota
  e := µ.GET( µ.Param(param.Float("foo", &value)) )
  e(mock.Input(mock.URL("/?foo=1"))) == nil && value == 1
  e(mock.Input(mock.URL("/?foo=bar"))) != nil

*/
func (key Param) Float(lens optics.Lens) Endpoint {
	return func(req Input) error {
		opt, exists := req.Params().Get(string(key))
		if !exists {
			return NoMatch{}
		}

		return req.Context().Put(lens, opt)
	}
}

/*

MaybeFloat matches a param key to closed variable of float64 type.
It does not fail if key is not defined.

  var value float64
  e := µ.GET( µ.Param(param.MaybeFloat("foo", &value)) )
  e(mock.Input(mock.URL("/?foo=1"))) == nil && value == 1
  e(mock.Input(mock.URL("/?foo=bar"))) == nil && value == 0

*/
func (key Param) MaybeFloat(lens optics.Lens) Endpoint {
	return func(req Input) error {
		opt, exists := req.Params().Get(string(key))
		if !exists {
			return nil
		}

		req.Context().Put(lens, opt)
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

		req.Context().Put(lens, str)
		return nil
	}
}

// MaybeJSON matches a param key to closed struct.
// It assumes that key holds JSON value as url encoded string.
// It does not fail if key is not defined.
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

/*

Or is a co-product of query param match arrows

  e := µ.Param(param.Or(param.Is("foo", "bar"), param.Is("bar", "foo")))
  e(mock.Input(mock.URL("/?foo=bar"))) == nil
  e(mock.Input(mock.URL("/?bar=foo"))) == nil
  e(mock.Input(mock.URL("/?foo=baz"))) != nil
*/
// func Or(arrows ...µ.ArrowParam) µ.ArrowParam {
// 	return func(ctx µ.Context, params µ.Params) error {
// 		for _, f := range arrows {
// 			if err := f(ctx, params); !errors.Is(err, µ.NoMatch{}) {
// 				return err
// 			}
// 		}
// 		return µ.NoMatch{}
// 	}
// }
