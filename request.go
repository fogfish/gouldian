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
	"github.com/fogfish/gouldian/core"
	"github.com/fogfish/gouldian/param"
	"github.com/fogfish/gouldian/path"
)

// APIGateway implements Endpoints to process AWS API Gateway request(s).
// There is a type constructor named after HTTP vers. It creates
// Endpoint to match HTTP verbs (methods).
type APIGateway struct {
	f core.Endpoint
}

// DELETE composes product Endpoint match HTTP DELETE request.
//   e := µ.DELETE()
//   e.IsMatch(mock.Input(mock.Method("DELETE"))) == true
//   e.IsMatch(mock.Input(mock.Method("OTHER"))) == false
func DELETE(arrows ...core.Endpoint) *APIGateway {
	return (&APIGateway{Method("DELETE")}).Do(arrows...)
}

// GET composes product Endpoint match HTTP GET request.
//   e := µ.GET()
//   e.IsMatch(mock.Input(mock.Method("GET"))) == true
//   e.IsMatch(mock.Input(mock.Method("OTHER"))) == false
func GET(arrows ...core.Endpoint) *APIGateway {
	return (&APIGateway{Method("GET")}).Do(arrows...)
}

// PATCH composes product Endpoint match HTTP PATCH request.
//   e := µ.PATCH()
//   e.IsMatch(mock.Input(mock.Method("PATCH"))) == true
//   e.IsMatch(mock.Input(mock.Method("OTHER"))) == false
func PATCH(arrows ...core.Endpoint) *APIGateway {
	return (&APIGateway{Method("PATCH")}).Do(arrows...)
}

// POST composes product Endpoint match HTTP POST request.
//   e := µ.POST()
//   e.IsMatch(mock.Input(mock.Method("POST"))) == true
//   e.IsMatch(mock.Input(mock.Method("OTHER"))) == false
func POST(arrows ...core.Endpoint) *APIGateway {
	return (&APIGateway{Method("POST")}).Do(arrows...)
}

// PUT composes product Endpoint match HTTP PUT request.
//   e := µ.PUT()
//   e.IsMatch(mock.Input(mock.Method("PUT"))) == true
//   e.IsMatch(mock.Input(mock.Method("OTHER"))) == false
func PUT(arrows ...core.Endpoint) *APIGateway {
	return (&APIGateway{Method("PUT")}).Do(arrows...)
}

// Method is an endpoint to match HTTP verb request
func Method(verb string) core.Endpoint {
	return func(http *core.Input) error {
		if http.HTTPMethod == verb {
			return nil
		}
		return core.NoMatch{}
	}
}

// Path is an endpoint to match URL of HTTP request. The function takes
// url matching primitives, which are define by the package `path`.
//
//   import "github.com/fogfish/gouldian/path"
//
//   e := µ.GET( µ.Path(path.Is("foo")) )
//   e.IsMatch(mock.Input(mock.URL("/foo"))) == true
//   e.IsMatch(mock.Input(mock.URL("/bar"))) == false
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
// The function takes url query matching primitives, which are define
// by the package `param`.
//
//   import "github.com/fogfish/gouldian/param"
//
//   e := µ.GET( µ.Param(param.Is("foo", "bar")) )
//   e.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == true
//   e.IsMatch(mock.Input(mock.URL("/?foo=baz"))) == false
func Param(arrows ...param.Arrow) core.Endpoint {
	return func(req *core.Input) error {
		for i, f := range arrows {
			if i > len(req.Path)-1 {
				return core.NoMatch{}
			}
			if err := f(req.APIGatewayProxyRequest.QueryStringParameters); err != nil {
				return err
			}
		}
		return nil
	}
}

// Do
func (state *APIGateway) Do(arrows ...core.Endpoint) *APIGateway {
	for _, f := range arrows {
		state.f = state.f.Then(f)
	}
	return state
}

// IsMatch evaluates Endpoint against mocked Input
func (state *APIGateway) IsMatch(in *core.Input) bool {
	return state.f(in) == nil
}

// FMap applies clojure to matched HTTP request.
// A business logic in gouldian is an endpoint transformation.
func (state *APIGateway) FMap(f func() error) core.Endpoint {
	return state.f.Then(func(req *core.Input) error { return f() })
}

/*
// HTTPHeader defines Endpoint(s) to match headers of HTTP request
type HTTPHeader interface {
	Head(name string, val string) HTTP
	HasHead(name string) HTTP
	HString(name string, val *string) HTTP
}

// HTTPBody defines Endpoint(s) to match body of HTTP Request
type HTTPBody interface {
	JSON(val interface{}) HTTP
	Text(val *string) HTTP
}

// HTTPAuthorize defines Endpoint(s) to match Access Token
type HTTPAuthorize interface {
	AccessToken(token *AccessToken) HTTP
}

// AccessToken is a container for user identity
type AccessToken struct {
	Sub   string
	Scope string
}

// HTTP defines Endpoint(s) to match elements of HTTP request
type HTTP interface {
	HTTPPath
	HTTPQuery
	HTTPHeader
	HTTPBody
	HTTPAuthorize

	FMap(f func() error) Endpoint
	IsMatch(req *Input) bool
}

// Head checks if HTTP header exists in request and
// it value has defined prefix.
func (state *APIGateway) Head(name string, val string) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		head, exists := req.Headers[name]
		if exists && strings.HasPrefix(head, val) {
			return nil
		}
		return NoMatch{}
	})
	return state
}

// HasHead checks if HTTP header exists in the request
func (state *APIGateway) HasHead(name string) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		_, exists := req.Headers[name]
		if exists {
			return nil
		}
		return NoMatch{}
	})
	return state
}

// HString matches HTTP header and lifts its value
func (state *APIGateway) HString(name string, val *string) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		head, exists := req.Headers[name]
		if exists {
			*val = head
			return nil
		}
		*val = ""
		return nil
	})
	return state
}

// JSON decodes HTTP payload to struct
func (state *APIGateway) JSON(val interface{}) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		err := json.Unmarshal([]byte(req.Body), val)
		if err == nil {
			return nil
		}
		// TODO: pass error to api client
		return NoMatch{}
	})
	return state
}

// Text decodes HTTP payload to text
func (state *APIGateway) Text(val *string) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		*val = req.Body
		return nil
	})
	return state
}

// AccessToken decodes JWT token associated with the request
func (state *APIGateway) AccessToken(val *AccessToken) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		if req.RequestContext.Authorizer != nil {
			if jwt, isJwt := req.RequestContext.Authorizer["claims"]; isJwt {
				switch tkn := jwt.(type) {
				case map[string]interface{}:
					*val = AccessToken{
						Sub:   tkn["sub"].(string),
						Scope: tkn["scope"].(string),
					}
					return nil
				}
			}
		}
		return NoMatch{}
	})
	return state
}


*/
