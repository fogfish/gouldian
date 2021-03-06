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
	endpoint(mock.Input(Json)) == nil

*/
package header

import (
	"errors"
	"fmt"
	"net/textproto"
	"strconv"
	"strings"

	µ "github.com/fogfish/gouldian"
)

// Or is a co-product of header match arrows
//   e := µ.GET(
//     µ.Header(
//       header.Or(
//         header.Is("Content-Type", "application/json"),
//         header.Is("Content-Type", "text/plain"),
//       )
//     )
//   )
//   e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil
//   e(mock.Input(mock.Header("Content-Type", "text/plain"))) == nil
//   e(mock.Input(mock.Header("Content-Type", "text/html"))) != nil
func Or(arrows ...µ.ArrowHeader) µ.ArrowHeader {
	return func(headers map[string]string) error {
		for _, f := range arrows {
			if err := f(headers); !errors.Is(err, µ.NoMatch{}) {
				return err
			}
		}
		return µ.NoMatch{}
	}
}

// value return header with lower case support
func value(headers map[string]string, key string) (string, bool) {
	k := textproto.CanonicalMIMEHeaderKey(key)
	v, exists := headers[k]
	if !exists {
		// Note: required due to browser behavior
		v, exists = headers[strings.ToLower(k)]
		return v, exists
	}
	return v, exists
}

// Is matches a header to defined literal value
//   e := µ.GET( µ.Header(header.Is("Content-Type", "application/json")) )
//   e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil
//   e(mock.Input(mock.Header("Content-Type", "text/plain"))) != nil
func Is(key string, val string) µ.ArrowHeader {
	return func(headers map[string]string) error {
		opt, exists := value(headers, key)
		if exists && strings.HasPrefix(opt, val) {
			return nil
		}
		return µ.NoMatch{}
	}
}

// ContentJSON is a syntax sugar to header.Is("Content-Type", "application/json")
func ContentJSON() µ.ArrowHeader {
	return Is("Content-Type", "application/json")
}

// ContentForm is a syntax sugar to header.Is("Content-Type", "application/x-www-form-urlencoded")
func ContentForm() µ.ArrowHeader {
	return Is("Content-Type", "application/x-www-form-urlencoded")
}

// Any is a wildcard matcher of header. It fails if header is not defined.
//   e := µ.GET( µ.Header(header.Any("Content-Type")) )
//   e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil
//   e(mock.Input(mock.Header("Content-Type", "text/plain"))) == nil
//   e(mock.Input()) != nil
func Any(key string) µ.ArrowHeader {
	return func(headers map[string]string) error {
		_, exists := value(headers, key)
		if exists {
			return nil
		}
		return µ.NoMatch{}
	}
}

// String matches a header value to closed variable of string type.
// It fails if header is not defined.
//   var value string
//   e := µ.GET( µ.Header(header.String("Content-Type", &value)) )
//   e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil && value == "application/json"
//   e(mock.Input()) != nil
func String(key string, val *string) µ.ArrowHeader {
	return func(headers map[string]string) error {
		opt, exists := value(headers, key)
		if exists {
			*val = opt
			return nil
		}
		return µ.NoMatch{}
	}
}

// MaybeString matches a header value to closed variable of string type.
// It does not fail if header is not defined.
//   var value string
//   e := µ.GET( µ.Header(header.String("foo", &value)) )
//   e(mock.Input(mock.Header("Content-Type", "application/json"))) == nil && value == "application/json"
//   e(mock.Input()) == nil
func MaybeString(key string, val *string) µ.ArrowHeader {
	return func(headers map[string]string) error {
		opt, exists := value(headers, key)
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
//   e(mock.Input(mock.Header("Content-Length", "1024"))) == nil && value == 1024
//   e(mock.Input()) != nil
func Int(key string, val *int) µ.ArrowHeader {
	return func(headers map[string]string) error {
		opt, exists := value(headers, key)
		if exists {
			if value, err := strconv.Atoi(opt); err == nil {
				*val = value
				return nil
			}
		}
		return µ.NoMatch{}
	}
}

// MaybeInt matches a header value to closed variable of int type.
// It does not fail if header is not defined.
//   var value int
//   e := µ.GET( µ.Header(header.MaybeInt("Content-Length", &value)) )
//   e(mock.Input(mock.Header("Content-Length", "1024"))) == nil && value == 1024
//   e(mock.Input()) == nil
func MaybeInt(key string, val *int) µ.ArrowHeader {
	return func(headers map[string]string) error {
		opt, exists := value(headers, key)
		*val = 0
		if exists {
			if value, err := strconv.Atoi(opt); err == nil {
				*val = value
			}
		}
		return nil
	}
}

// Authorize validates content of HTTP Authorization header
func Authorize(authType string, f func(string) error) µ.Endpoint {
	return func(req *µ.Input) error {
		auth, exists := req.Header("Authorization")
		if !exists {
			return µ.Unauthorized(fmt.Errorf("Unauthorized %v", req.APIGatewayProxyRequest.Path))
		}

		cred := strings.Split(auth, " ")
		if len(cred) != 2 {
			return µ.Unauthorized(fmt.Errorf("Unauthorized %v", req.APIGatewayProxyRequest.Path))
		}

		if strings.ToLower(cred[0]) != strings.ToLower(authType) {
			return µ.Unauthorized(fmt.Errorf("Unauthorized %v", req.APIGatewayProxyRequest.Path))
		}

		if err := f(cred[1]); err != nil {
			return µ.Unauthorized(err)
		}

		return nil
	}
}
