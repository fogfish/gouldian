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
	"log"

	"github.com/aws/aws-lambda-go/events"
)

// Serve HTTP service
func Serve(seq ...Endpoint) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := Or(seq...)

	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		http := Request(req)
		switch v := api(http).(type) {
		case *Output:
			return events.APIGatewayProxyResponse{
				Body:       v.Body,
				StatusCode: v.Status,
				Headers:    joinHead(defaultCORS(http), v.Headers),
			}, nil
		case *Issue:
			return recoverIssue(http, v), nil
		case NoMatch:
			iss := NotImplemented(fmt.Errorf("NoMatch %v", http.APIGatewayProxyRequest.Path))
			// TODO: output JSON only in debug mode
			// b, _ := json.Marshal(req)
			// log.Printf("ERROR %v, NoMatch: %v", iss.ID, string(b))
			return recoverIssue(http, iss), nil
		default:
			iss := InternalServerError(v)
			return recoverIssue(http, iss), nil
		}
	}
}

func recoverIssue(http *Input, issue *Issue) events.APIGatewayProxyResponse {
	log.Printf("ERROR %v: %d %s, %v", issue.ID, issue.Status, issue.Title, issue.Failure)
	text, _ := json.Marshal(issue)
	return events.APIGatewayProxyResponse{
		Body:       string(text),
		StatusCode: issue.Status,
		Headers:    joinHead(defaultCORS(http), map[string]string{"Content-Type": "application/json"}),
	}
}

func defaultCORS(req *Input) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":  defaultOrigin(req),
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Max-Age":       "600",
	}
}

func defaultOrigin(req *Input) string {
	origin, exists := req.Header("Origin")
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
