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

Package header defines primitives to match Headers of HTTP requests.

	import "github.com/fogfish/gouldian/header"

	endpoint := µ.GET(
		µ.Header(
			header.Is("Content-Type", "application/json"),
			...
		)
	)
	Json := mock.Header("Content-Type", "application/json")
	endpoint.IsMatch(mock.Input(Json)) == true

*/
package header

import (
	"errors"
	"strconv"

	"github.com/fogfish/gouldian/core"
)

// Arrow is a type-safe definition of URL Query matcher
type Arrow func(map[string]string) error

// Or is a co-product of query param match arrows
//   e := µ.GET(
//     µ.Header(
//       header.Or(
//         header.Is("Content-Type", "application/json"),
//         header.Is("Content-Type", "text/plain"),
//       )
//     )
//   )
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "application/json"))) == true
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "text/plain"))) == true
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "text/html"))) == false
func Or(arrows ...Arrow) Arrow {
	return func(headers map[string]string) error {
		for _, f := range arrows {
			if err := f(headers); !errors.Is(err, core.NoMatch{}) {
				return err
			}
		}
		return core.NoMatch{}
	}
}

// Is matches a header to defined literal value
//   e := µ.GET( µ.Header(header.Is("Content-Type", "application/json")) )
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "application/json"))) == true
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "text/plain"))) == false
func Is(key string, val string) Arrow {
	return func(headers map[string]string) error {
		opt, exists := headers[key]
		if exists && opt == val {
			return nil
		}
		return core.NoMatch{}
	}
}

// Any is a wildcard matcher of header. It fails if header is not defined.
//   e := µ.GET( µ.Header(header.Any("Content-Type")) )
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "application/json"))) == true
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "text/plain"))) == true
//   e.IsMatch(mock.Input()) == false
func Any(key string) Arrow {
	return func(headers map[string]string) error {
		_, exists := headers[key]
		if exists {
			return nil
		}
		return core.NoMatch{}
	}
}

// String matches a header value to closed variable of string type.
// It fails if header is not defined.
//   var value string
//   e := µ.GET( µ.Header(header.String("Content-Type", &value)) )
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "application/json"))) == true && value == "application/json"
//   e.IsMatch(mock.Input()) == false
func String(key string, val *string) Arrow {
	return func(headers map[string]string) error {
		opt, exists := headers[key]
		if exists {
			*val = opt
			return nil
		}
		return core.NoMatch{}
	}
}

// MaybeString matches a header value to closed variable of string type.
// It does not fail if header is not defined.
//   var value string
//   e := µ.GET( µ.Header(header.String("foo", &value)) )
//   e.IsMatch(mock.Input(mock.Header("Content-Type", "application/json"))) == true && value == "application/json"
//   e.IsMatch(mock.Input()) == true
func MaybeString(key string, val *string) Arrow {
	return func(params map[string]string) error {
		opt, exists := params[key]
		*val = ""
		if exists {
			*val = opt
		}
		return nil
	}
}

// Int matches a header value to closed variable of int type.
// It fails if header is not defined.
//   var value int
//   e := µ.GET( µ.Header(header.Int("Content-Length", &value)) )
//   e.IsMatch(mock.Input(mock.Header("Content-Length", "1024"))) == true && value == 1024
//   e.IsMatch(mock.Input()) == false
func Int(key string, val *int) Arrow {
	return func(headers map[string]string) error {
		opt, exists := headers[key]
		if exists {
			if value, err := strconv.Atoi(opt); err == nil {
				*val = value
				return nil
			}
		}
		return core.NoMatch{}
	}
}

// MaybeInt matches a header value to closed variable of int type.
// It does not fail if header is not defined.
//   var value int
//   e := µ.GET( µ.Header(header.MaybeInt("Content-Length", &value)) )
//   e.IsMatch(mock.Input(mock.Header("Content-Length", "1024"))) == true && value == 1024
//   e.IsMatch(mock.Input()) == true
func MaybeInt(key string, val *int) Arrow {
	return func(headers map[string]string) error {
		opt, exists := headers[key]
		*val = 0
		if exists {
			if value, err := strconv.Atoi(opt); err == nil {
				*val = value
			}
		}
		return nil
	}
}
