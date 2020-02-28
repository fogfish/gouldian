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

/*

Package path defines primitives to match URL of HTTP requests.

	import "github.com/fogfish/gouldian/path"

	endpoint := µ.GET( µ.Path(path.Is("foo"), ...) )
	endpoint(mock.Input(mock.URL("/foo"))) == nil

*/
package path

import (
	"errors"
	"strconv"

	"github.com/fogfish/gouldian/core"
)

// Arrow is a type-safe definition of URL segment matcher
type Arrow func(string) error

// Or is a co-product of path match arrows
//   e := µ.GET( µ.Path(path.Or(path.Is("foo"), path.Is("bar"))) )
//   e(mock.Input(mock.URL("/foo"))) == nil
//   e(mock.Input(mock.URL("/bar"))) == nil
//   e(mock.Input(mock.URL("/baz"))) != nil
func Or(arrows ...Arrow) Arrow {
	return func(segment string) error {
		for _, f := range arrows {
			if err := f(segment); !errors.Is(err, core.NoMatch{}) {
				return err
			}
		}
		return core.NoMatch{}
	}
}

// Is matches a path segment to defined literal
//   e := µ.GET( µ.Path(path.Is("foo")) )
//   e(mock.Input(mock.URL("/foo"))) == nil
//   e(mock.Input(mock.URL("/bar"))) != nil
func Is(val string) Arrow {
	if val == "*" {
		return func(string) error { return nil }
	}

	return func(segment string) error {
		if segment == val {
			return nil
		}
		return core.NoMatch{}
	}
}

// Any is a wildcard matcher of path segment
//   e := µ.GET( µ.Path(path.Any()) )
//   e(mock.Input(mock.URL("/foo"))) == nil
//   e(mock.Input(mock.URL("/bar"))) == nil
func Any() Arrow {
	return func(string) error {
		return nil
	}
}

// String matches a path segment to closed variable of string type.
//   var value string
//   e := µ.GET( µ.Path(path.String(&value)) )
//   e(mock.Input(mock.URL("/foo"))) == nil && value == "foo"
//   e(mock.Input(mock.URL("/1"))) == nil && value == "1"
func String(val *string) Arrow {
	return func(segment string) error {
		*val = segment
		return nil
	}
}

// Int matches a path segment to closed variable of int type
//   var value int
//   e := µ.GET( µ.Path(path.Int(&value)) )
//   e(mock.Input(mock.URL("/1"))) == nil && value == 1
//   e(mock.Input(mock.URL("/foo"))) != nil
func Int(val *int) Arrow {
	return func(segment string) error {
		if value, err := strconv.Atoi(segment); err == nil {
			*val = value
			return nil
		}
		return core.NoMatch{}
	}
}
