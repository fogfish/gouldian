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
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/optics"
	ƒ "github.com/fogfish/gouldian/output"
	"testing"
)

//
// Microbenchmark
//

type MyT1 struct {
	Name string `lens:"application/json"`
}

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
// Route with 128 bytes of payload

/* */

type MyT128 struct{ Text string }

var (
	body = optics.Lenses1(MyT128{})

	foo128 = µ.GET(
		// µ.Param("a").To(a),
		// µ.Param("b").To(b),
		// µ.Param("c").To(c),
		// µ.Param("d").To(d),
		// µ.Param("e").To(e),
		µ.Body(body),
		// µ.FMap(func(c µ.Context) error {
		// 	var myt MyT5
		// 	c.Get(&myt)
		// 	return nil
		// }),
	)

	req128 = mock.Input(mock.Text("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	// req5 = mock.Input(mock.URL("/?a=a&b=b&c=c&d=d&e=e"))
)

//
// Route with 5 Params (no write)
/* */
func BenchmarkBody128(mb *testing.B) {
	mb.ReportAllocs()
	mb.ResetTimer()

	for i := 0; i < mb.N; i++ {
		foo128(req128)
	}
}

type TEcho struct {
	Echo string
}

var lensEcho = optics.Lenses1(TEcho{})

var echo = µ.GET(
	µ.Path("echo", lensEcho),
	func(r *µ.Input) error {
		var req TEcho
		if err := r.Context.Get(&req); err != nil {
			return µ.Status.BadRequest(ƒ.Issue(err))
		}

		return µ.Status.OK(
			ƒ.ContentType.Text,
			ƒ.Server.Is("echo"),
			ƒ.Text(req.Echo),
		)
	},
)

var reqEcho = mock.Input(mock.URL("/echo/123456"))

func BenchmarkEcho(mb *testing.B) {
	for i := 0; i < mb.N; i++ {
		echo(reqEcho)
	}
}
