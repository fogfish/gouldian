package apigatewayv2_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian/v3"
	ø "github.com/fogfish/gouldian/v3/output"
	"github.com/fogfish/gouldian/v3/server/aws/apigatewayv2"
	"github.com/fogfish/it/v2"
)

func request(method, path string) events.APIGatewayV2HTTPRequest {
	evt := events.APIGatewayV2HTTPRequest{RawPath: path, Headers: map[string]string{}}
	evt.RequestContext.HTTP.Method = method
	return evt
}

func TestServe(t *testing.T) {
	api := apigatewayv2.Serve(µ.GET(µ.URI(µ.Path("echo")), func(*µ.Context) error {
		return ø.Status.OK(ø.ContentType.TextPlain, ø.Send("echo"))
	}))

	out, err := api(context.Background(), request("GET", "/echo"))
	it.Then(t).
		Must(it.Nil(err)).
		Should(it.Equal(out.StatusCode, http.StatusOK)).
		Should(it.Equal(out.Headers["Content-Type"], "text/plain")).
		Should(it.Equal(out.Body, "echo")).
		Should(it.Equal(out.IsBase64Encoded, false))

	out, _ = api(context.Background(), request("GET", "/other"))
	it.Then(t).
		Should(it.Equal(out.StatusCode, http.StatusNotFound)).
		Should(it.Equal(out.Headers["Content-Type"], "application/problem+json"))
}

func TestRequestMapping(t *testing.T) {
	type req struct{ Raw []byte }
	var token µ.Token
	var cookie, rawQuery, requestID string
	var raw []byte

	api := apigatewayv2.Serve(µ.POST(µ.URI(µ.Path("a"), µ.Path("b c")),
		µ.Body(µ.Optics1[req, []byte]()),
		µ.FMap(func(c *µ.Context, r *req) error {
			token, raw = c.JWT, r.Raw
			cookie = c.Request.Header.Get("Cookie")
			rawQuery = c.Request.URL.RawQuery
			requestID = c.RequestID
			return nil
		})))

	evt := request("POST", "/a/b%20c")
	evt.RawQueryString = "x=1&y=2"
	evt.Cookies = []string{"a=1", "b=2"}
	evt.Body = base64.StdEncoding.EncodeToString([]byte{0, 1, 2})
	evt.IsBase64Encoded = true
	evt.RequestContext.RequestID = "rid-1"
	evt.RequestContext.Authorizer = &events.APIGatewayV2HTTPRequestContextAuthorizerDescription{
		JWT: &events.APIGatewayV2HTTPRequestContextAuthorizerJWTDescription{
			Claims: map[string]string{"sub": "alice", "exp": "1234"},
			Scopes: []string{"ro", "rw"},
		},
	}

	out, _ := api(context.Background(), evt)
	it.Then(t).
		Should(it.Equal(out.StatusCode, http.StatusNoContent)).
		Should(it.Equal(token.Sub(), "alice")).
		Should(it.Equal(token.Scope(), "ro rw")).
		Should(it.Equal(cookie, "a=1; b=2")).
		Should(it.Equal(rawQuery, "x=1&y=2")).
		Should(it.Equal(requestID, "rid-1")).
		Should(it.Equiv(raw, []byte{0, 1, 2}))
}

func TestAuthorizers(t *testing.T) {
	var token µ.Token
	api := apigatewayv2.Serve(µ.GET(µ.URI(), func(c *µ.Context) error { token = c.JWT; return nil }))

	evt := request("GET", "/")
	evt.RequestContext.Authorizer = &events.APIGatewayV2HTTPRequestContextAuthorizerDescription{
		IAM: &events.APIGatewayV2HTTPRequestContextAuthorizerIAMDescription{UserARN: "arn:aws:iam::1:user/a", UserID: "AID"},
	}
	api(context.Background(), evt)
	it.Then(t).
		Should(it.Equal(token.Username(), "arn:aws:iam::1:user/a")).
		Should(it.Equal(token.Sub(), "AID"))

	evt.RequestContext.Authorizer = &events.APIGatewayV2HTTPRequestContextAuthorizerDescription{
		Lambda: map[string]any{"role": "admin", "level": 3.0},
	}
	api(context.Background(), evt)
	level, _ := token.Number("level")
	it.Then(t).Should(it.Equal(token.String("role"), "admin")).Should(it.Equal(level, 3.0))
}

func TestResponseMapping(t *testing.T) {
	out := µ.NewOutput(200)
	out.AddHeader("Set-Cookie", "a=1")
	out.AddHeader("Set-Cookie", "b=2")
	out.AddHeader("X-Multi", "1")
	out.AddHeader("X-Multi", "2")
	out.SetHeader("Content-Type", "image/png")
	out.Body = []byte{0x89, 'P', 'N', 'G', 0xff}

	evt := apigatewayv2.Response(out)
	it.Then(t).
		Should(it.Equiv(evt.Cookies, []string{"a=1", "b=2"})).
		Should(it.Equal(evt.Headers["X-Multi"], "1,2")).
		Should(it.Equal(evt.IsBase64Encoded, true)).
		Should(it.Equal(evt.Body, base64.StdEncoding.EncodeToString(out.Body)))
}

func TestBadRequest(t *testing.T) {
	api := apigatewayv2.Serve()

	evt := request("GET", "/x")
	evt.Body, evt.IsBase64Encoded = "***", true
	out, _ := api(context.Background(), evt)
	it.Then(t).Should(it.Equal(out.StatusCode, http.StatusBadRequest))

	out, _ = api(context.Background(), request("G E T", "/x"))
	it.Then(t).Should(it.Equal(out.StatusCode, http.StatusBadRequest))
}
