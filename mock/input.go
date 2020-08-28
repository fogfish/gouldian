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

package mock

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian"
)

// Mock is an option type to customize mock event
type Mock func(*µ.Input) *µ.Input

// Input mocks HTTP event, takes mock options to customize event
func Input(spec ...Mock) *µ.Input {
	input := µ.Request(
		events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
			Path:       "/",
			Headers:    map[string]string{},
		},
	)
	for _, f := range spec {
		input = f(input)
	}
	return input
}

// Method changes the verb of mocked HTTP request
func Method(verb string) Mock {
	return func(mock *µ.Input) *µ.Input {
		mock.APIGatewayProxyRequest.HTTPMethod = verb
		return mock
	}
}

// URL changes URL of mocked HTTP request
func URL(httpURL string) Mock {
	return func(mock *µ.Input) *µ.Input {
		uri, _ := url.Parse(httpURL)
		query := map[string]string{}
		for key, val := range uri.Query() {
			query[key] = strings.Join(val, "")
		}
		mock.APIGatewayProxyRequest.Path = uri.Path
		mock.APIGatewayProxyRequest.QueryStringParameters = query
		segments := strings.Split(uri.Path, "/")[1:]
		if len(segments) == 1 && segments[0] == "" {
			segments = []string{}
		}
		mock.Path = segments
		return mock
	}
}

// Param add raw param string to mocked HTTP request
func Param(key, val string) Mock {
	return func(mock *µ.Input) *µ.Input {
		if mock.APIGatewayProxyRequest.QueryStringParameters == nil {
			mock.APIGatewayProxyRequest.Path = "/"
			mock.APIGatewayProxyRequest.QueryStringParameters = map[string]string{}
		}
		mock.APIGatewayProxyRequest.QueryStringParameters[key] = val
		return mock
	}
}

// Header adds Header to mocked HTTP request
func Header(header string, value string) Mock {
	return func(mock *µ.Input) *µ.Input {
		mock.APIGatewayProxyRequest.Headers[header] = value
		return mock
	}
}

// JSON adds payload to mocked HTTP request
func JSON(val interface{}) Mock {
	return func(mock *µ.Input) *µ.Input {
		body, _ := json.Marshal(val)
		mock.APIGatewayProxyRequest.Headers["Content-Type"] = "application/json"
		mock.APIGatewayProxyRequest.Body = string(body)
		return mock
	}
}

// Text adds payload to mocked HTTP request
func Text(val string) Mock {
	return func(mock *µ.Input) *µ.Input {
		mock.APIGatewayProxyRequest.Body = val
		return mock
	}
}

// Auth adds Authorizer payload to mocked HTTP request
func Auth(claims map[string]interface{}) Mock {
	return func(mock *µ.Input) *µ.Input {
		mock.RequestContext = events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": claims,
			},
		}
		return mock
	}
}
