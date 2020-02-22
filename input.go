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
	"strings"

	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

// Input wraps HTTP request
type Input struct {
	events.APIGatewayProxyRequest
	segment int
	path    []string
	body    string
}

// Mock creates new Input - HTTP GET request
func Mock(httpURL string) *Input {
	return MockVerb("GET", httpURL)
}

// MockVerb creates new Input with any verb
func MockVerb(verb string, httpURL string) *Input {
	uri, _ := url.Parse(httpURL)
	query := map[string]string{}
	for key, val := range uri.Query() {
		query[key] = strings.Join(val, "")
	}

	return NewRequest(
		events.APIGatewayProxyRequest{
			HTTPMethod:            verb,
			Path:                  uri.Path,
			Headers:               map[string]string{},
			QueryStringParameters: query,
		},
	)
}

// NewRequest creates new Input from API Gateway request
func NewRequest(req events.APIGatewayProxyRequest) *Input {
	return &Input{req, 1, strings.Split(req.Path, "/"), ""}
}

// With adds HTTP header to mocked request
func (input *Input) With(head string, value string) *Input {
	input.Headers[head] = value
	return input
}

// WithJSON adds Json payload to mocked request
func (input *Input) WithJSON(val interface{}) *Input {
	body, _ := json.Marshal(val)
	input.Body = string(body)
	return input
}

// WithText adds Text payload to mocked request
func (input *Input) WithText(val string) *Input {
	input.Body = val
	return input
}

// WithAuthorizer adds Authorizer payload to mocked request
func (input *Input) WithAuthorizer(claims map[string]interface{}) *Input {
	input.RequestContext = events.APIGatewayProxyRequestContext{
		Authorizer: map[string]interface{}{
			"claims": claims,
		},
	}
	return input
}
