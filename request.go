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
	"encoding/json"

	"github.com/fogfish/gouldian/core"
	"github.com/fogfish/gouldian/header"
	"github.com/fogfish/gouldian/param"
	"github.com/fogfish/gouldian/path"
)

// DELETE composes product Endpoint match HTTP DELETE request.
//   e := µ.DELETE()
//   e(mock.Input(mock.Method("DELETE"))) == nil
//   e(mock.Input(mock.Method("OTHER"))) != nil
func DELETE(arrows ...core.Endpoint) core.Endpoint {
	return Method("DELETE").Then(core.Join(arrows...))
}

// GET composes product Endpoint match HTTP GET request.
//   e := µ.GET()
//   e(mock.Input(mock.Method("GET"))) == nil
//   e(mock.Input(mock.Method("OTHER"))) != nil
func GET(arrows ...core.Endpoint) core.Endpoint {
	return Method("GET").Then(core.Join(arrows...))
}

// PATCH composes product Endpoint match HTTP PATCH request.
//   e := µ.PATCH()
//   e(mock.Input(mock.Method("PATCH"))) == nil
//   e(mock.Input(mock.Method("OTHER"))) != nil
func PATCH(arrows ...core.Endpoint) core.Endpoint {
	return Method("PATCH").Then(core.Join(arrows...))
}

// POST composes product Endpoint match HTTP POST request.
//   e := µ.POST()
//   e(mock.Input(mock.Method("POST"))) == nil
//   e(mock.Input(mock.Method("OTHER"))) != nil
func POST(arrows ...core.Endpoint) core.Endpoint {
	return Method("POST").Then(core.Join(arrows...))
}

// PUT composes product Endpoint match HTTP PUT request.
//   e := µ.PUT()
//   e(mock.Input(mock.Method("PUT"))) == nil
//   e(mock.Input(mock.Method("OTHER"))) != nil
func PUT(arrows ...core.Endpoint) core.Endpoint {
	return Method("PUT").Then(core.Join(arrows...))
}

// ANY composes product Endpoint match HTTP PUT request.
//   e := µ.ANY()
//   e(mock.Input(mock.Method("PUT"))) == nil
//   e(mock.Input(mock.Method("OTHER"))) == nil
func ANY(arrows ...core.Endpoint) core.Endpoint {
	return Method("*").Then(core.Join(arrows...))
}

// Method is an endpoint to match HTTP verb request
func Method(verb string) core.Endpoint {
	if verb == "*" {
		return func(http *core.Input) error { return nil }
	}

	return func(http *core.Input) error {
		if http.HTTPMethod == verb {
			return nil
		}
		return core.NoMatch{}
	}
}

// Path is an endpoint to match URL of HTTP request. The function takes
// url matching primitives, which are defined by the package `path`.
//
//   import "github.com/fogfish/gouldian/path"
//
//   e := µ.GET( µ.Path(path.Is("foo")) )
//   e(mock.Input(mock.URL("/foo"))) == nil
//   e(mock.Input(mock.URL("/bar"))) != nil
func Path(arrows ...path.Arrow) core.Endpoint {
	return func(req *core.Input) error {
		for i, f := range arrows {
			if i > len(req.Path)-1 {
				return core.NoMatch{}
			}
			if err := f(req.Path[i]); err != nil {
				return err
			}
		}
		return nil
	}
}

// Param is an endpoint to match URL Query parameters of HTTP request.
// The function takes url query matching primitives, which are defined
// by the package `param`.
//
//   import "github.com/fogfish/gouldian/param"
//
//   e := µ.GET( µ.Param(param.Is("foo", "bar")) )
//   e(mock.Input(mock.URL("/?foo=bar"))) == nil
//   e(mock.Input(mock.URL("/?foo=baz"))) != nil
func Param(arrows ...param.Arrow) core.Endpoint {
	return func(req *core.Input) error {
		for _, f := range arrows {
			if err := f(req.APIGatewayProxyRequest.QueryStringParameters); err != nil {
				return err
			}
		}
		return nil
	}
}

// Header is an endpoint to match Header(s) of HTTP request.
// The function takes header matching primitives, which are defined
// by the package `header`.
//
//   import "github.com/fogfish/gouldian/header"
//
//   e := µ.GET(
//     µ.Header(
//       param.Header("Content-Type", "application/json"),
//     ),
//   )
//   Json := mock.Header("Content-Type", "application/json")
//   e(mock.Input(Json)) == nil
//
//   Text := mock.Header("Content-Type", "text/plain")
//   e(mock.Input(Text)) != nil
func Header(arrows ...header.Arrow) core.Endpoint {
	return func(req *core.Input) error {
		for _, f := range arrows {
			if err := f(req.APIGatewayProxyRequest.Headers); err != nil {
				return err
			}
		}
		return nil
	}
}

// AccessToken decodes JWT token associated with the request.
// Endpoint fails if Authentication context is not found in the request.
func AccessToken(val *core.AccessToken) core.Endpoint {
	return func(req *core.Input) error {
		if req.RequestContext.Authorizer == nil {
			return core.NoMatch{}
		}

		if jwt, isJwt := req.RequestContext.Authorizer["claims"]; isJwt {
			switch tkn := jwt.(type) {
			case map[string]interface{}:
				*val = core.AccessToken{
					Sub:   tkn["sub"].(string),
					Scope: tkn["scope"].(string),
				}
				return nil
			}
		}

		return core.NoMatch{}
	}
}

// JSON decodes HTTP payload to struct
func JSON(val interface{}) core.Endpoint {
	return func(req *core.Input) error {
		err := json.Unmarshal([]byte(req.Body), val)
		if err == nil {
			return nil
		}
		// TODO: pass error to api client
		return core.NoMatch{}
	}
}

// Text decodes HTTP payload to closed variable
func Text(val *string) core.Endpoint {
	return func(req *core.Input) error {
		*val = ""
		if req.Body != "" {
			*val = req.Body
			return nil
		}
		return core.NoMatch{}
	}
}

// FMap applies clojure to matched HTTP request.
// A business logic in gouldian is an endpoint transformation.
func FMap(f func() error) core.Endpoint {
	return func(*core.Input) error { return f() }
}
