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
	"reflect"
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestSuccess(t *testing.T) {
	output(t, µ.Success(200), µ.Success(200))
	output(t, µ.Ok(), µ.Success(200))
	output(t, µ.Created(), µ.Success(201))
	output(t, µ.Accepted(), µ.Success(202))
	output(t, µ.NoContent(), µ.Success(204))
}

func TestIssue(t *testing.T) {
	issue(t, µ.Failure(500), µ.Failure(500))
	issue(t, µ.BadRequest("issue"), µ.Failure(400).Reason("issue"))
	issue(t, µ.Unauthorized("issue"), µ.Failure(401).Reason("issue"))
	issue(t, µ.Forbidden("issue"), µ.Failure(403).Reason("issue"))
	issue(t, µ.NotFound("issue"), µ.Failure(404).Reason("issue"))
	issue(t, µ.InternalServerError("issue"), µ.Failure(500).Reason("issue"))
	issue(t, µ.NotImplemented("issue"), µ.Failure(501).Reason("issue"))
	issue(t, µ.ServiceUnavailable("issue"), µ.Failure(503).Reason("issue"))
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
