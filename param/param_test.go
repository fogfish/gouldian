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

package param_test

import (
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/param"
	"github.com/fogfish/it"
)

func TestParamIs(t *testing.T) {
	foo := µ.GET(µ.Param(param.Is("foo", "bar")))
	success := mock.Input(mock.URL("/?foo=bar"))
	failure := mock.Input(mock.URL("/?bar=foo"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestParamAny(t *testing.T) {
	foo := µ.GET(µ.Param(param.Any("foo")))
	success1 := mock.Input(mock.URL("/?foo"))
	success2 := mock.Input(mock.URL("/?foo=bar"))
	success3 := mock.Input(mock.URL("/?foo=baz"))
	failure := mock.Input(mock.URL("/?bar=foo"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(foo(success3)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestParamString(t *testing.T) {
	var value string
	foo := µ.GET(µ.Param(param.String("foo", &value)))
	success1 := mock.Input(mock.URL("/?foo=bar"))
	success2 := mock.Input(mock.URL("/?foo=1"))
	failure := mock.Input(mock.URL("/?bar=foo"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(value).Should().Equal("bar").
		//
		If(foo(success2)).Should().Equal(nil).
		If(value).Should().Equal("1").
		//
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestParamMaybeString(t *testing.T) {
	var value string
	foo := µ.GET(µ.Param(param.MaybeString("foo", &value)))
	success1 := mock.Input(mock.URL("/?foo=bar"))
	success2 := mock.Input(mock.URL("/?bar=foo"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(value).Should().Equal("bar").
		//
		If(foo(success2)).Should().Equal(nil).
		If(value).Should().Equal("")
}

func TestParamInt(t *testing.T) {
	var value int
	foo := µ.GET(µ.Param(param.Int("foo", &value)))
	success := mock.Input(mock.URL("/?foo=1"))
	failure1 := mock.Input(mock.URL("/?foo=bar"))
	failure2 := mock.Input(mock.URL("/?bar=foo"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(1).
		//
		If(foo(failure1)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil)
}

func TestParamMaybeInt(t *testing.T) {
	var value int
	foo := µ.GET(µ.Param(param.MaybeInt("foo", &value)))
	success := mock.Input(mock.URL("/?foo=1"))
	failure1 := mock.Input(mock.URL("/?foo=bar"))
	failure2 := mock.Input(mock.URL("/?bar=foo"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(1).
		//
		If(foo(failure1)).Should().Equal(nil).
		If(value).Should().Equal(0).
		If(foo(failure2)).Should().Equal(nil).
		If(value).Should().Equal(0)
}

type MyT struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func TestParamJSON(t *testing.T) {
	var value MyT
	foo := µ.GET(µ.Param(param.JSON("foo", &value)))
	success := mock.Input(mock.URL("/?foo=%7B%22a%22%3A%22abc%22%2C%22b%22%3A10%7D"))
	failure1 := mock.Input(mock.URL("/?foo=bar"))
	failure2 := mock.Input(mock.URL("/?bar=foo"))
	failure3 := mock.Input(mock.Param("foo", "%7"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(MyT{A: "abc", B: 10}).
		//
		If(foo(failure1)).ShouldNot().Equal(nil).
		If(foo(failure2)).ShouldNot().Equal(nil).
		If(foo(failure3)).ShouldNot().Equal(nil)
}

func TestParamMaybeJSON(t *testing.T) {
	var value MyT
	foo := µ.GET(µ.Param(param.MaybeJSON("foo", &value)))
	success := mock.Input(mock.URL("/?foo=%7B%22a%22%3A%22abc%22%2C%22b%22%3A10%7D"))
	failure1 := mock.Input(mock.URL("/?foo=bar"))
	failure2 := mock.Input(mock.URL("/?bar=foo"))
	failure3 := mock.Input(mock.Param("foo", "%7"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(MyT{A: "abc", B: 10}).
		//
		If(foo(failure1)).Should().Equal(nil).
		If(foo(failure2)).Should().Equal(nil).
		If(foo(failure3)).Should().Equal(nil)
}

func TestParamOr(t *testing.T) {
	foo := µ.GET(µ.Param(
		param.Or(param.Any("foo"), param.Any("bar")),
	))
	success1 := mock.Input(mock.URL("/?foo"))
	success2 := mock.Input(mock.URL("/?bar"))
	failure := mock.Input(mock.URL("/?baz"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}
