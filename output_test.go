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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestSuccess(t *testing.T) {
	uri, _ := url.Parse("/")

	output(t, µ.Success(200), µ.Success(200))
	output(t, µ.Ok(), µ.Success(200))
	output(t, µ.Created(), µ.Success(201))
	output(t, µ.Accepted(), µ.Success(202))
	output(t, µ.NoContent(), µ.Success(204))
	output(t, µ.MovedPermanently(*uri), µ.Success(301).With("Location", "/"))
	output(t, µ.Found(*uri), µ.Success(302).With("Location", "/"))
	output(t, µ.SeeOther(*uri), µ.Success(303).With("Location", "/"))
	output(t, µ.NotModified(*uri), µ.Success(304).With("Location", "/"))
	output(t, µ.TemporaryRedirect(*uri), µ.Success(307).With("Location", "/"))
	output(t, µ.PermanentRedirect(*uri), µ.Success(308).With("Location", "/"))
}

func TestIssue(t *testing.T) {
	err := fmt.Errorf("issue")
	issue(t, µ.Failure(500), µ.Failure(500))
	issue(t, µ.BadRequest(err), µ.Failure(400).Reason(err))
	issue(t, µ.Unauthorized(err), µ.Failure(401).Reason(err))
	issue(t, µ.Forbidden(err), µ.Failure(403).Reason(err))
	issue(t, µ.NotFound(err), µ.Failure(404).Reason(err))
	issue(t, µ.InternalServerError(err), µ.Failure(500).Reason(err))
	issue(t, µ.NotImplemented(err), µ.Failure(501).Reason(err))
	issue(t, µ.ServiceUnavailable(err), µ.Failure(503).Reason(err))
}

type myT struct {
	A string
}

func TestJSON(t *testing.T) {
	output(t,
		µ.Ok().JSON(myT{"Hello"}),
		µ.Ok().JSON(myT{"Hello"}),
	)
}

func TestErrorOnJSON(t *testing.T) {
	output := µ.Ok().JSON(make(chan int))

	it.Ok(t).
		If(output.Status).Should().Equal(http.StatusInternalServerError)
}

func output(t *testing.T, a, b *µ.Output) {
	t.Helper()
	foo := µ.GET(µ.FMap(func() error { return a }))
	req := mock.Input()

	it.Ok(t).
		If(foo(req)).Should().
		Assert(
			func(be interface{}) bool {
				var out *µ.Output
				if errors.As(be.(error), &out) {
					return reflect.DeepEqual(b, out)
				}
				return false
			},
		)
}

func issue(t *testing.T, a, b *µ.Issue) {
	t.Helper()
	foo := µ.GET(µ.FMap(func() error { return a }))
	req := mock.Input()

	it.Ok(t).
		If(foo(req)).Should().
		Assert(
			func(be interface{}) bool {
				var out *µ.Issue
				if errors.As(be.(error), &out) {
					return reflect.DeepEqual(b, out)
				}
				return false
			},
		)
}
