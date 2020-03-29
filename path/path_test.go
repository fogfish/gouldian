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

package path_test

import (
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/gouldian/path"
	"github.com/fogfish/it"
)

func TestPathIs(t *testing.T) {
	foo := µ.GET(µ.Path(path.Is("foo")))
	success := mock.Input(mock.URL("/foo"))
	failure := mock.Input(mock.URL("/bar"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestPathAny(t *testing.T) {
	foo := µ.GET(µ.Path(path.Is("foo"), path.Any()))
	bar := µ.GET(µ.Path(path.Is("foo"), path.Is("*")))
	success1 := mock.Input(mock.URL("/foo/bar"))
	success2 := mock.Input(mock.URL("/foo/foo"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(bar(success1)).Should().Equal(nil).
		If(bar(success2)).Should().Equal(nil)
}

func TestPathString(t *testing.T) {
	var value string
	foo := µ.GET(µ.Path(path.Is("foo"), path.String(&value)))
	success1 := mock.Input(mock.URL("/foo/bar"))
	success2 := mock.Input(mock.URL("/foo/1"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(value).Should().Equal("bar").
		//
		If(foo(success2)).Should().Equal(nil).
		If(value).Should().Equal("1")
}

func TestPathInt(t *testing.T) {
	var value int
	foo := µ.GET(µ.Path(path.Is("foo"), path.Int(&value)))
	success := mock.Input(mock.URL("/foo/1"))
	failure := mock.Input(mock.URL("/foo/bar"))

	it.Ok(t).
		If(foo(success)).Should().Equal(nil).
		If(value).Should().Equal(1).
		//
		If(foo(failure)).ShouldNot().Equal(nil)
}

func TestPathOr(t *testing.T) {
	foo := µ.GET(µ.Path(path.Or(path.Is("foo"), path.Is("bar"))))
	success1 := mock.Input(mock.URL("/foo"))
	success2 := mock.Input(mock.URL("/bar"))
	failure := mock.Input(mock.URL("/baz"))

	it.Ok(t).
		If(foo(success1)).Should().Equal(nil).
		If(foo(success2)).Should().Equal(nil).
		If(foo(failure)).ShouldNot().Equal(nil)
}
