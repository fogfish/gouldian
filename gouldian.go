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
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

// Serve HTTP service
func Serve(seq ...Endpoint) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := Or(seq...)

	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var output *Output
		var issue *Issue

		http := Request(req)
		err := api(http)
		if errors.As(err, &output) {
			return events.APIGatewayProxyResponse{
				Body:       output.Body,
				StatusCode: output.Status,
				Headers:    joinHead(defaultCORS(req), output.Headers),
			}, nil
		} else if errors.As(err, &issue) {
			text, _ := json.Marshal(issue)
			return events.APIGatewayProxyResponse{
				Body:       string(text),
				StatusCode: issue.Status,
				Headers:    joinHead(defaultCORS(req), map[string]string{"Content-Type": "application/json"}),
			}, nil
		} else if errors.Is(err, NoMatch{}) {
			return events.APIGatewayProxyResponse{
				StatusCode: 501,
				Headers:    defaultCORS(req),
			}, nil
		}
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
}

func defaultCORS(req events.APIGatewayProxyRequest) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":  defaultOrigin(req),
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Max-Age":       "600",
	}
}

func defaultOrigin(req events.APIGatewayProxyRequest) string {
	origin, exists := req.Headers["Origin"]
	if exists {
		return origin
	}
	return "*"
}

func joinHead(a, b map[string]string) map[string]string {
	for keyA, valA := range a {
		if _, ok := b[keyA]; !ok {
			b[keyA] = valA
		}
	}
	return b
}

// NoMatchLogger logs the request
func NoMatchLogger() Endpoint {
	return func(req *Input) error {
		bytes, _ := json.Marshal(req)
		log.Printf("No Match:\n%v\n", string(bytes))
		return NoMatch{}
	}
}
