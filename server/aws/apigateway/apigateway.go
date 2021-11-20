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

package apigateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian"
)

/*

Request is events.APIGatewayProxyRequest ⟼ µ.Input
*/
func Request(r *events.APIGatewayProxyRequest) *µ.Input {
	return &µ.Input{
		Context:  µ.NewContext(context.Background()),
		Method:   r.HTTPMethod,
		Resource: strings.Split(r.Path, "/")[1:],
		Params:   µ.Params(r.MultiValueQueryStringParameters),
		Headers:  µ.Headers(r.MultiValueHeaders),
		Payload:  r.Body,
	}
}

// Serve HTTP service
func Serve(endpoints ...µ.Endpoint) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := µ.Or(endpoints...)

	return func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req := Request(&r)
		switch v := api(req).(type) {
		case *µ.Output:
			return output(v, req)
		case µ.NoMatch:
			failure := µ.Status.NotImplemented(
				µ.WithIssue(fmt.Errorf("NoMatch %s", r.Path)),
			).(*µ.Output)
			return output(failure, req)
		default:
			failure := µ.Status.InternalServerError(
				µ.WithIssue(fmt.Errorf("Unknown response %s", r.Path)),
			).(*µ.Output)
			return output(failure, req)
		}
	}
}

func output(out *µ.Output, req *µ.Input) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       out.Body,
		StatusCode: out.Status,
		Headers:    joinHead(defaultCORS(req), out.Headers),
	}, nil
}

func defaultCORS(req *µ.Input) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":  defaultOrigin(req),
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Max-Age":       "600",
	}
}

func defaultOrigin(req *µ.Input) string {
	origin, exists := req.Headers.Get("Origin")
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

// JWT decodes token associated with the request.
// Endpoint fails if Authentication context is not found in the request.
/*
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
*/

// TODO: Gone
// Input wraps HTTP request
/*
type Input struct {
	events.APIGatewayProxyRequest
	Path []string
	Body string
}
*/

// Request creates new Input from API Gateway request
/*
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
*/
