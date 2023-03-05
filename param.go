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

	"github.com/fogfish/gouldian/v2/internal/optics"
)

// Params lifts all Query parameters to struct
//
//	type MyRequest struct {
//		Params MyType `content:"form"`
//	}
//	var params = µ.Optics1[MyRequest, MyType]()
//	µ.GET(µ.Params(params))
func Params(lens Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.params == nil {
			ctx.params = Query(ctx.Request.URL.Query())
		}

		return ctx.Put(lens, ctx.Request.URL.RawQuery)
	}
}

/*
Param combinator defines primitives to match query param in the HTTP requests.

	  endpoint := µ.GET(
	    µ.Param("foo", "bar"),
	  )

	  endpoint(
	    mock.Input(
				mock.URL("/?foo=bar")
	    )
	  ) == nil
*/
func Param[T Pattern](key string, val T) Endpoint {
	switch v := any(val).(type) {
	case string:
		return param(key).Is(v)
	case Lens:
		return param(key).To(v)
	default:
		panic("type system failure")
	}
}

/*
ParamAny is a wildcard matcher of param key. It fails if key is not defined.

	e := µ.GET( µ.ParamAny("foo") )
	e(mock.Input(mock.URL("/?foo"))) == nil
	e(mock.Input(mock.URL("/?foo=bar"))) == nil
	e(mock.Input(mock.URL("/?foo=baz"))) == nil
	e(mock.Input()) != nil
*/
func ParamAny(key string) Endpoint {
	return param(key).Any
}

/*
ParamMaybe matches param value to the request context. It uses lens abstraction to
decode value into Golang type. The Endpoint does not cause no-match
if header value cannot be decoded to the target type. See optics.Lens type for details.

	  type myT struct{ Val string }

	  x := µ.Optics1[myT, string]()
	  e := µ.GET( µ.ParamMaybe("foo", x) )
	  e(mock.Input(mock.URL("/?foo=bar"))) == nil
		e(mock.Input(mock.URL("/?foo"))) == nil
		e(mock.Input(mock.URL("/"))) == nil
*/
func ParamMaybe(key string, lens Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.params == nil {
			ctx.params = Query(ctx.Request.URL.Query())
		}

		if opt, exists := ctx.params.Get(string(key)); exists {
			ctx.Put(lens, opt)
		}
		return nil
	}
}

/*
JSON matches a param key to struct.
It assumes that key holds JSON value as url encoded string
*/
func ParamJSON(key string, lens Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.params == nil {
			ctx.params = Query(ctx.Request.URL.Query())
		}

		opt, exists := ctx.params.Get(string(key))
		if !exists {
			return ErrNoMatch
		}

		str, err := url.QueryUnescape(opt)
		if err != nil {
			return ErrNoMatch
		}

		return ctx.Put(lens, str)
	}
}

/*
MaybeJSON matches a param key to closed struct.
It assumes that key holds JSON value as url encoded string.
It does not fail if key is not defined.
*/
func ParamMaybeJSON(key string, lens Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.params == nil {
			ctx.params = Query(ctx.Request.URL.Query())
		}

		opt, exists := ctx.params.Get(string(key))
		if !exists {
			return nil
		}

		str, err := url.QueryUnescape(opt)
		if err != nil {
			return nil
		}

		ctx.Put(lens, str)
		return nil
	}
}

// Internal type
type param string

/*
Is matches a param key to defined literal value
*/
func (key param) Is(val string) Endpoint {
	if val == Any {
		return key.Any
	}

	return func(ctx *Context) error {
		if ctx.params == nil {
			ctx.params = Query(ctx.Request.URL.Query())
		}

		opt, exists := ctx.params.Get(string(key))
		if exists && opt == val {
			return nil
		}
		return ErrNoMatch
	}
}

/*
Any is a wildcard matcher of param key. It fails if key is not defined.
*/
func (key param) Any(ctx *Context) error {
	if ctx.params == nil {
		ctx.params = Query(ctx.Request.URL.Query())
	}

	_, exists := ctx.params.Get(string(key))
	if exists {
		return nil
	}
	return ErrNoMatch
}

/*
To matches param value to the request context. It uses lens abstraction to
decode value into Golang type. The Endpoint causes no-match if param
value cannot be decoded to the target type. See optics.Lens type for details.
*/
func (key param) To(lens optics.Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.params == nil {
			ctx.params = Query(ctx.Request.URL.Query())
		}

		opt, exists := ctx.params.Get(string(key))
		if exists {
			return ctx.Put(lens, opt)
		}
		return ErrNoMatch
	}
}
