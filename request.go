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
	"fmt"
	"strconv"
	"strings"
)

// HTTPPath defines Endpoint(s) to match path elements of HTTP request
type HTTPPath interface {
	Path(segment string) HTTP
	String(val *string) HTTP
	Int(val *int) HTTP
}

// HTTPQuery defines Endpoint(s) to match query elements of HTTP request
type HTTPQuery interface {
	Param(name string, val string) HTTP
	HasParam(name string) HTTP
	QString(name string, val *string) HTTP
	QInt(name string, val *int) HTTP
}

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

// HTTP defines Endpoint(s) to match elements of HTTP request
type HTTP interface {
	HTTPPath
	HTTPQuery
	HTTPHeader
	HTTPBody

	FMap(f func() error) Endpoint
	IsMatch(req *Input) bool
}

// NoMatch is returned by Endpoint if Input is not matched.
type NoMatch struct{}

func (err NoMatch) Error() string {
	return fmt.Sprintf("No Match")
}

// APIGateway implements Endpoints to process AWS API Gateway request(s).
// There is a type constructor named after HTTP vers. It creates
// Endpoint to match HTTP verbs (methods).
type APIGateway struct {
	f Endpoint
}

// Delete is Endpoint to match HTTP DELETE verb.
func Delete() HTTP {
	return &APIGateway{isVerb("DELETE")}
}

// Get is Endpoint to match HTTP GET verb.
func Get() HTTP {
	return &APIGateway{isVerb("GET")}
}

// Patch is Endpoint to match HTTP PATCH verb.
func Patch() HTTP {
	return &APIGateway{isVerb("PATCH")}
}

// Post is Endpoint to match HTTP POST verb.
func Post() HTTP {
	return &APIGateway{isVerb("POST")}
}

// Put is Endpoint to match HTTP PUT verb.
func Put() HTTP {
	return &APIGateway{isVerb("PUT")}
}

func isVerb(verb string) Endpoint {
	return func(http *Input) error {
		if http.HTTPMethod == verb {
			http.segment = 1
			return nil
		}
		return NoMatch{}
	}
}

func hasSegment(req *Input) error {
	if len(req.path) > req.segment {
		return nil
	}
	return NoMatch{}
}

// Path matches the current path segment
//   e := gouldian.Get().Path("foo")
//   e.IsMatch(gouldian.New("/foo")) == true
//   e.IsMatch(gouldian.New("/bar")) == false
func (state *APIGateway) Path(segment string) HTTP {
	state.f = state.f.Then(hasSegment).Then(func(req *Input) error {
		if req.path[req.segment] == segment {
			req.segment++
			return nil
		}
		return NoMatch{}
	})
	return state
}

// String matches the current path segment to string type,
// matched segment is returned to closure
//   var value string
//   e := gouldian.Get().String(&value)
//   e.IsMatch(gouldian.New("/foo")) == true && value == "foo"
//   e.IsMatch(gouldian.New("/1")) == true && value == "1"
func (state *APIGateway) String(val *string) HTTP {
	state.f = state.f.Then(hasSegment).Then(func(req *Input) error {
		*val = req.path[req.segment]
		req.segment++
		return nil
	})
	return state
}

// Int matches the current path segment to int type,
// matched segment is returned to closure.
// Endpoint fails if value cannot be converted to int
//   var value int
//   e := gouldian.Get().String(&value)
//   e.IsMatch(gouldian.New("/1")) == true && value == 1
//   e.IsMatch(gouldian.New("/foo")) == false
func (state *APIGateway) Int(val *int) HTTP {
	state.f = state.f.Then(hasSegment).Then(func(req *Input) error {
		value, err := strconv.Atoi(req.path[req.segment])
		if err != nil {
			return NoMatch{}
		}
		*val = value
		req.segment++
		return nil
	})
	return state
}

// Param checks if query param present in URL query string and
// it value equals defined one.
func (state *APIGateway) Param(name string, val string) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		opt, exists := req.QueryStringParameters[name]
		if exists && opt == val {
			return nil
		}
		return NoMatch{}
	})
	return state
}

// HasParam check presence of URL query parameter
func (state *APIGateway) HasParam(name string) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		_, exists := req.QueryStringParameters[name]
		if exists {
			return nil
		}
		return NoMatch{}
	})
	return state
}

// QString matches parameter and lifts its value
func (state *APIGateway) QString(name string, val *string) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		opt, exists := req.QueryStringParameters[name]
		if exists {
			*val = opt
			return nil
		}
		*val = ""
		return nil
	})
	return state
}

// QInt matches parameter to int type and lifts its value
func (state *APIGateway) QInt(name string, val *int) HTTP {
	state.f = state.f.Then(func(req *Input) error {
		opt, exists := req.QueryStringParameters[name]
		if exists {
			value, err := strconv.Atoi(opt)
			if err != nil {
				return NoMatch{}
			}
			*val = value
			return nil
		}
		*val = 0
		return nil
	})
	return state
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

// IsMatch evaluates Endpoint against mocked Input
func (state *APIGateway) IsMatch(in *Input) bool {
	return state.f(in) == nil
}

// FMap applies clojure to matched HTTP request.
// A business logic in gouldian is an endpoint transformation.
func (state *APIGateway) FMap(f func() error) Endpoint {
	return state.f.Then(func(req *Input) error { return f() })
}
