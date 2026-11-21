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

package apigateway_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian/v3"
	ø "github.com/fogfish/gouldian/v3/output"

	"github.com/fogfish/gouldian/v3/server/aws/apigateway"
	"github.com/fogfish/it/v2"
)

func TestServeMatch(t *testing.T) {
	api := apigateway.Serve(mock("echo"))
	req := events.APIGatewayProxyRequest{
		HTTPMethod:            "GET",
		Path:                  "/echo",
		Headers:               map[string]string{"Accept": "*/*"},
		QueryStringParameters: map[string]string{"foo": "bar"},
	}

	out, err1 := api(context.Background(), req)
	it.Then(t).
		Must(it.Equal(err1, nil))

	it.Then(t).
		Should(it.Equal(out.StatusCode, http.StatusOK)).
		Should(it.Equal(out.Headers["Server"], "echo")).
		Should(it.Equal(out.Headers["Content-Type"], "text/plain")).
		Should(it.Equal(out.Body, "echo"))

	// 	"Access-Control-Allow-Origin":  "*",
	// 	"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
	// 	"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
	// 	"Access-Control-Max-Age":       "600",
	// }
}

func TestServeNoMatch(t *testing.T) {
	api := apigateway.Serve(mock("echo"))
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/foo",
	}

	out, err1 := api(context.Background(), req)
	it.Then(t).
		Must(it.Equal(err1, nil))

	it.Then(t).
		Should(it.Equal(out.StatusCode, http.StatusNotFound)).
		Should(it.Equal(out.Headers["Content-Type"], "application/problem+json")).
		ShouldNot(it.Equal(out.Body, ""))

	// "Access-Control-Allow-Origin":  "*",
	// "Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
	// "Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
	// "Access-Control-Max-Age":       "600",
}

func TestServeInvalidMethod(t *testing.T) {
	api := apigateway.Serve(mock("echo"))
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "G E T",
		Path:       "/echo",
	}

	out, err1 := api(context.Background(), req)
	it.Then(t).
		Must(it.Equal(err1, nil))

	it.Then(t).
		Should(it.Equal(out.StatusCode, http.StatusBadRequest)).
		Should(it.Equal(out.Headers["Content-Type"], "application/problem+json"))
}

func TestServeAndCommit(t *testing.T) {
	cnt := 0
	api := apigateway.ServeAndCommit(
		func() { cnt = cnt + 1 },
		mock("echo"),
	)
	req := events.APIGatewayProxyRequest{
		HTTPMethod:            "GET",
		Path:                  "/echo",
		Headers:               map[string]string{"Accept": "*/*"},
		QueryStringParameters: map[string]string{"foo": "bar"},
	}

	out, err1 := api(context.Background(), req)
	it.Then(t).
		Must(it.Equal(err1, nil))

	it.Then(t).
		Should(it.Equal(out.StatusCode, http.StatusOK)).
		Should(it.Equal(out.Headers["Server"], "echo")).
		Should(it.Equal(out.Headers["Content-Type"], "text/plain")).
		Should(it.Equal(out.Body, "echo")).
		Should(it.Equal(cnt, 1))
}

func mock(path string) µ.Routable {
	return µ.GET(
		µ.URI(µ.Path(path)),
		func(ctx *µ.Context) error {
			return ø.Status.OK(
				ø.Server.Set("echo"),
				ø.ContentType.TextPlain,
				ø.Send("echo"),
			)
		},
	)
}

func TestRequestMapping(t *testing.T) {
	type req struct{ Raw []byte }
	var token µ.Token
	var accept []string
	var query string
	var requestID string
	var raw []byte

	api := apigateway.Serve(µ.POST(µ.URI(µ.Path("a")),
		µ.Body(µ.Optics1[req, []byte]()),
		µ.FMap(func(c *µ.Context, r *req) error {
			token, raw = c.JWT, r.Raw
			accept = c.Request.Header.Values("Accept")
			query = c.Request.URL.RawQuery
			requestID = c.RequestID
			return nil
		})))

	evt := events.APIGatewayProxyRequest{
		HTTPMethod:                      "POST",
		Path:                            "/a",
		MultiValueHeaders:               map[string][]string{"Accept": {"text/html", "application/json"}},
		MultiValueQueryStringParameters: map[string][]string{"rels": {"a", "b"}},
		Body:                            "AAEC",
		IsBase64Encoded:                 true,
	}
	evt.RequestContext.RequestID = "rid-1"
	evt.RequestContext.Authorizer = map[string]interface{}{"claims": map[string]interface{}{"sub": "alice", "exp": 1234.0}}

	out, err := api(context.Background(), evt)
	it.Then(t).
		Must(it.Nil(err)).
		Should(it.Equal(out.StatusCode, http.StatusNoContent)).
		Should(it.Equal(token.Sub(), "alice")).
		Should(it.Equal(token.Exp(), "1234")).
		Should(it.Equiv(accept, []string{"text/html", "application/json"})).
		Should(it.Equal(query, "rels=a&rels=b")).
		Should(it.Equal(requestID, "rid-1")).
		Should(it.Equiv(raw, []byte{0, 1, 2}))

	// IAM identity
	evt.RequestContext.Authorizer = nil
	evt.RequestContext.Identity.UserArn = "arn:aws:iam::1:user/a"
	api(context.Background(), evt)
	it.Then(t).Should(it.Equal(token.Username(), "arn:aws:iam::1:user/a"))
}

func TestResponseMapping(t *testing.T) {
	out := µ.NewOutput(200)
	out.AddHeader("Set-Cookie", "a=1")
	out.AddHeader("Set-Cookie", "b=2")
	out.SetHeader("Content-Type", "image/png")
	out.Body = []byte{0x89, 'P', 'N', 'G', 0xff}

	evt := apigateway.Response(out)
	it.Then(t).
		Should(it.Equiv(evt.MultiValueHeaders["Set-Cookie"], []string{"a=1", "b=2"})).
		Should(it.Equal(evt.Headers["Content-Type"], "image/png")).
		Should(it.Equal(evt.IsBase64Encoded, true))
}

func TestInvalidBase64(t *testing.T) {
	api := apigateway.Serve()
	out, _ := api(context.Background(), events.APIGatewayProxyRequest{HTTPMethod: "POST", Path: "/", Body: "***", IsBase64Encoded: true})
	it.Then(t).Should(it.Equal(out.StatusCode, http.StatusBadRequest))
}
