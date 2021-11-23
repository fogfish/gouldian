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
	"unsafe"

	"github.com/fogfish/gouldian/optics"
)

const (
	// Any constant matches any term
	Any = "_"
)

/*

DELETE composes product Endpoint match HTTP DELETE request.
  e := µ.DELETE()
  e(mock.Input(mock.Method("DELETE"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func DELETE(arrows ...Endpoint) Endpoint {
	return Method("DELETE").Then(Join(arrows...))
}

/*

GET composes product Endpoint match HTTP GET request.
  e := µ.GET()
  e(mock.Input(mock.Method("GET"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func GET(arrows ...Endpoint) Endpoint {
	return Method("GET").Then(Join(arrows...))
}

/*

PATCH composes product Endpoint match HTTP PATCH request.
  e := µ.PATCH()
  e(mock.Input(mock.Method("PATCH"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func PATCH(arrows ...Endpoint) Endpoint {
	return Method("PATCH").Then(Join(arrows...))
}

/*

POST composes product Endpoint match HTTP POST request.
  e := µ.POST()
  e(mock.Input(mock.Method("POST"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func POST(arrows ...Endpoint) Endpoint {
	return Method("POST").Then(Join(arrows...))
}

/*

PUT composes product Endpoint match HTTP PUT request.
  e := µ.PUT()
  e(mock.Input(mock.Method("PUT"))) == nil
  e(mock.Input(mock.Method("OTHER"))) != nil
*/
func PUT(arrows ...Endpoint) Endpoint {
	return Method("PUT").Then(Join(arrows...))
}

/*

ANY composes product Endpoint match HTTP PUT request.
  e := µ.ANY()
  e(mock.Input(mock.Method("PUT"))) == nil
  e(mock.Input(mock.Method("OTHER"))) == nil
*/
func ANY(arrows ...Endpoint) Endpoint {
	return Method(Any).Then(Join(arrows...))
}

/*

Method is an endpoint to match HTTP verb request
*/
func Method(verb string) Endpoint {
	if verb == Any {
		return func(ctx *Context) error {
			// req.Context.Free()
			ctx.free()
			return nil
		}
	}

	return func(ctx *Context) error {
		if ctx.Request.Method == verb {
			// req.Context.Free()
			ctx.free()
			return nil
		}
		return NoMatch{}
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
				return NoMatch{}
			}

			return ctx.Put(lens, *(*string)(unsafe.Pointer(&ctx.payload)))
		}

		return NoMatch{}
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
			return NoMatch{}
		}

		if key(ctx.JWT) != val {
			return NoMatch{}
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
			return NoMatch{}
		}

		if val := key(ctx.JWT); val != "" {
			return ctx.Put(lens, val)
		}

		return NoMatch{}
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
			return NoMatch{}
		}

		if val := key(ctx.JWT); val != "" {
			ctx.Put(lens, val)
		}

		return nil
	}
}
