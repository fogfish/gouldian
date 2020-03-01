//
//   Copyright 2019 Dmitry Kolesnikov, All Rights Reserved
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package gouldian_test

import (
	"testing"

	"github.com/fogfish/gouldian"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/core"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/path"
	"github.com/fogfish/it"
)

func TestServeSuccess(t *testing.T) {
	fun := gouldian.Serve(hello())
	req := mock.Input(mock.URL("/hello"))
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

func hello() core.Endpoint {
	return µ.GET(
		µ.Path(path.Is("hello")),
		µ.FMap(
			func() error { return µ.Ok().Text("Hello World!") },
		),
	)
}

/*
func TestServeFailure(t *testing.T) {
	fun := gouldian.Serve(unauthorized())
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/issue",
	}
	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
		"Content-Type":                 "application/json",
	}
	rsp, _ := fun(req)

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(401).
		If(rsp.Headers).Should().Equal(head).
		If(rsp.Body).Should().Equal("{\"type\":\"https://httpstatuses.com/401\",\"status\":401,\"title\":\"Unauthorized\",\"details\":\"some reason\"}")
}

func unauthorized() gouldian.Endpoint {
	return gouldian.Get().Path("issue").FMap(
		func() error { return gouldian.Unauthorized("some reason") },
	)
}

func TestServeNoMatch(t *testing.T) {
	fun := gouldian.Serve(hello())
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/issue",
	}
	head := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization, Accept",
		"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, OPTIONS",
		"Access-Control-Max-Age":       "600",
	}
	rsp, _ := fun(req)

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(501).
		If(rsp.Headers).Should().Equal(head)
}
*/
