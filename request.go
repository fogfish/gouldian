//
//   Copyright 2019 Dmitry Kolesnikov, All Rights Reserved
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package gouldian

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/ajg/form"
)

// ArrowHeader is a type-safe definition of Header matcher
type ArrowHeader func(Context, Headers) error

// ArrowParam is a type-safe definition of URL Query matcher
type ArrowParam func(Context, Params) error

// ArrowPath is a type-safe definition of URL segment matcher
type ArrowPath func(Context, Segments) error

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
	return Method("*").Then(Join(arrows...))
}

// Method is an endpoint to match HTTP verb request
func Method(verb string) Endpoint {
	if verb == "*" {
		return func(req Input) error {
			req.Context().Free()
			return nil
		}
	}

	return func(req Input) error {
		if req.Method() == verb {
			req.Context().Free()
			return nil
		}
		return NoMatch{}
	}
}

/*

Path is an endpoint to match URL of HTTP request. The function takes
url matching primitives, which are defined by the package `path`.

  import "github.com/fogfish/gouldian/path"

  e := µ.GET( µ.Path(path.Is("foo")) )
  e(mock.Input(mock.URL("/foo"))) == nil
  e(mock.Input(mock.URL("/bar"))) != nil
*/
func Path(arrows ...ArrowPath) Endpoint {
	return func(req Input) error {
		seq := req.Resource()

		sz := len(seq)
		at := 0

		for _, f := range arrows {
			if sz <= at {
				return NoMatch{}
			}
			switch err := f(req.Context(), seq[at:]).(type) {
			case nil:
				at++
			case Match:
				at = at + err.N
			default:
				return err
			}
		}

		if sz > at {
			return NoMatch{}
		}

		return nil
	}
}

/*

Param is an endpoint to match URL Query parameters of HTTP request.
The function takes url query matching primitives, which are defined
by the package `param`.

  import "github.com/fogfish/gouldian/param"

  e := µ.GET( µ.Param(param.Is("foo", "bar")) )
  e(mock.Input(mock.URL("/?foo=bar"))) == nil
  e(mock.Input(mock.URL("/?foo=baz"))) != nil
*/
func Param(arrows ...ArrowParam) Endpoint {
	return func(req Input) error {
		for _, f := range arrows {
			if err := f(req.Context(), req.Params()); err != nil {
				return err
			}
		}
		return nil
	}
}

/*

Header is an endpoint to match Header(s) of HTTP request.
The function takes header matching primitives, which are defined
by the package `header`.

  import "github.com/fogfish/gouldian/header"

  e := µ.GET(
    µ.Header(
      param.Header("Content-Type", "application/json"),
    ),
  )
  Json := mock.Header("Content-Type", "application/json")
  e(mock.Input(Json)) == nil

  Text := mock.Header("Content-Type", "text/plain")
  e(mock.Input(Text)) != nil
*/
func Header(arrows ...ArrowHeader) Endpoint {
	return func(req Input) error {
		for _, f := range arrows {
			if err := f(req.Context(), req.Headers()); err != nil {
				return err
			}
		}
		return nil
	}
}

// JWT decodes token associated with the request.
// Endpoint fails if Authentication context is not found in the request.
func JWT(val *AccessToken) Endpoint {
	return func(req Input) error {
		if req.RequestContext.Authorizer == nil {
			return NoMatch{}
		}

		if jwt, isJwt := req.RequestContext.Authorizer["claims"]; isJwt {
			switch tkn := jwt.(type) {
			case map[string]interface{}:
				*val = mkAccessToken(tkn)
				return nil
			}
		}

		return NoMatch{}
	}
}

// Body decodes HTTP request body to struct
func Body(val interface{}) Endpoint {
	return func(req Input) error {
		p := reflect.ValueOf(val).Elem()
		p.Set(reflect.Zero(p.Type()))

		content, _ := req.Header("Content-Type")
		switch {
		case strings.HasPrefix(content, "application/json"):
			return decodeJSON(req.APIGatewayProxyRequest.Body, &val)
		case strings.HasPrefix(content, "application/x-www-form-urlencoded"):
			return decodeForm(req.APIGatewayProxyRequest.Body, &val)
		}
		return NoMatch{}
	}
}

func decodeJSON(body string, val interface{}) error {
	if err := json.Unmarshal([]byte(body), val); err != nil {
		// TODO: pass error to api client
		return NoMatch{}
	}
	return nil
}

func decodeForm(body string, val interface{}) error {
	buf := bytes.NewBuffer([]byte(body))
	if err := form.NewDecoder(buf).Decode(val); err != nil {
		// TODO: pass error to api client
		return NoMatch{}
	}
	return nil
}

// Text decodes HTTP payload to closed variable
func Text(symbol Symbol) Endpoint {
	return func(req Input) error {
		payload := req.Payload()
		if *payload != "" {
			req.Context().Put(symbol, payload)
			return nil
		}
		return NoMatch{}
	}
}

// FMap applies clojure to matched HTTP request.
// A business logic in gouldian is an endpoint transformation.
func FMap(f func(Context) error) Endpoint {
	return func(req Input) error { return f(req.Context()) }
}
