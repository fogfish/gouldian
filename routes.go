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
	"net/http"
	"strings"
)

/*
Route is endpoint with routing metadata: HTTP verb, path and the endpoint
itself. The metadata allows the server to distinguish "no such resource" (404)
from "no such method" (405) and to respond to OPTIONS requests.
*/
type Route struct {
	// Method of HTTP request (GET, POST, ...), Any matches every verb
	Method string
	// Path is a sequence of segments: literal, ":" for a variable captured
	// by a lens, "_" for any segment, "*" for the rest of the path.
	Path     []string
	Endpoint Endpoint
}

// String is "METHOD /path", ":" is a variable, "_" any segment and "*" the
// rest of the path.
func (r Route) String() string {
	method := strings.ToUpper(r.Method)
	if method == "" || method == Any {
		method = "ANY"
	}
	return method + " /" + strings.Join(r.Path, "/")
}

// Routable is a lazy route, see URI, GET, POST, ...
type Routable func() Route

// Routes is a routing table. It is built once and it is safe for concurrent use.
type Routes struct {
	root *node
}

type node struct {
	lit     []literal
	param   *node // :
	any     *node // _
	rest    *node // *
	methods []handler
}

type literal struct {
	segment string
	node    *node
}

type handler struct {
	verb     string
	endpoint Endpoint
}

// NewRoutes creates new routing table
func NewRoutes(seq ...Routable) *Routes {
	r := &Routes{root: &node{}}
	for _, route := range seq {
		r.append(route())
	}
	return r
}

func (r *Routes) append(route Route) {
	n := r.root
	for i, segment := range route.Path {
		switch segment {
		case ":":
			if n.param == nil {
				n.param = &node{}
			}
			n = n.param
		case "_":
			if n.any == nil {
				n.any = &node{}
			}
			n = n.any
		case "*":
			if i != len(route.Path)-1 {
				panic(fmt.Errorf("wildcard * must be the last segment of the path %v", route.Path))
			}
			if n.rest == nil {
				n.rest = &node{}
			}
			n = n.rest
		default:
			n = n.literal(segment)
		}
	}

	verb := strings.ToUpper(route.Method)
	if verb == "" {
		verb = Any
	}

	for i := range n.methods {
		if n.methods[i].verb == verb {
			n.methods[i].endpoint = n.methods[i].endpoint.Or(route.Endpoint)
			return
		}
	}
	n.methods = append(n.methods, handler{verb, route.Endpoint})
}

func (n *node) literal(segment string) *node {
	for _, l := range n.lit {
		if l.segment == segment {
			return l.node
		}
	}
	child := &node{}
	n.lit = append(n.lit, literal{segment, child})
	return child
}

// handler finds the endpoint for HTTP verb
func (n *node) handler(verb string) Endpoint {
	var wildcard, get Endpoint
	for _, h := range n.methods {
		switch h.verb {
		case verb:
			return h.endpoint
		case Any:
			wildcard = h.endpoint
		case http.MethodGet:
			get = h.endpoint
		}
	}

	if wildcard != nil {
		return wildcard
	}
	if verb == http.MethodHead {
		return get
	}
	return nil
}

// match finds the endpoint(s) for the path and executes it. The method looks
// through the alternatives (literal, variable, any, rest) in the order of
// the priority until an endpoint accepts the request. The outcome other than
// NoMatch is terminal.
func (n *node) match(ctx *Context, path string) error {
	if path == "" {
		return n.run(ctx)
	}

	// split the path to the first segment and the rest
	i := 0
	for i < len(path) && path[i] != '/' {
		i++
	}
	segment, rest := path[:i], ""
	if i < len(path) {
		rest = path[i+1:]
	}
	reason := NoMatchPath

	for _, l := range n.lit {
		if l.segment == segment {
			err := l.node.match(ctx, rest)
			if no, ok := err.(NoMatch); !ok {
				return err
			} else {
				reason = max(reason, no)
			}
			break
		}
	}

	if segment == "" {
		return reason
	}

	if n.param != nil {
		mark := len(ctx.values)
		ctx.values = append(ctx.values, segment)
		err := n.param.match(ctx, rest)
		if no, ok := err.(NoMatch); !ok {
			return err
		} else {
			reason = max(reason, no)
		}
		ctx.values = ctx.values[:mark]
	}

	if n.any != nil {
		err := n.any.match(ctx, rest)
		if no, ok := err.(NoMatch); !ok {
			return err
		} else {
			reason = max(reason, no)
		}
	}

	if n.rest != nil {
		mark := len(ctx.values)
		ctx.values = append(ctx.values, path)
		err := n.rest.run(ctx)
		if no, ok := err.(NoMatch); !ok {
			return err
		} else {
			reason = max(reason, no)
		}
		ctx.values = ctx.values[:mark]
	}

	return reason
}

// run the endpoint of HTTP verb
func (n *node) run(ctx *Context) error {
	if len(n.methods) == 0 {
		return NoMatchPath
	}

	endpoint := n.handler(ctx.Request.Method)
	if endpoint == nil {
		for _, h := range n.methods {
			ctx.allow = append(ctx.allow, h.verb)
		}
		return NoMatchMethod
	}

	mark := len(ctx.morphism)
	err := endpoint(ctx)
	if _, ok := err.(NoMatch); ok {
		ctx.rollback(mark)
	}
	return err
}

// Exec the routing table against the request
func (r *Routes) exec(ctx *Context, path string) error {
	ctx.free()
	ctx.allow = ctx.allow[:0]
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	// /path/ is not /path, see TrailingSlash policy of the server
	if len(path) > 0 && path[len(path)-1] == '/' {
		return NoMatchPath
	}
	return r.root.match(ctx, path)
}

// Endpoint converts routing table to Endpoint
func (r *Routes) Endpoint() Endpoint {
	return func(ctx *Context) error {
		if ctx.Request == nil {
			return NoMatchAny
		}

		return r.exec(ctx, ctx.Request.URL.Path)
	}
}
