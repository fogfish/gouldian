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

package header_test

import (
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/header"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestHeaderIs(t *testing.T) {
	foo := µ.GET(µ.Header(header.Is("Content-Type", "application/json")))
	success := mock.Input(mock.Header("Content-Type", "application/json"))
	failure := mock.Input(mock.Header("Content-Type", "text/plain"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestHeaderIsLowerCase(t *testing.T) {
	foo := µ.GET(µ.Header(header.Is("Content-Type", "application/json")))
	success := mock.Input(mock.Header("content-type", "application/json"))
	failure := mock.Input(mock.Header("content-type", "text/plain"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestContentJSON(t *testing.T) {
	foo := µ.GET(µ.Header(header.ContentJSON()))
	success := mock.Input(mock.Header("Content-Type", "application/json"))
	failure := mock.Input(mock.Header("Content-Type", "text/plain"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestContentForm(t *testing.T) {
	foo := µ.GET(µ.Header(header.ContentForm()))
	success := mock.Input(mock.Header("Content-Type", "application/x-www-form-urlencoded"))
	failure := mock.Input(mock.Header("Content-Type", "text/plain"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestHeaderAny(t *testing.T) {
	foo := µ.GET(µ.Header(header.Any("Content-Type")))
	success1 := mock.Input(mock.Header("Content-Type", "application/json"))
	success2 := mock.Input(mock.Header("Content-Type", "text/plain"))
	failure := mock.Input()

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestHeaderString(t *testing.T) {
	var value string
	foo := µ.GET(µ.Header(header.String("Content-Type", &value)))
	success := mock.Input(mock.Header("Content-Type", "application/json"))
	failure := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal("application/json").
		//
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestHeaderMaybeString(t *testing.T) {
	var value string
	foo := µ.GET(µ.Header(header.MaybeString("Content-Type", &value)))
	success1 := mock.Input(mock.Header("Content-Type", "application/json"))
	success2 := mock.Input()

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(value).Should().Equal("application/json").
		//
		If(foo(success2)).Should().Equal(nil).
		If(value).Should().Equal("")
}

func TestHeaderInt(t *testing.T) {
	var value int
	foo := µ.GET(µ.Header(header.Int("Content-Length", &value)))
	success := mock.Input(mock.Header("Content-Length", "1024"))
	failure := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(1024).
		//
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestHeaderMaybeInt(t *testing.T) {
	var value int
	foo := µ.GET(µ.Header(header.MaybeInt("Content-Length", &value)))
	success := mock.Input(mock.Header("Content-Length", "1024"))
	failure := mock.Input()

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(1024).
		//
		If(foo(failure)).Should().Equal(nil).
		If(value).Should().Equal(0)
}

func TestParamOr(t *testing.T) {
	foo := µ.GET(µ.Header(
		header.Or(
			header.Is("Content-Type", "application/json"),
			header.Is("Content-Type", "text/html"),
		),
	))

	success1 := mock.Input(mock.Header("Content-Type", "application/json"))
	success2 := mock.Input(mock.Header("Content-Type", "text/html"))
	failure := mock.Input(mock.Header("Content-Type", "text/plain"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}
