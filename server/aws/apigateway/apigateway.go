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
	"net/url"
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
		Resource: splitPath(r.Path),
		Params:   µ.Params(r.MultiValueQueryStringParameters),
		Headers:  µ.Headers(r.MultiValueHeaders),
		JWT:      jwtFromAuthorizer(r),
		Payload:  r.Body,
	}
}

func splitPath(path string) µ.Segments {
	seq := strings.Split(path, "/")[1:]
	segments := make(µ.Segments, 0, len(seq))
	for _, x := range seq {
		if val, err := url.PathUnescape(x); err != nil {
			segments = append(segments, x)
		} else {
			segments = append(segments, val)
		}
	}

	if len(segments) == 1 && segments[0] == "" {
		segments = segments[:0]
	}

	return segments
}

func jwtFromAuthorizer(r *events.APIGatewayProxyRequest) µ.JWT {
	if r.RequestContext.Authorizer == nil {
		return nil
	}

	if jwt, isJwt := r.RequestContext.Authorizer["claims"]; isJwt {
		switch tkn := jwt.(type) {
		case map[string]interface{}:
			return µ.NewJWT(tkn)
		}
	}

	return nil

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
