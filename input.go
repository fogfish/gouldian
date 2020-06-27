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
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// Input wraps HTTP request
type Input struct {
	events.APIGatewayProxyRequest
	Path []string
	Body string
}

// Request creates new Input from API Gateway request
func Request(req events.APIGatewayProxyRequest) *Input {
	segments := []string{}
	for _, x := range strings.Split(req.Path, "/")[1:] {
		if val, err := url.PathUnescape(x); err != nil {
			segments = append(segments, x)
		} else {
			segments = append(segments, val)
		}
	}

	if len(segments) == 1 && segments[0] == "" {
		segments = []string{}
	}

	return &Input{req, segments, ""}
}

// Header returns header value
func (req *Input) Header(key string) (string, bool) {
	v, exists := req.APIGatewayProxyRequest.Headers[key]
	if !exists {
		v, exists = req.APIGatewayProxyRequest.Headers[strings.ToLower(key)]
		return v, exists
	}
	return v, exists
}
