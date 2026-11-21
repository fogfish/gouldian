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
)

const (
	// Any constant matches any term
	Any = "_"
)

// HTTP composes Endpoints into Routable. The verb is the routing metadata,
// the server responds with 405 if the path is known but not the verb.
func HTTP(verb string, path Routable, arrows ...Endpoint) Routable {
	return func() Route {
		route := path()
		if len(arrows) > 0 {
			route.Endpoint = append(Endpoints{route.Endpoint}, arrows...).Join
		}
		route.Method = verb
		return route
	}
}

/*
DELETE composes Endpoints into Routable that matches HTTP DELETE request.

	e := µ.DELETE(
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
	  ...
	)
*/
func DELETE(path Routable, arrows ...Endpoint) Routable {
	return HTTP(http.MethodDelete, path, arrows...)
}

/*
GET composes Endpoints into Routable that matches HTTP GET request.

	e := µ.GET(
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
	  ...
	)
*/
func GET(path Routable, arrows ...Endpoint) Routable {
	return HTTP(http.MethodGet, path, arrows...)
}

/*
PATCH composes Endpoints into Routable that matches HTTP PATCH request.

	e := µ.PATCH(
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
	  ...
	)
*/
func PATCH(path Routable, arrows ...Endpoint) Routable {
	return HTTP(http.MethodPatch, path, arrows...)
}

/*
POST composes Endpoints into Routable that matches HTTP POST request.

	e := µ.POST(
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
	  ...
	)
*/
func POST(path Routable, arrows ...Endpoint) Routable {
	return HTTP(http.MethodPost, path, arrows...)
}

/*
PUT composes Endpoints into Routable that matches HTTP PUT request.

	e := µ.PUT(
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
	  ...
	)
*/
func PUT(path Routable, arrows ...Endpoint) Routable {
	return HTTP(http.MethodPut, path, arrows...)
}

/*
ANY composes Endpoints into Routable that matches HTTP any request.

	e := µ.ANY(
	  µ.URI(µ.Path("foo"), µ.Path("bar")),
	  ...
	)
*/
func ANY(path Routable, arrows ...Endpoint) Routable {
	return HTTP(Any, path, arrows...)
}

// Method is an endpoint to match HTTP verb request. Routes built by GET, POST,
// etc. carry the verb as metadata, use the endpoint only to constrain
// a request inside the route (e.g. ANY).
func Method(verb string) Endpoint {
	if verb == Any {
		return func(ctx *Context) error {
			return nil
		}
	}

	return func(ctx *Context) error {
		if ctx.Request == nil {
			return NoMatchMethod
		}

		if ctx.Request.Method == verb {
			return nil
		}

		return NoMatchMethod
	}
}

// Body decodes HTTP request body and lifts it to the structure.
// The body is read once and shared by all endpoints of the request, so
// several lenses (e.g. raw bytes and typed struct) can read the same body.
//
// An empty body is Failure (400, "body is required"). A body over the limit
// of the server is Failure (413). A body that cannot be decoded by the lens
// is Failure (400) with JSON pointer to the invalid value.
func Body(lens Lens) Endpoint {
	return func(ctx *Context) error {
		buf, err := ctx.body()
		if err != nil {
			return err
		}

		if len(buf) == 0 {
			return newFailure(SourceBody, "", ErrBodyRequired)
		}

		return ctx.putBody(lens, buf)
	}
}

// BodyMaybe is like Body, but an empty body is not an error.
func BodyMaybe(lens Lens) Endpoint {
	return func(ctx *Context) error {
		buf, err := ctx.body()
		if err != nil {
			return err
		}

		if len(buf) == 0 {
			return nil
		}

		return ctx.putBody(lens, buf)
	}
}

// FMap applies clojure to matched HTTP request,
// taking the execution context as the input to closure
func FMap[A any](f func(*Context, *A) error) Endpoint {
	return func(req *Context) error {
		var a A
		if err := FromContext(req, &a); err != nil {
			return err
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
			return err
		}

		b, err := f(req, &a)
		if err != nil {
			return err
		}

		val, err := json.Marshal(b)
		if err != nil {
			out := NewOutput(http.StatusInternalServerError)
			out.SetIssue(fmt.Errorf("serialization is failed for <%T>: %w", b, err))
			return out
		}

		out := NewOutput(http.StatusOK)
		out.SetHeader("Content-Type", "application/json")
		out.Body = val
		return out
	}
}
