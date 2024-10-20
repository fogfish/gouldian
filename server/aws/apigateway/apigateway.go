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
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian/v2"
	ø "github.com/fogfish/gouldian/v2/output"
	"github.com/fogfish/logger"
)

// Request is events.APIGatewayProxyRequest ⟼ µ.Input
func Request(r *events.APIGatewayProxyRequest) *µ.Context {
	ctx := µ.NewContext(context.Background())
	req, err := http.NewRequest(r.HTTPMethod, r.Path, requestBody(r))
	if err != nil {
		return nil
	}

	for header, value := range r.Headers {
		req.Header.Set(header, value)
	}

	q := req.URL.Query()
	for key, val := range r.QueryStringParameters {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	ctx.Request = req
	ctx.JWT = jwtFromAuthorizer(r)

	return ctx
}

func requestBody(r *events.APIGatewayProxyRequest) io.ReadCloser {
	reader := strings.NewReader(r.Body)

	if r.IsBase64Encoded {
		return io.NopCloser(
			base64.NewDecoder(base64.StdEncoding, reader),
		)
	}

	return io.NopCloser(reader)
}

func jwtFromAuthorizer(r *events.APIGatewayProxyRequest) µ.Token {
	if r.RequestContext.Authorizer != nil {
		if jwt, isJwt := r.RequestContext.Authorizer["claims"]; isJwt {
			switch tkn := jwt.(type) {
			case map[string]interface{}:
				return µ.NewToken(tkn)
			}
		}

		return nil
	}

	if r.RequestContext.Identity.UserArn != "" {
		return µ.Token{
			"iss":      "https://aws.amazon.com/iam",
			"sub":      r.RequestContext.Identity.User,
			"username": r.RequestContext.Identity.UserArn,
		}
	}

	return nil

}

func ServeAndCommit(commit func(), endpoints ...µ.Routable) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h := Serve(endpoints...)

	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		ret, err := h(req)
		commit()
		return ret, err
	}
}

// Serve HTTP service
func Serve(endpoints ...µ.Routable) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := µ.NewRoutes(endpoints...).Endpoint()

	return func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req := Request(&r)
		if req == nil {
			failure := ø.Status.BadRequest(
				ø.Error(fmt.Errorf("unknown response %s", r.Path)),
			).(*µ.Output)
			return output(failure, req)
		}

		switch v := api(req).(type) {
		case *µ.Output:
			return output(v, req)
		case µ.NoMatch:
			failure := ø.Status.NotImplemented(
				ø.Error(fmt.Errorf("NoMatch %s", r.Path)),
			).(*µ.Output)
			return output(failure, req)
		default:
			failure := ø.Status.InternalServerError(
				ø.Error(fmt.Errorf("unknown response %s", r.Path)),
			).(*µ.Output)
			return output(failure, req)
		}
	}
}

func output(out *µ.Output, req *µ.Context) (events.APIGatewayProxyResponse, error) {
	if out.Failure != nil && req != nil && req.Request != nil {
		logger.Error("%s %v", req.Request.URL, out.Failure)
	}

	evt := events.APIGatewayProxyResponse{
		Body:       out.Body,
		StatusCode: out.Status,
		Headers:    joinHead(defaultCORS(req), out.Headers),
	}
	out.Free()

	return evt, nil
}

func defaultCORS(req *µ.Context) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":  defaultOrigin(req),
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Max-Age":       "600",
	}
}

func defaultOrigin(req *µ.Context) string {
	if req == nil {
		return "*"
	}

	origin := req.Request.Header.Get("Origin")
	if origin != "" {
		return origin
	}
	return "*"
}

func joinHead(a map[string]string, b []struct{ Header, Value string }) map[string]string {
	for _, v := range b {
		if _, ok := a[v.Header]; !ok {
			a[v.Header] = v.Value
		}
	}
	return a
}
