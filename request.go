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
	"encoding/json"
	"fmt"
	"net/http"
	"unsafe"
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
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
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
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
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
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
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
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
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
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
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
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
	  ...
	)
	e(mock.Input(mock.Method("PUT"))) == nil
	e(mock.Input(mock.Method("OTHER"))) == nil
*/
func ANY(path Routable, arrows ...Endpoint) Routable {
	seq := append(Endpoints{Method(Any)}, arrows...)
	return Route(path, seq...)
}

// HTTP composes Endpoints into Routable
func HTTP(verb string, path Routable, arrows ...Endpoint) Routable {
	seq := append(Endpoints{Method(verb)}, arrows...)
	return Route(path, seq...)
}

// Method is an endpoint to match HTTP verb request
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

// Body decodes HTTP request body and lifts it to the structure
func Body(lens Lens) Endpoint {
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

// FMap applies clojure to matched HTTP request,
// taking the execution context as the input to closure
func FMap[A any](f func(*Context, *A) error) Endpoint {
	return func(req *Context) error {
		var a A
		if err := FromContext(req, &a); err != nil {
			out := NewOutput(http.StatusBadRequest)
			out.SetIssue(err)
			return out
		}

		return f(req, &a)
	}
}

// Map applies clojure to matched HTTP request,
// taking the execution context and matched parameters as the input to closure.
// The output is always returned as JSON.
func Map[A, B any](f func(*Context, *A) (*B, error)) Endpoint {
	return func(req *Context) error {
		var a A
		if err := FromContext(req, &a); err != nil {
			out := NewOutput(http.StatusBadRequest)
			out.SetIssue(err)
			return out
		}

		b, err := f(req, &a)
		if err != nil {
			return err
		}

		out := NewOutput(http.StatusOK)

		val, err := json.Marshal(b)
		if err != nil {
			out.SetIssue(fmt.Errorf("Serialization is failed for <%T>", val))
			return out
		}

		out.SetHeader("Content-Type", "application/json")
		out.Body = string(val)
		return out
	}
}
