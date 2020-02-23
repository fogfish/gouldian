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

func TestParamIs(t *testing.T) {
	foo := µ.GET(µ.Header(header.Is("Content-Type", "application/json")))
	success := mock.Input(mock.Header("Content-Type", "application/json"))
	failure := mock.Input(mock.Header("Content-Type", "text/plain"))

	it.Ok(t).
		If(foo.IsMatch(success)).Should().Equal(true).
		If(foo.IsMatch(failure)).Should().Equal(false)
}

func TestParamAny(t *testing.T) {
	foo := µ.GET(µ.Header(header.Any("Content-Type")))
	success1 := mock.Input(mock.Header("Content-Type", "application/json"))
	success2 := mock.Input(mock.Header("Content-Type", "text/plain"))
	failure := mock.Input()

	it.Ok(t).
		If(foo.IsMatch(success1)).Should().Equal(true).
		If(foo.IsMatch(success2)).Should().Equal(true).
		If(foo.IsMatch(failure)).Should().Equal(false)
}

func TestParamString(t *testing.T) {
	var value string
	foo := µ.GET(µ.Header(header.String("Content-Type", &value)))
	success := mock.Input(mock.Header("Content-Type", "application/json"))
	failure := mock.Input()

	it.Ok(t).
		If(foo.IsMatch(success)).Should().Equal(true).
		If(value).Should().Equal("application/json").
		//
		If(foo.IsMatch(failure)).Should().Equal(false)
}

func TestParamMaybeString(t *testing.T) {
	var value string
	foo := µ.GET(µ.Header(header.MaybeString("Content-Type", &value)))
	success1 := mock.Input(mock.Header("Content-Type", "application/json"))
	success2 := mock.Input()

	it.Ok(t).
		If(foo.IsMatch(success1)).Should().Equal(true).
		If(value).Should().Equal("application/json").
		//
		If(foo.IsMatch(success2)).Should().Equal(true).
		If(value).Should().Equal("")
}

func TestParamInt(t *testing.T) {
	var value int
	foo := µ.GET(µ.Header(header.Int("Content-Length", &value)))
	success := mock.Input(mock.Header("Content-Length", "1024"))
	failure := mock.Input()

	it.Ok(t).
		If(foo.IsMatch(success)).Should().Equal(true).
		If(value).Should().Equal(1024).
		//
		If(foo.IsMatch(failure)).Should().Equal(false)
}

func TestParamMaybeInt(t *testing.T) {
	var value int
	foo := µ.GET(µ.Header(header.MaybeInt("Content-Length", &value)))
	success := mock.Input(mock.Header("Content-Length", "1024"))
	failure := mock.Input()

	it.Ok(t).
		If(foo.IsMatch(success)).Should().Equal(true).
		If(value).Should().Equal(1024).
		//
		If(foo.IsMatch(failure)).Should().Equal(true).
		If(value).Should().Equal(0)
}
