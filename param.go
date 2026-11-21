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
		return ctx.put(SourceParam, "", lens, ctx.Request.URL.RawQuery)
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

An absent param or differing literal is NoMatch. A param that is present but
cannot be decoded by the lens is Failure (400).
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
decode value into Golang type. The Endpoint matches if the param is not defined.
A param that is present but cannot be decoded to the target type is Failure.

	type myT struct{ Val string }

	x := µ.Optics1[myT, string]()
	e := µ.GET( µ.ParamMaybe("foo", x) )
	e(mock.Input(mock.URL("/?foo=bar"))) == nil
	e(mock.Input(mock.URL("/?foo"))) == nil
	e(mock.Input(mock.URL("/"))) == nil
*/
func ParamMaybe(key string, lens Lens) Endpoint {
	return func(ctx *Context) error {
		if _, exists := ctx.Query().Get(key); exists {
			return putParam(ctx, key, lens)
		}
		return nil
	}
}

/*
ParamJSON matches a param key to struct.
It assumes that key holds JSON value as url encoded string
*/
func ParamJSON(key string, lens Lens) Endpoint {
	return func(ctx *Context) error {
		opt, exists := ctx.Query().Get(key)
		if !exists {
			return NoMatchParam
		}

		return putParamJSON(ctx, key, lens, opt)
	}
}

/*
ParamMaybeJSON matches a param key to closed struct.
It assumes that key holds JSON value as url encoded string.
It does not fail if key is not defined, a malformed value is Failure.
*/
func ParamMaybeJSON(key string, lens Lens) Endpoint {
	return func(ctx *Context) error {
		opt, exists := ctx.Query().Get(key)
		if !exists {
			return nil
		}

		return putParamJSON(ctx, key, lens, opt)
	}
}

func putParamJSON(ctx *Context, key string, lens Lens, opt string) error {
	str, err := url.QueryUnescape(opt)
	if err != nil {
		return newFailure(SourceParam, key, err)
	}

	return ctx.put(SourceParam, key, lens, str)
}

// putParam binds all values of the param if lens is multi-value (slice),
// the first value otherwise.
func putParam(ctx *Context, key string, lens Lens) error {
	if lens.Multi() {
		vals, _ := ctx.Query().GetAll(key)
		return ctx.putMulti(SourceParam, key, lens, vals)
	}

	opt, _ := ctx.Query().Get(key)
	return ctx.put(SourceParam, key, lens, opt)
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
		opt, exists := ctx.Query().Get(string(key))
		if exists && opt == val {
			return nil
		}
		return NoMatchParam
	}
}

/*
Any is a wildcard matcher of param key. It fails if key is not defined.
*/
func (key param) Any(ctx *Context) error {
	_, exists := ctx.Query().Get(string(key))
	if exists {
		return nil
	}
	return NoMatchParam
}

/*
To matches param value to the request context. It uses lens abstraction to
decode value into Golang type. An absent param is NoMatch, a param that
cannot be decoded to the target type is Failure.
*/
func (key param) To(lens Lens) Endpoint {
	return func(ctx *Context) error {
		if _, exists := ctx.Query().Get(string(key)); exists {
			return putParam(ctx, string(key), lens)
		}
		return NoMatchParam
	}
}
