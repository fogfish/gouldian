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

Package param defines primitives to match URL Query parameters of HTTP requests.

	import "github.com/fogfish/gouldian/param"

	endpoint := µ.GET( µ.Param(param.Is("foo", "bar"), ...) )
	endpoint.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == true

*/
package param

import (
	"errors"
	"strconv"

	"github.com/fogfish/gouldian/core"
)

// Arrow is a type-safe definition of URL Query matcher
type Arrow func(map[string]string) error

// Or is a co-product of query param match arrows
//   e := µ.GET( µ.Param(param.Or(param.Is("foo", "bar"), param.Is("bar", "foo"))) )
//   e.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == true
//   e.IsMatch(mock.Input(mock.URL("/?bar=foo"))) == true
//   e.IsMatch(mock.Input(mock.URL("/?foo=baz"))) == false
func Or(arrows ...Arrow) Arrow {
	return func(params map[string]string) error {
		for _, f := range arrows {
			if err := f(params); !errors.Is(err, core.NoMatch{}) {
				return err
			}
		}
		return core.NoMatch{}
	}
}

// Is matches a param key to defined literal value
//   e := µ.GET( µ.Param(param.Is("foo", "bar")) )
//   e.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == true
//   e.IsMatch(mock.Input(mock.URL("/?bar=foo"))) == false
func Is(key string, val string) Arrow {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if exists && opt == val {
			return nil
		}
		return core.NoMatch{}
	}
}

// Any is a wildcard matcher of param key. It fails if key is not defined.
//   e := µ.GET( µ.Param(param.Any("foo")) )
//   e.IsMatch(mock.Input(mock.URL("/?foo"))) == true
//   e.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == true
//   e.IsMatch(mock.Input(mock.URL("/?foo=baz"))) == true
//   e.IsMatch(mock.Input(mock.URL("/?bar=foo"))) == false
func Any(key string) Arrow {
	return func(params map[string]string) error {
		_, exists := params[key]
		if exists {
			return nil
		}
		return core.NoMatch{}
	}
}

// String matches a param key to closed variable of string type.
// It fails if key is not defined.
//   var value string
//   e := µ.GET( µ.Param(param.String("foo", &value)) )
//   e.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == true && value == "bar"
//   e.IsMatch(mock.Input(mock.URL("/?foo=1"))) == true && value == "1"
func String(key string, val *string) Arrow {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if exists {
			*val = opt
			return nil
		}
		return core.NoMatch{}
	}
}

// MaybeString matches a param key to closed variable of string type.
// It does not fail if key is not defined.
//   var value string
//   e := µ.GET( µ.Param(param.String("foo", &value)) )
//   e.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == true && value == "bar"
//   e.IsMatch(mock.Input(mock.URL("/?foo=1"))) == true && value == "1"
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

// Int matches a param key to closed variable of int type.
// It fails if key is not defined.
//   var value int
//   e := µ.GET( µ.Param(param.Int("foo", &value)) )
//   e.IsMatch(mock.Input(mock.URL("/?foo=1"))) == true && value == 1
//   e.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == false
func Int(key string, val *int) Arrow {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if exists {
			if value, err := strconv.Atoi(opt); err == nil {
				*val = value
				return nil
			}
		}
		return core.NoMatch{}
	}
}

// MaybeInt matches a param key to closed variable of int type.
// It does not fail if key is not defined.
//   var value int
//   e := µ.GET( µ.Param(param.Int("foo", &value)) )
//   e.IsMatch(mock.Input(mock.URL("/?foo=1"))) == true && value == 1
//   e.IsMatch(mock.Input(mock.URL("/?foo=bar"))) == false
func MaybeInt(key string, val *int) Arrow {
	return func(params map[string]string) error {
		opt, exists := params[key]
		*val = 0
		if exists {
			if value, err := strconv.Atoi(opt); err == nil {
				*val = value
			}
		}
		return nil
	}
}
