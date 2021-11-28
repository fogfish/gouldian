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
	return v(verb).match
	// if verb == Any {
	// 	return func(ctx *Context) error {
	// 		ctx.free()
	// 		return nil
	// 	}
	// }

	// return func(ctx *Context) error {
	// 	if ctx.Request == nil {
	// 		return ErrNoMatch
	// 	}

	// 	if ctx.Request.Method == verb {
	// 		ctx.free()
	// 		return nil
	// 	}

	// 	return ErrNoMatch
	// }
}

type v string

func (verb v) match(ctx *Context) error {
	if ctx.Request.Method == string(verb) {
		ctx.free()
		return nil
	}

	return ErrNoMatch
}

// Method2 ...
// func Method2(verb string) Builder {
// 	return func(root *Node) *Node {
// 		return root.mkByID(verb, func(ctx *Context) error {
// 			if ctx.Request.Method == verb {
// 				ctx.free()
// 				return nil
// 			}

// 			return ErrNoMatch
// 		})
// 	}
// }

// Path2 ...
func Path2(path ...string) Builder {
	return func(root *Node) (n *Node) {
		root.appendEndpoint(path, func(c *Context) error { return nil })
		// Put(root, path)
		return root
		// n = root
		// for i, term := range seq {
		// 	n = n.mkByID(term, func(c *Context) error {
		// 		if len(c.segments) < i+1 {
		// 			return ErrNoMatch
		// 		}

		// 		if len(c.segments[i]) != len(term) {
		// 			return ErrNoMatch
		// 		}

		// 		if c.segments[i] == term {
		// 			return nil
		// 		}
		// 		return ErrNoMatch
		// 	})
		// }
		// return
	}
}

/*
TODO: How to make a node selector so that

join takes root, but each next is hierarchical "node injection"

GET -> "a" -> "b" -> "c" -> ...

µ.Join(
	µ.Method(...)
	µ.Path(...)
	µ....
)

*/

// func Path22(path string) *Node {
// 	return &Node{
// 		Endpoint: p(path).match2,
// 		Children: make([]*Node, 0),
// 	}
// }

// func Path23(path string) *Node {
// 	return &Node{
// 		Endpoint: p(path).match3,
// 		Children: make([]*Node, 0),
// 	}
// }

// type p string

// func (path p) match1(ctx *Context) error {
// 	if ctx.segments[0] == string(path) {
// 		return nil
// 	}

// 	// if strings.HasPrefix(ctx.Request.RequestURI, string(path)) {
// 	// 	return nil
// 	// }
// 	// if ctx.Request.RequestURI == string(path) {
// 	// 	return nil
// 	// }

// 	return ErrNoMatch
// }

// func (path p) match2(ctx *Context) error {
// 	// fmt.Println(ctx.segments[1], string(path))
// 	if ctx.segments[1] == string(path) {
// 		return nil
// 	}

// 	// if strings.HasPrefix(ctx.Request.RequestURI, string(path)) {
// 	// 	return nil
// 	// }
// 	// if ctx.Request.RequestURI == string(path) {
// 	// 	return nil
// 	// }

// 	return ErrNoMatch
// }

// func (path p) match3(ctx *Context) error {
// 	// fmt.Println(ctx.segments[1], string(path))
// 	if ctx.segments[2] == string(path) {
// 		return nil
// 	}

// 	// if strings.HasPrefix(ctx.Request.RequestURI, string(path)) {
// 	// 	return nil
// 	// }
// 	// if ctx.Request.RequestURI == string(path) {
// 	// 	return nil
// 	// }

// 	return ErrNoMatch
// }

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
