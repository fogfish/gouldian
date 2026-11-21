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

// Package apigateway adapts the gouldian server to AWS API Gateway
// (REST API, payload format 1.0).
package apigateway

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/internal/media"
)

// Func is AWS Lambda handler
type Func = func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

/*
Serve builds Lambda handler for sequence of endpoints

	lambda.Start(apigateway.Serve( ... ))
*/
func Serve(endpoints ...µ.Routable) Func {
	return Handler(µ.NewServer().Mount(endpoints...))
}

// ServeAndCommit builds Lambda handler that executes commit after each request.
func ServeAndCommit(commit func(), endpoints ...µ.Routable) Func {
	h := Serve(endpoints...)

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		ret, err := h(ctx, req)
		commit()
		return ret, err
	}
}

/*
Handler adapts the server to AWS Lambda. The context of the invocation (the
deadline, the cancellation) is propagated to endpoints.

	api := µ.NewServer(µ.WithCORS(...)).Mount( ... )
	lambda.Start(apigateway.Handler(api))
*/
func Handler(srv *µ.Server) Func {
	return func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req, err := Request(ctx, &r)
		if err != nil {
			return badRequest(r.Path, err), nil
		}

		c := srv.NewContext(ctx, req)
		c.RequestID = r.RequestContext.RequestID
		c.JWT = jwtFromAuthorizer(&r)

		out := srv.Dispatch(c)
		evt := Response(out)
		out.Free()
		srv.Release(c)

		return evt, nil
	}
}

// Request is events.APIGatewayProxyRequest ⟼ http.Request
func Request(ctx context.Context, r *events.APIGatewayProxyRequest) (*http.Request, error) {
	body, err := requestBody(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, r.HTTPMethod, "/", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.URL = &url.URL{Path: r.Path}

	if len(r.MultiValueHeaders) != 0 {
		for header, values := range r.MultiValueHeaders {
			for _, value := range values {
				req.Header.Add(header, value)
			}
		}
	} else {
		for header, value := range r.Headers {
			req.Header.Set(header, value)
		}
	}
	req.Host = req.Header.Get("Host")

	q := make(url.Values)
	if len(r.MultiValueQueryStringParameters) != 0 {
		for key, values := range r.MultiValueQueryStringParameters {
			for _, value := range values {
				q.Add(key, value)
			}
		}
	} else {
		for key, value := range r.QueryStringParameters {
			q.Add(key, value)
		}
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func requestBody(r *events.APIGatewayProxyRequest) ([]byte, error) {
	if !r.IsBase64Encoded {
		return []byte(r.Body), nil
	}

	buf, err := base64.StdEncoding.DecodeString(r.Body)
	if err != nil {
		return nil, fmt.Errorf("body is not valid base64: %w", err)
	}
	return buf, nil
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

// Response is µ.Output ⟼ events.APIGatewayProxyResponse. Headers with
// multiple values (Set-Cookie) use MultiValueHeaders, binary bodies are base64
// encoded.
func Response(out *µ.Output) events.APIGatewayProxyResponse {
	evt := events.APIGatewayProxyResponse{
		StatusCode: out.Status,
	}

	count := make(map[string]int, len(out.Headers))
	for _, h := range out.Headers {
		count[h.Header]++
	}

	var contentType string
	for _, h := range out.Headers {
		if h.Header == "Content-Type" {
			contentType = h.Value
		}

		if count[h.Header] == 1 {
			if evt.Headers == nil {
				evt.Headers = make(map[string]string, len(out.Headers))
			}
			evt.Headers[h.Header] = h.Value
			continue
		}

		if evt.MultiValueHeaders == nil {
			evt.MultiValueHeaders = make(map[string][]string)
		}
		evt.MultiValueHeaders[h.Header] = append(evt.MultiValueHeaders[h.Header], h.Value)
	}

	if media.IsBinary(contentType, out.Body) {
		evt.IsBase64Encoded = true
		evt.Body = base64.StdEncoding.EncodeToString(out.Body)
	} else {
		evt.Body = string(out.Body)
	}

	return evt
}

func badRequest(path string, err error) events.APIGatewayProxyResponse {
	issue := µ.NewIssue(http.StatusBadRequest)
	issue.Detail = err.Error()
	issue.Instance = path
	ct, body := µ.DefaultIssue(issue, err)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Headers:    map[string]string{"Content-Type": ct},
		Body:       string(body),
	}
}
