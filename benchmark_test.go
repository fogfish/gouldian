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

package gouldian_test

import (
	"context"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/headers"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/optics"
	"github.com/fogfish/gouldian/server/httpd"
	"net/http"
	"path/filepath"
	"testing"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

//
// Microbenchmark
//

/*

Path Pattern with 1 param

*/

type MyT1 struct {
	Name string
}

var (
	name  = optics.ForProduct1(MyT1{})
	path1 = µ.Path("user", name)
	foo1  = µ.GET(path1)
	req1  = mock.Input(mock.URL("/user/123456"))
)

func BenchmarkPathParam1(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		foo1(req1)
	}
}

func BenchmarkServerParam1(mb *testing.B) {
	w := new(mockResponseWriter)
	router := httpd.Serve(
		µ.Join(
			foo1,
			func(c *µ.Context) error { return nil },
		),
	)
	r, _ := http.NewRequest("GET", "/user/123456", nil)

	mb.ReportAllocs()
	mb.ResetTimer()
	for i := 0; i < mb.N; i++ {
		router.ServeHTTP(w, r)
	}
}

/*

Path Pattern with 5 param

*/

type MyT5 struct{ A, B, C, D, E string }

var (
	a, b, c, d, e = optics.ForProduct5(MyT5{})
	path5         = µ.Path(a, b, c, d, e)
	foo5          = µ.GET(path5)
	req5          = mock.Input(mock.URL("/a/b/c/d/e"))
)

func BenchmarkPathParam5(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		foo5(req5)
	}
}

func BenchmarkServerParam5(mb *testing.B) {
	w := new(mockResponseWriter)
	router := httpd.Serve(
		µ.Join(
			foo5,
			func(c *µ.Context) error { return nil },
		),
	)
	r, _ := http.NewRequest("GET", "/a/b/c/d/e", nil)

	mb.ReportAllocs()
	mb.ResetTimer()
	for i := 0; i < mb.N; i++ {
		router.ServeHTTP(w, r)
	}
}

/*

Lens decode with 1 param

*/

func BenchmarkLensForProduct1(mb *testing.B) {
	ctx := µ.NewContext(context.Background())
	ctx.Put(name, "123456")

	var val MyT1

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		ctx.Get(&val)
	}
}

/*

Lens decode with 1 param

*/

func BenchmarkLensForProduct5(mb *testing.B) {
	ctx := µ.NewContext(context.Background())
	ctx.Put(a, "a")
	ctx.Put(b, "b")
	ctx.Put(c, "c")
	ctx.Put(d, "d")
	ctx.Put(e, "e")

	var val MyT5

	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		ctx.Get(&val)
	}
}

/*

Endpoint decode with 1 param

*/

var endpoint1 = µ.GET(
	path1,
	µ.FMap(func(ctx *µ.Context) error {
		var req MyT1
		if err := ctx.Get(&req); err != nil {
			return µ.Status.BadRequest(µ.WithIssue(err))
		}

		return µ.Status.OK(
			headers.ContentType.Value(headers.TextPlain),
			headers.Server.Value("echo"),
			µ.WithText(req.Name),
		)
	}),
)

func BenchmarkEndpoint1(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		endpoint1(req1)
	}
}

/*

Endpoint decode with 5 param

*/

var endpoint5 = µ.GET(
	path5,
	µ.FMap(func(ctx *µ.Context) error {
		var req MyT5
		if err := ctx.Get(&req); err != nil {
			return µ.Status.BadRequest(µ.WithIssue(err))
		}

		return µ.Status.OK(
			headers.ContentType.Value(headers.TextPlain),
			headers.Server.Value("echo"),
			µ.WithText(filepath.Join(req.A, req.B, req.C, req.D, req.E)),
		)
	}),
)

func BenchmarkEndpoint5(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		endpoint5(req5)
	}
}
