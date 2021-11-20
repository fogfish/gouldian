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
	api := mock()
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
	api := mock()
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

/*



func TestServeUnescapedPath(t *testing.T) {
	fun := µ.Serve(unescaped())
	req := mock.Input(mock.URL("/"))
	req.APIGatewayProxyRequest.Path = "/h%rt"
	req.Path = []string{"h%rt"}
	rsp, _ := fun(req.APIGatewayProxyRequest)

	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "text/plain",
	}

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(200).
		If(rsp.Headers).Should().Equal(head).
		If(rsp.Body).Should().Equal("Hello World!")
}

func unescaped() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("h%rt")),
		µ.FMap(
			func() error { return µ.Ok().Text("Hello World!") },
		),
	)
}

func TestServeUnknownError(t *testing.T) {
	fun := µ.Serve(unknown())
	req := mock.Input(mock.URL("/"))
	req.APIGatewayProxyRequest.Path = "/h%rt"
	req.Path = []string{"h%rt"}
	rsp, _ := fun(req.APIGatewayProxyRequest)

	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "application/json",
	}

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(500).
		If(rsp.Headers).Should().Equal(head)
}

func unknown() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("h%rt")),
		µ.FMap(
			func() error { return fmt.Errorf("Unknown error") },
		),
	)
}
*/

func mock() func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return apigateway.Serve(
		µ.GET(
			µ.Path("echo"),
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
