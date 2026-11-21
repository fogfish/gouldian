// Package exampletest runs examples through the HTTP adapters
// (net/http and AWS API Gateway HTTP API payload 2.0).
package exampletest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/server/aws/apigatewayv2"
	"github.com/fogfish/gouldian/v3/server/httpd"
)

// Case is a request and the expectations about the response
type Case struct {
	Name   string
	Method string // GET by default
	Path   string
	Query  url.Values
	Header map[string]string
	Body   string

	Status   int    // 200 by default
	Contains string // substring of the response body
}

// Run executes cases against the routes using every adapter
func Run(t *testing.T, routes []µ.Routable, cases []Case) {
	t.Helper()

	for _, c := range cases {
		if c.Method == "" {
			c.Method = "GET"
		}
		if c.Status == 0 {
			c.Status = http.StatusOK
		}

		t.Run("httpd/"+c.Name, func(t *testing.T) {
			target := c.Path
			if len(c.Query) != 0 {
				target += "?" + c.Query.Encode()
			}
			req := httptest.NewRequest(c.Method, target, strings.NewReader(c.Body))
			for k, v := range c.Header {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			httpd.Serve(routes...).ServeHTTP(w, req)

			body, _ := io.ReadAll(w.Result().Body)
			check(t, c, w.Code, string(body))
		})

		t.Run("apigatewayv2/"+c.Name, func(t *testing.T) {
			evt := events.APIGatewayV2HTTPRequest{
				RawPath:        c.Path,
				RawQueryString: c.Query.Encode(),
				Headers:        map[string]string{},
				Body:           c.Body,
			}
			evt.RequestContext.HTTP.Method = c.Method
			for k, v := range c.Header {
				evt.Headers[strings.ToLower(k)] = v
			}

			out, err := apigatewayv2.Serve(routes...)(context.Background(), evt)
			if err != nil {
				t.Fatal(err)
			}
			check(t, c, out.StatusCode, out.Body)
		})
	}
}

func check(t *testing.T, c Case, status int, body string) {
	t.Helper()
	if status != c.Status {
		t.Fatalf("status %d, want %d, body %q", status, c.Status, body)
	}
	if !strings.Contains(body, c.Contains) {
		t.Fatalf("body %q does not contain %q", body, c.Contains)
	}
}
