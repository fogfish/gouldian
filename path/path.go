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
	"strconv"

	µ "github.com/fogfish/gouldian"
)

// Is matches a path segment to defined literal
//   e := µ.GET( µ.Path(path.Is("foo")) )
//   e(mock.Input(mock.URL("/foo"))) == nil
//   e(mock.Input(mock.URL("/bar"))) != nil
func Is(val string) µ.ArrowPath {
	if val == "*" {
		return func([]string) error { return nil }
	}

	return func(segments []string) error {
		if segments[0] == val {
			return nil
		}
		return µ.NoMatch{}
	}
}

// Any is a wildcard matcher of path segment
//   e := µ.GET( µ.Path(path.Any()) )
//   e(mock.Input(mock.URL("/foo"))) == nil
//   e(mock.Input(mock.URL("/bar"))) == nil
func Any() µ.ArrowPath {
	return func([]string) error {
		return nil
	}
}

// String matches a path segment to closed variable of string type.
//   var value string
//   e := µ.GET( µ.Path(path.String(&value)) )
//   e(mock.Input(mock.URL("/foo"))) == nil && value == "foo"
//   e(mock.Input(mock.URL("/1"))) == nil && value == "1"
func String(val *string) µ.ArrowPath {
	return func(segments []string) error {
		*val = segments[0]
		return nil
	}
}

// Int matches a path segment to closed variable of int type
//   var value int
//   e := µ.GET( µ.Path(path.Int(&value)) )
//   e(mock.Input(mock.URL("/1"))) == nil && value == 1
//   e(mock.Input(mock.URL("/foo"))) != nil
func Int(val *int) µ.ArrowPath {
	return func(segments []string) error {
		if value, err := strconv.Atoi(segments[0]); err == nil {
			*val = value
			return nil
		}
		return µ.NoMatch{}
	}
}
