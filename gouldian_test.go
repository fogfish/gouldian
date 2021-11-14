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
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/optics"
	"testing"
)

/*
func TestServeSuccess(t *testing.T) {
	fun := µ.Serve(hello())
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

func hello() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("hello")),
		µ.FMap(
			func() error { return µ.Ok().Text("Hello World!") },
		),
	)
}

func TestServeFailure(t *testing.T) {
	fun := µ.Serve(unauthorized())
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
	var issue µ.Issue

	it.Ok(t).
		If(rsp.StatusCode).Should().Equal(401).
		If(rsp.Headers).Should().Equal(head).
		If(json.Unmarshal([]byte(rsp.Body), &issue)).Should().Equal(nil).
		If(issue.Type).Should().Equal("https://httpstatuses.com/401").
		If(issue.Status).Should().Equal(401).
		If(issue.Title).Should().Equal("Unauthorized").
		If(issue.ID).ShouldNot().Equal("")
}

func unauthorized() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("issue")),
		µ.FMap(func() error { return µ.Unauthorized(errors.New("some reason")) }),
	)
}

func TestServeNoMatch(t *testing.T) {
	fun := µ.Serve(hello())
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
		If(rsp.StatusCode).Should().Equal(501).
		If(rsp.Headers).Should().Equal(head)
}

func TestServeNoMatchLogger(t *testing.T) {
	fun := µ.Serve(µ.NoMatchLogger())
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
		If(rsp.StatusCode).Should().Equal(501).
		If(rsp.Headers).Should().Equal(head)
}

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

//
// Microbenchmark
//

type MyT1 struct{ Name string }

var (
	uid = optics.Lenses1(MyT1{})

	foo1 = µ.GET(
		// µ.Param("user").To(uid),
		µ.Path("user", uid),
		// µ.FMap(func(c µ.Context) error {
		// 	var myt MyT1
		// 	c.Get(&myt)
		// 	return nil
		// }),
	)

	req1 = mock.Input(mock.URL("/user/123456"))
	// req1 = mock.Input(mock.URL("/?user=gordon"))
)

//
// Route with Param (no write)
/* */
func BenchmarkPathParam1(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		foo1(req1)
	}
}

/* */

type MyT5 struct{ A, B, C, D, E string }

var (
	a, b, c, d, e = optics.Lenses5(MyT5{})

	foo5 = µ.GET(
		// µ.Param("a").To(a),
		// µ.Param("b").To(b),
		// µ.Param("c").To(c),
		// µ.Param("d").To(d),
		// µ.Param("e").To(e),
		µ.Path(a, b, c, d, e),
		// µ.FMap(func(c µ.Context) error {
		// 	var myt MyT5
		// 	c.Get(&myt)
		// 	return nil
		// }),
	)

	req5 = mock.Input(mock.URL("/a/b/c/d/e"))
	// req5 = mock.Input(mock.URL("/?a=a&b=b&c=c&d=d&e=e"))
)

//
// Route with 5 Params (no write)
/* */
func BenchmarkPathParam5(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		foo5(req5)
	}
}

/* */

//
// Route with 20 Params (no write)
