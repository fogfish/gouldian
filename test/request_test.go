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
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestPath(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(µ.URI(µ.Path("foo"))),
	)
	bar := mock.Endpoint(
		µ.GET(µ.URI(µ.Path("bar"))),
	)
	foobar := mock.Endpoint(µ.GET(µ.URI(µ.Path("foo"), µ.Path("bar"))))

	req := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestPathRoot(t *testing.T) {
	root := mock.Endpoint(
		µ.GET(µ.URI()),
	)

	success := mock.Input(mock.URL("/"))
	failure := mock.Input(mock.URL("/foo"))

	it.Ok(t).
		If(root(success)).Should().Equal(nil).
		If(root(failure)).ShouldNot().Equal(nil)
}

func TestParam(t *testing.T) {
	foo := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Param("foo", "bar"),
		),
	)
	bar := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Param("bar", "foo"),
		),
	)
	foobar := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Param("foo", "bar"),
			µ.Param("bar", "foo"),
		),
	)

	req := mock.Input(mock.URL("/?foo=bar"))

	it.Ok(t).
		If(foo(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}

func TestHeader(t *testing.T) {
	foo1 := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Header("foo", "bar"),
		),
	)

	bar := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Header("bar", "foo"),
		),
	)
	foobar := mock.Endpoint(
		µ.GET(
			µ.URI(),
			µ.Header("foo", "bar"),
			µ.Header("bar", "foo"),
		),
	)

	req := mock.Input(mock.Header("foo", "bar"))

	it.Ok(t).
		If(foo1(req)).Should().Equal(nil).
		If(bar(req)).ShouldNot().Equal(nil).
		If(foobar(req)).ShouldNot().Equal(nil)
}
