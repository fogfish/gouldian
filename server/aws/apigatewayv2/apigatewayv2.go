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

// Package apigatewayv2 adapts the gouldian server to AWS API Gateway HTTP API
// and Lambda Function URL (payload format 2.0).
package apigatewayv2

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/internal/media"
)

// Func is AWS Lambda handler
type Func = func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

/*
Serve builds Lambda handler for sequence of endpoints

	lambda.Start(apigatewayv2.Serve( ... ))
*/
func Serve(endpoints ...µ.Routable) Func {
	return Handler(µ.NewServer().Mount(endpoints...))
}

/*
Handler adapts the server to AWS Lambda (HTTP API, Function URL). The context
of the invocation (the deadline, the cancellation) is propagated to endpoints.

	api := µ.NewServer(µ.WithCORS(...)).Mount( ... )
	lambda.Start(apigatewayv2.Handler(api))
*/
func Handler(srv *µ.Server) Func {
	return func(ctx context.Context, r events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		req, err := Request(ctx, &r)
		if err != nil {
			return badRequest(r.RawPath, err), nil
		}

		c := srv.NewContext(ctx, req)
		c.RequestID = r.RequestContext.RequestID
		c.JWT = tokenFromAuthorizer(&r)

		out := srv.Dispatch(c)
		evt := Response(out)
		out.Free()
		srv.Release(c)

		return evt, nil
	}
}

// Request is events.APIGatewayV2HTTPRequest ⟼ http.Request
func Request(ctx context.Context, r *events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	body := []byte(r.Body)
	if r.IsBase64Encoded {
		buf, err := base64.StdEncoding.DecodeString(r.Body)
		if err != nil {
			return nil, fmt.Errorf("body is not valid base64: %w", err)
		}
		body = buf
	}

	method := r.RequestContext.HTTP.Method
	req, err := http.NewRequestWithContext(ctx, method, "/", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	path := r.RawPath
	if path == "" {
		path = r.RequestContext.HTTP.Path
	}
	// RawPath is the escaped form
	if unescaped, err := url.PathUnescape(path); err == nil {
		req.URL = &url.URL{Path: unescaped, RawPath: path, RawQuery: r.RawQueryString}
	} else {
		req.URL = &url.URL{Path: path, RawQuery: r.RawQueryString}
	}

	// HTTP API joins multiple values of header with comma
	for header, value := range r.Headers {
		req.Header.Set(header, value)
	}
	if len(r.Cookies) != 0 {
		req.Header.Set("Cookie", strings.Join(r.Cookies, "; "))
	}
	req.Host = r.RequestContext.DomainName

	return req, nil
}

func tokenFromAuthorizer(r *events.APIGatewayV2HTTPRequest) µ.Token {
	auth := r.RequestContext.Authorizer
	if auth == nil {
		return nil
	}

	switch {
	case auth.JWT != nil:
		token := make(µ.Token, len(auth.JWT.Claims)+1)
		for k, v := range auth.JWT.Claims {
			token[k] = v
		}
		if _, has := token["scope"]; !has && len(auth.JWT.Scopes) != 0 {
			token["scope"] = strings.Join(auth.JWT.Scopes, " ")
		}
		return token

	case auth.IAM != nil:
		return µ.Token{
			"iss":      "https://aws.amazon.com/iam",
			"sub":      auth.IAM.UserID,
			"username": auth.IAM.UserARN,
		}

	case auth.Lambda != nil:
		return µ.NewToken(auth.Lambda)
	}

	return nil
}

// Response is µ.Output ⟼ events.APIGatewayV2HTTPResponse. Set-Cookie headers
// are sent as cookies, headers with multiple values are joined, binary bodies
// are base64 encoded.
func Response(out *µ.Output) events.APIGatewayV2HTTPResponse {
	evt := events.APIGatewayV2HTTPResponse{
		StatusCode: out.Status,
	}

	var contentType string
	for _, h := range out.Headers {
		switch h.Header {
		case "Set-Cookie":
			evt.Cookies = append(evt.Cookies, h.Value)
			continue
		case "Content-Type":
			contentType = h.Value
		}

		if evt.Headers == nil {
			evt.Headers = make(map[string]string, len(out.Headers))
		}
		if prev, exists := evt.Headers[h.Header]; exists {
			evt.Headers[h.Header] = prev + "," + h.Value
		} else {
			evt.Headers[h.Header] = h.Value
		}
	}

	if media.IsBinary(contentType, out.Body) {
		evt.IsBase64Encoded = true
		evt.Body = base64.StdEncoding.EncodeToString(out.Body)
	} else {
		evt.Body = string(out.Body)
	}

	return evt
}

func badRequest(path string, err error) events.APIGatewayV2HTTPResponse {
	issue := µ.NewIssue(http.StatusBadRequest)
	issue.Detail = err.Error()
	issue.Instance = path
	ct, body := µ.DefaultIssue(issue, err)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusBadRequest,
		Headers:    map[string]string{"Content-Type": ct},
		Body:       string(body),
	}
}
