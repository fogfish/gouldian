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

package core_test

import (
	"errors"
	"testing"

	"github.com/fogfish/gouldian/core"
	"github.com/fogfish/gouldian/mock"
	"github.com/fogfish/it"
)

func TestEndpointThen(t *testing.T) {
	var ok = errors.New("b")
	var a core.Endpoint = func(x *core.Input) error { return nil }
	var b core.Endpoint = func(x *core.Input) error { return ok }
	var c core.Endpoint = a.Then(b)

	it.Ok(t).
		If(c(mock.Input())).Should().Equal(ok)
}

func TestEndpointOr(t *testing.T) {
	var ok = errors.New("a")
	var a core.Endpoint = func(x *core.Input) error { return ok }
	var b core.Endpoint = func(x *core.Input) error { return nil }
	var c core.Endpoint = a.Or(b)

	it.Ok(t).
		If(c(mock.Input())).Should().Equal(ok)
}
