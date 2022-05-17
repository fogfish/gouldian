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
	"unsafe"

	"github.com/fogfish/gouldian/optics"
)

const (
	// Any constant matches any term
	Any = "_"
)

/*

Route converts sequence ot Endpoints into Routable element
*/
func Route(
	path Routable,
	seq ...Endpoint,
) Routable {
	return func() ([]string, Endpoint) {
		route, pathEndpoint := path()
		endpoints := append(Endpoints{pathEndpoint}, seq...)
		return route, endpoints.Join
	}
}

/*

DELETE composes Endpoints into Routable that matches HTTP DELETE request.
  e := µ.DELETE(
    µ.Path("foo", "bar"),
    ...
  )
  e(mock.Input(mock.Method("DELETE"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func DELETE(path Routable, arrows ...Endpoint) Routable {
	seq := append(Endpoints{Method("DELETE")}, arrows...)
	return Route(path, seq...)
}

/*

GET composes Endpoints into Routable that matches HTTP GET request.
  e := µ.GET(
    µ.Path("foo", "bar"),
    ...
  )
  e(mock.Input(mock.Method("GET"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func GET(path Routable, arrows ...Endpoint) Routable {
	seq := append(Endpoints{Method("GET")}, arrows...)
	return Route(path, seq...)
}

/*

PATCH composes Endpoints into Routable that matches HTTP PATCH request.
  e := µ.PATCH(
    µ.Path("foo", "bar"),
    ...
  )
  e(mock.Input(mock.Method("PATCH"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func PATCH(path Routable, arrows ...Endpoint) Routable {
	seq := append(Endpoints{Method("PATCH")}, arrows...)
	return Route(path, seq...)
}

/*

POST composes Endpoints into Routable that matches HTTP POST request.
  e := µ.POST(
    µ.Path("foo", "bar"),
    ...
  )
  e(mock.Input(mock.Method("POST"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func POST(path Routable, arrows ...Endpoint) Routable {
	seq := append(Endpoints{Method("POST")}, arrows...)
	return Route(path, seq...)
}

/*

PUT composes Endpoints into Routable that matches HTTP PUT request.
  e := µ.PUT(
    µ.Path("foo", "bar"),
    ...
  )
  e(mock.Input(mock.Method("PUT"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func PUT(path Routable, arrows ...Endpoint) Routable {
	seq := append(Endpoints{Method("PUT")}, arrows...)
	return Route(path, seq...)
}

/*

ANY composes Endpoints into Routable that matches HTTP any request.
  e := µ.ANY(
    µ.Path("foo", "bar"),
    ...
  )
  e(mock.Input(mock.Method("PUT"))) == nil
  e(mock.Input(mock.Method("OTHER"))) == nil
*/
func ANY(path Routable, arrows ...Endpoint) Routable {
	seq := append(Endpoints{Method(Any)}, arrows...)
	return Route(path, seq...)
}

/*

Method is an endpoint to match HTTP verb request
*/
func Method(verb string) Endpoint {
	if verb == Any {
		return func(ctx *Context) error {
			return nil
		}
	}

	return func(ctx *Context) error {
		if ctx.Request == nil {
			return ErrNoMatch
		}

		if ctx.Request.Method == verb {
			return nil
		}

		return ErrNoMatch
	}
}

/*

Body decodes HTTP request body and lifts it to the structure
*/
func Body(lens optics.Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.payload == nil {
			if err := ctx.cacheBody(); err != nil {
				return err
			}

			if ctx.payload == nil {
				return ErrNoMatch
			}

			return ctx.Put(lens, *(*string)(unsafe.Pointer(&ctx.payload)))
		}

		return ErrNoMatch
	}
}

/*

FMap applies clojure to matched HTTP request,
taking the execution context as the input to closure
*/
func FMap(f func(*Context) error) Endpoint {
	return func(req *Context) error { return f(req) }
}

/*

Access type defines primitives to match JWT token in the HTTP requests.

  endpoint := µ.GET(
    µ.Access(µ.JWT.Username).Is("joedoe"),
  )

  endpoint(
    mock.Input(
			mock.JWT(µ.JWT{"username": "joedoe"})
    )
  ) == nil

*/
type Access func(JWT) string

/*

Is matches a key of JWT to defined literal value

  e := µ.GET( µ.Access(µ.JWT.Username).Is("joedoe") )
  e(mock.Input(mock.JWT(µ.JWT{"username": "joedoe"})) == nil
  e(mock.Input(mock.JWT(µ.JWT{"username": "landau"})) != nil
*/
func (key Access) Is(val string) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return ErrNoMatch
		}

		if key(ctx.JWT) != val {
			return ErrNoMatch
		}

		return nil
	}
}

/*

OneOf matches a key of JWT if it contains one of the tokens

  µ.GET( µ.Access(µ.JWT.Scope).OneOf("ro", "rw") )
*/
func (key Access) OneOf(vals ...string) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return ErrNoMatch
		}

		val := key(ctx.JWT)
		for _, x := range vals {
			if strings.Index(val, x) != -1 {
				return nil
			}
		}

		return ErrNoMatch
	}
}

/*

AllOf matches a key of JWT if it contains one of the tokens

  µ.GET( µ.Access(µ.JWT.Scope).AllOf("ro", "rw") )
*/
func (key Access) AllOf(vals ...string) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return ErrNoMatch
		}

		val := key(ctx.JWT)
		for _, x := range vals {
			if strings.Index(val, x) == -1 {
				return ErrNoMatch
			}
		}

		return nil
	}
}

/*

To matches key of JWT value to the request context. It uses lens abstraction to
decode value into Golang type. The Endpoint causes no-match if param
value cannot be decoded to the target type. See optics.Lens type for details.

  type MyT struct{ Username string }

  username := optics.Lenses1(MyT{})
  e := µ.GET( µ.Access(µ.JWT.Sub).To(username) )
  e(mock.Input(mock.JWT(µ.JWT{"username": "joedoe"}))) == nil
*/
func (key Access) To(lens optics.Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return ErrNoMatch
		}

		if val := key(ctx.JWT); val != "" {
			return ctx.Put(lens, val)
		}

		return ErrNoMatch
	}
}

/*

Maybe matches key of JWT to the request context. It uses lens abstraction to
decode value into Golang type. The Endpoint does not cause no-match
if header value cannot be decoded to the target type. See optics.Lens type for details.

  type MyT struct{ Username string }

  userna,e := optics.Lenses1(MyT{})
  e := µ.GET( µ.Access(µ.JWT.Sub).Maybe(username) )
  e(mock.Input(mock.JWT(µ.JWT{"username": "joedoe"}))) == nil

*/
func (key Access) Maybe(lens optics.Lens) Endpoint {
	return func(ctx *Context) error {
		if ctx.JWT == nil {
			return ErrNoMatch
		}

		if val := key(ctx.JWT); val != "" {
			ctx.Put(lens, val)
		}

		return nil
	}
}
