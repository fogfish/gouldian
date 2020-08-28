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
	endpoint(mock.Input(mock.URL("/?foo=bar"))) == nil

*/
package param

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	µ "github.com/fogfish/gouldian"
)

// Or is a co-product of query param match arrows
//   e := µ.GET( µ.Param(param.Or(param.Is("foo", "bar"), param.Is("bar", "foo"))) )
//   e(mock.Input(mock.URL("/?foo=bar"))) == nil
//   e(mock.Input(mock.URL("/?bar=foo"))) == nil
//   e(mock.Input(mock.URL("/?foo=baz"))) != nil
func Or(arrows ...µ.ArrowParam) µ.ArrowParam {
	return func(params map[string]string) error {
		for _, f := range arrows {
			if err := f(params); !errors.Is(err, µ.NoMatch{}) {
				return err
			}
		}
		return µ.NoMatch{}
	}
}

// Is matches a param key to defined literal value
//   e := µ.GET( µ.Param(param.Is("foo", "bar")) )
//   e(mock.Input(mock.URL("/?foo=bar"))) == nil
//   e(mock.Input(mock.URL("/?bar=foo"))) != nil
func Is(key string, val string) µ.ArrowParam {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if exists && opt == val {
			return nil
		}
		return µ.NoMatch{}
	}
}

// Any is a wildcard matcher of param key. It fails if key is not defined.
//   e := µ.GET( µ.Param(param.Any("foo")) )
//   e(mock.Input(mock.URL("/?foo"))) == nil
//   e(mock.Input(mock.URL("/?foo=bar"))) == nil
//   e(mock.Input(mock.URL("/?foo=baz"))) == nil
//   e(mock.Input(mock.URL("/?bar=foo"))) != nil
func Any(key string) µ.ArrowParam {
	return func(params map[string]string) error {
		_, exists := params[key]
		if exists {
			return nil
		}
		return µ.NoMatch{}
	}
}

// String matches a param key to closed variable of string type.
// It fails if key is not defined.
//   var value string
//   e := µ.GET( µ.Param(param.String("foo", &value)) )
//   e(mock.Input(mock.URL("/?foo=bar"))) == nil && value == "bar"
//   e(mock.Input(mock.URL("/?foo=1"))) == nil && value == "1"
func String(key string, val *string) µ.ArrowParam {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if exists {
			*val = opt
			return nil
		}
		return µ.NoMatch{}
	}
}

// MaybeString matches a param key to closed variable of string type.
// It does not fail if key is not defined.
//   var value string
//   e := µ.GET( µ.Param(param.String("foo", &value)) )
//   e(mock.Input(mock.URL("/?foo=bar"))) == nil && value == "bar"
//   e(mock.Input(mock.URL("/?bar=1"))) == nil && value == ""
func MaybeString(key string, val *string) µ.ArrowParam {
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
//   e(mock.Input(mock.URL("/?foo=1"))) == nil && value == 1
//   e(mock.Input(mock.URL("/?foo=bar"))) != nil
func Int(key string, val *int) µ.ArrowParam {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if !exists {
			return µ.NoMatch{}
		}

		value, err := strconv.Atoi(opt)
		if err != nil {
			return µ.NoMatch{}
		}

		*val = value
		return nil
	}
}

// MaybeInt matches a param key to closed variable of int type.
// It does not fail if key is not defined.
//   var value int
//   e := µ.GET( µ.Param(param.Int("foo", &value)) )
//   e(mock.Input(mock.URL("/?foo=1"))) == nil && value == 1
//   e(mock.Input(mock.URL("/?foo=bar"))) == nil && value == 0
func MaybeInt(key string, val *int) µ.ArrowParam {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if !exists {
			return nil
		}

		value, err := strconv.Atoi(opt)
		if err != nil {
			*val = 0
			return nil
		}

		*val = value
		return nil
	}
}

// JSON matches a param key to closed struct.
// It assumes that key holds JSON value as url encoded string
func JSON(key string, val interface{}) µ.ArrowParam {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if !exists {
			return µ.NoMatch{}
		}

		str, err := url.QueryUnescape(opt)
		if err != nil {
			return µ.NoMatch{}
		}

		if err := json.Unmarshal([]byte(str), val); err != nil {
			// TODO: pass error to api client
			return µ.NoMatch{}
		}
		return nil
	}
}

// MaybeJSON matches a param key to closed struct.
// It assumes that key holds JSON value as url encoded string.
// It does not fail if key is not defined.
func MaybeJSON(key string, val interface{}) µ.ArrowParam {
	return func(params map[string]string) error {
		opt, exists := params[key]
		if !exists {
			return nil
		}

		str, err := url.QueryUnescape(opt)
		if err != nil {
			return nil
		}

		json.Unmarshal([]byte(str), val)
		return nil
	}
}
