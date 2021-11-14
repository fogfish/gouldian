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

package mock

import (
	"context"
	"encoding/json"
	"net/textproto"
	"net/url"
	"strings"

	µ "github.com/fogfish/gouldian"
)

/*

µMock is abstract container of HTTP terms for testing purposes
*/
type µMock struct {
	ctx      µ.Context
	method   string
	resource µ.Segments
	params   µ.Params
	headers  µ.Headers
	payload  []byte
}

var _ µ.Input = (*µMock)(nil)

func (mock *µMock) Context() µ.Context { return mock.ctx }

func (mock *µMock) Method() string { return mock.method }

func (mock *µMock) Resource() µ.Segments { return mock.resource }

func (mock *µMock) Params() µ.Params { return mock.params }

func (mock *µMock) Headers() µ.Headers { return mock.headers }

func (mock *µMock) Payload() []byte { return mock.payload }

// Mock is an option type to customize mock event
type Mock func(*µMock) *µMock

// Input mocks HTTP request, takes mock options to customize event
func Input(spec ...Mock) µ.Input {
	input := &µMock{
		ctx:      µ.NewContext(context.Background()),
		method:   "GET",
		resource: µ.Segments{},
		params:   µ.Params{},
		headers:  µ.Headers{},
		payload:  nil,
	}

	for _, f := range spec {
		input = f(input)
	}
	return input
}

// Method changes the verb of mocked HTTP request
func Method(verb string) Mock {
	return func(mock *µMock) *µMock {
		mock.method = verb
		return mock
	}
}

// URL changes URL of mocked HTTP request
func URL(httpURL string) Mock {
	return func(mock *µMock) *µMock {
		uri, err := url.Parse(httpURL)
		if err != nil {
			panic(err)
		}

		segments := strings.Split(uri.Path, "/")[1:]
		if len(segments) == 1 && segments[0] == "" {
			segments = []string{}
		}
		mock.resource = segments

		params := µ.Params{}
		for key, val := range uri.Query() {
			params[key] = []string{strings.Join(val, "")}
		}
		mock.params = params

		return mock
	}
}

// Param add raw param string to mocked HTTP request
func Param(key, val string) Mock {
	return func(mock *µMock) *µMock {
		mock.params[key] = []string{val}
		return mock
	}
}

// Header adds Header to mocked HTTP request
func Header(header string, value string) Mock {
	return func(mock *µMock) *µMock {
		head := textproto.CanonicalMIMEHeaderKey(header)
		mock.headers[head] = value
		return mock
	}
}

// JSON adds payload to mocked HTTP request
func JSON(val interface{}) Mock {
	return func(mock *µMock) *µMock {
		body, err := json.Marshal(val)
		if err != nil {
			panic(err)
		}
		mock.headers["Content-Type"] = "application/json"
		mock.payload = body
		return mock
	}
}

// Text adds payload to mocked HTTP request
func Text(val string) Mock {
	return func(mock *µMock) *µMock {
		mock.payload = []byte(val)
		return mock
	}
}

// Auth adds Authorizer payload to mocked HTTP request
/*

TODO

func Auth(token µ.AccessToken) Mock {
	return func(mock *µ.Input) *µ.Input {
		bin, err := json.Marshal(token)
		if err != nil {
			panic(err)
		}
		var claims map[string]interface{}
		if err := json.Unmarshal(bin, &claims); err != nil {
			panic(err)
		}

		mock.RequestContext = events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": claims,
			},
		}
		return mock
	}
}
*/
