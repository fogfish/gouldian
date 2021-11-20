package apigateway_test

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/server/aws/apigateway"
	"github.com/fogfish/it"
)

func TestServeMatch(t *testing.T) {
	api := mock("echo")
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/echo",
	}

	out, err1 := api(req)
	it.Ok(t).If(err1).Must().Equal(nil)

	it.Ok(t).
		If(out.StatusCode).Should().Equal(http.StatusOK).
		If(out.Headers["Server"]).Should().Equal("echo").
		If(out.Headers["Content-Type"]).Should().Equal("text/plain").
		If(out.Body).Should().Equal("echo")

	// 	"Access-Control-Allow-Origin":  "*",
	// 	"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
	// 	"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
	// 	"Access-Control-Max-Age":       "600",
	// }
}

func TestServeNoMatch(t *testing.T) {
	api := mock("echo")
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/foo",
	}

	out, err1 := api(req)
	it.Ok(t).If(err1).Must().Equal(nil)

	it.Ok(t).
		If(out.StatusCode).Should().Equal(http.StatusNotImplemented).
		If(out.Headers["Content-Type"]).Should().Equal("application/json").
		If(out.Body).ShouldNot().Equal("")

	// "Access-Control-Allow-Origin":  "*",
	// "Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
	// "Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
	// "Access-Control-Max-Age":       "600",
}

func TestServeMatchUnescaped(t *testing.T) {
	api := mock("h%rt")
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/h%rt",
	}

	out, err1 := api(req)
	it.Ok(t).If(err1).Must().Equal(nil)

	it.Ok(t).
		If(out.StatusCode).Should().Equal(http.StatusOK).
		If(out.Headers["Server"]).Should().Equal("echo").
		If(out.Headers["Content-Type"]).Should().Equal("text/plain").
		If(out.Body).Should().Equal("echo")
}

func mock(path string) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return apigateway.Serve(
		µ.GET(
			µ.Path(path),
			µ.FMap(func(ctx µ.Context) error {
				return µ.Status.OK(
					headers.ContentType.Value(headers.TextPlain),
					headers.Server.Value("echo"),
					µ.WithText("echo"),
				)
			}),
		),
	)
}
