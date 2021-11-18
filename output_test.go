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
	"errors"
	"reflect"
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	ƒ "github.com/fogfish/gouldian/output"
	"github.com/fogfish/it"
)

// func TestSuccess(t *testing.T) {
// 	output(t, µ.NewSuccess(200), µ.NewSuccess(200))
// 	output(t, µ.Status.OK(), µ.NewSuccess(200))
// 	output(t, µ.Status.Created(), µ.NewSuccess(201))
// 	output(t, µ.Status.Accepted(), µ.NewSuccess(202))
// 	output(t, µ.Status.NoContent(), µ.NewSuccess(204))
// 	output(t, µ.Status.MovedPermanently("/"), µ.NewSuccess(301).With("Location", "/"))
// 	output(t, µ.Status.Found("/"), µ.NewSuccess(302).With("Location", "/"))
// 	output(t, µ.Status.SeeOther("/"), µ.NewSuccess(303).With("Location", "/"))
// 	output(t, µ.Status.NotModified("/"), µ.NewSuccess(304).With("Location", "/"))
// 	output(t, µ.Status.TemporaryRedirect("/"), µ.NewSuccess(307).With("Location", "/"))
// 	output(t, µ.Status.PermanentRedirect("/"), µ.NewSuccess(308).With("Location", "/"))
// }

// func TestIssue(t *testing.T) {
// 	err := fmt.Errorf("issue")
// 	issue(t, µ.NewFailure(500, err), µ.NewFailure(500, err))
// 	issue(t, µ.Status.BadRequest(err), µ.NewFailure(400, err))
// 	issue(t, µ.Status.Unauthorized(err), µ.NewFailure(401, err))
// 	issue(t, µ.Status.Forbidden(err), µ.NewFailure(403, err))
// 	issue(t, µ.Status.NotFound(err), µ.NewFailure(404, err))
// 	issue(t, µ.Status.InternalServerError(err), µ.NewFailure(500, err))
// 	issue(t, µ.Status.NotImplemented(err), µ.NewFailure(501, err))
// 	issue(t, µ.Status.ServiceUnavailable(err), µ.NewFailure(503, err))
// }

type myT struct {
	A string
}

func TestJSON(t *testing.T) {
	output(t,
		µ.Status.OK(ƒ.JSON(myT{"Hello"})),
		µ.Status.OK(ƒ.JSON(myT{"Hello"})),
	)
}

func TestOutputText(t *testing.T) {
	output(t,
		µ.Status.OK(ƒ.Text("Hello")),
		µ.Status.OK(ƒ.Text("Hello")),
	)
}

func TestOutputBytes(t *testing.T) {
	output(t,
		µ.Status.OK(ƒ.Bytes([]byte("Hello"))),
		µ.Status.OK(ƒ.Bytes([]byte("Hello"))),
	)
}

// func TestErrorOnJSON(t *testing.T) {
// 	output := µ.Status.OK().JSON(make(chan int))

// 	it.Ok(t).
// 		If(output.Status).Should().Equal(http.StatusInternalServerError)
// }

func output(t *testing.T, a, b error) {
	t.Helper()
	foo := µ.GET(func(*µ.Input) error { return a })
	req := mock.Input()

	it.Ok(t).
		If(foo(req)).Should().
		Assert(
			func(be interface{}) bool {
				var out error
				if errors.As(be.(error), &out) {
					return reflect.DeepEqual(b, out)
				}
				return false
			},
		)
}

// func issue(t *testing.T, a, b µ.Output) {
// 	t.Helper()
// 	foo := µ.GET(func(*µ.Input) error { return a })
// 	req := mock.Input()

// 	it.Ok(t).
// 		If(foo(req)).Should().
// 		Assert(
// 			func(be interface{}) bool {
// 				if v, ok := be.(error); ok {
// 					if v.Error() == b.Error() {
// 						return true
// 					}
// 				}
// 				return false
// 			},
// 		)
// }
