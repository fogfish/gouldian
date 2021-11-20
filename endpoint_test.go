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
	"testing"

	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestEndpointThen(t *testing.T) {
	var ok = errors.New("b")
	var a µ.Endpoint = func(x *µ.Input) error { return nil }
	var b µ.Endpoint = func(x *µ.Input) error { return ok }
	var c µ.Endpoint = a.Then(b)

	it.Ok(t).
		If(c(mock.Input())).Should().Equal(ok)
}

func TestEndpointOr(t *testing.T) {
	var ok = errors.New("a")
	var a µ.Endpoint = func(x *µ.Input) error { return ok }
	var b µ.Endpoint = func(x *µ.Input) error { return µ.NoMatch{} }

	t.Run("a", func(t *testing.T) {
		var c µ.Endpoint = a.Or(b)

		it.Ok(t).
			If(c(mock.Input())).Should().Equal(ok)
	})

	t.Run("b", func(t *testing.T) {
		var c µ.Endpoint = b.Or(a)

		it.Ok(t).
			If(c(mock.Input())).Should().Equal(ok)
	})

}
