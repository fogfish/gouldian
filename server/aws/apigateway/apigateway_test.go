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
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	µ "github.com/fogfish/gouldian/v2"
	ø "github.com/fogfish/gouldian/v2/output"

	"github.com/fogfish/gouldian/v2/server/aws/apigateway"
	"github.com/fogfish/it"
)

func TestServeMatch(t *testing.T) {
	api := apigateway.Serve(mock("echo"))
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/echo?foo=bar",
		Headers:    map[string]string{"Accept": "*/*"},
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
	api := apigateway.Serve(mock("echo"))
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
	api := apigateway.Serve(mock("h%rt"))
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/h%rt",
	}

	out, err1 := api(req)
	it.Ok(t).If(err1).Must().Equal(nil)

	it.Ok(t).
		If(out.StatusCode).Should().Equal(http.StatusBadRequest)
}

func TestServeAndCommit(t *testing.T) {
	cnt := 0
	api := apigateway.ServeAndCommit(
		func() { cnt = cnt + 1 },
		mock("echo"),
	)
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/echo?foo=bar",
		Headers:    map[string]string{"Accept": "*/*"},
	}

	out, err1 := api(req)
	it.Ok(t).If(err1).Must().Equal(nil)

	it.Ok(t).
		If(out.StatusCode).Should().Equal(http.StatusOK).
		If(out.Headers["Server"]).Should().Equal("echo").
		If(out.Headers["Content-Type"]).Should().Equal("text/plain").
		If(out.Body).Should().Equal("echo").
		If(cnt).Should().Equal(1)
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
