// Package conformance_test runs the same request through all adapters
// (httpd, apigateway, apigatewayv2). Bare endpoints are not enough: routing,
// failure mapping, CORS and logging live in the server and the adapters.
package conformance_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian/v3"
	"github.com/fogfish/gouldian/v3/server/aws/apigateway"
	"github.com/fogfish/gouldian/v3/server/aws/apigatewayv2"
	"github.com/fogfish/gouldian/v3/server/httpd"
)

// Req is the adapter independent request
type Req struct {
	Method string
	Path   string
	Query  url.Values
	Header http.Header
	Body   []byte
	JWT    map[string]string
	Ctx    context.Context
}

// Resp is the adapter independent response
type Resp struct {
	Status int
	Header http.Header
	Body   []byte
}

type adapter struct {
	name string
	run  func(*testing.T, *µ.Server, Req) Resp
}

var adapters = []adapter{
	{"httpd", runHTTPD},
	{"apigateway", runAPIGW},
	{"apigatewayv2", runAPIGWv2},
}

// each runs the check against every adapter
func each(t *testing.T, srv func() *µ.Server, req Req, check func(*testing.T, Resp)) {
	t.Helper()
	for _, a := range adapters {
		t.Run(a.name, func(t *testing.T) {
			check(t, a.run(t, srv(), req))
		})
	}
}

func ctxOf(r Req) context.Context {
	if r.Ctx != nil {
		return r.Ctx
	}
	return context.Background()
}

func runHTTPD(t *testing.T, srv *µ.Server, r Req) Resp {
	target := r.Path
	if len(r.Query) != 0 {
		target += "?" + r.Query.Encode()
	}
	req := httptest.NewRequest(r.Method, target, strings.NewReader(string(r.Body))).WithContext(ctxOf(r))
	for k, v := range r.Header {
		req.Header[k] = v
	}
	w := httptest.NewRecorder()
	httpd.Handler(srv).ServeHTTP(w, req)

	res := w.Result()
	return Resp{Status: res.StatusCode, Header: res.Header, Body: w.Body.Bytes()}
}

func runAPIGW(t *testing.T, srv *µ.Server, r Req) Resp {
	evt := events.APIGatewayProxyRequest{
		HTTPMethod:                      r.Method,
		Path:                            r.Path,
		MultiValueQueryStringParameters: r.Query,
		MultiValueHeaders:               r.Header,
	}
	if len(r.Body) != 0 {
		evt.Body = base64.StdEncoding.EncodeToString(r.Body)
		evt.IsBase64Encoded = true
	}
	if r.JWT != nil {
		claims := map[string]interface{}{}
		for k, v := range r.JWT {
			claims[k] = v
		}
		evt.RequestContext.Authorizer = map[string]interface{}{"claims": claims}
	}

	out, err := apigateway.Handler(srv)(ctxOf(r), evt)
	if err != nil {
		t.Fatal(err)
	}

	h := http.Header{}
	for k, v := range out.Headers {
		h.Add(k, v)
	}
	for k, vs := range out.MultiValueHeaders {
		for _, v := range vs {
			h.Add(k, v)
		}
	}
	return Resp{Status: out.StatusCode, Header: h, Body: decode(t, out.Body, out.IsBase64Encoded)}
}

func runAPIGWv2(t *testing.T, srv *µ.Server, r Req) Resp {
	evt := events.APIGatewayV2HTTPRequest{
		RawPath:        r.Path,
		RawQueryString: r.Query.Encode(),
		Headers:        map[string]string{},
	}
	evt.RequestContext.HTTP.Method = r.Method
	for k, vs := range r.Header {
		if strings.EqualFold(k, "Cookie") {
			evt.Cookies = append(evt.Cookies, vs...)
			continue
		}
		evt.Headers[strings.ToLower(k)] = strings.Join(vs, ",")
	}
	if len(r.Body) != 0 {
		evt.Body = base64.StdEncoding.EncodeToString(r.Body)
		evt.IsBase64Encoded = true
	}
	if r.JWT != nil {
		evt.RequestContext.Authorizer = &events.APIGatewayV2HTTPRequestContextAuthorizerDescription{
			JWT: &events.APIGatewayV2HTTPRequestContextAuthorizerJWTDescription{Claims: r.JWT},
		}
	}

	out, err := apigatewayv2.Handler(srv)(ctxOf(r), evt)
	if err != nil {
		t.Fatal(err)
	}

	h := http.Header{}
	for k, v := range out.Headers {
		h.Add(k, v)
	}
	for _, c := range out.Cookies {
		h.Add("Set-Cookie", c)
	}
	return Resp{Status: out.StatusCode, Header: h, Body: decode(t, out.Body, out.IsBase64Encoded)}
}

func decode(t *testing.T, body string, isBase64 bool) []byte {
	if !isBase64 {
		return []byte(body)
	}
	b, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// eachAWS runs the check against adapters that carry authorizer claims
func eachAWS(t *testing.T, srv func() *µ.Server, req Req, check func(*testing.T, Resp)) {
	t.Helper()
	for _, a := range adapters[1:] {
		t.Run(a.name, func(t *testing.T) {
			check(t, a.run(t, srv(), req))
		})
	}
}
